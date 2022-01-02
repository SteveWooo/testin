package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Analysis struct {
	Blocks       []map[string]interface{}
	TpsPerBlocks []float64
}

func (a *Analysis) FetchBlockInfo() map[string]interface{} {
	resp, err := http.Get("http://localhost:9024/sdk/consensus/get_info_by_name?bcag_name=test")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	jsonData := map[string]interface{}{}
	json.Unmarshal(body, &jsonData)
	data := jsonData["Data"].(map[string]interface{})
	return data
}

func (a *Analysis) GetBlock(blockNumber int) map[string]interface{} {
	resp, err := http.Get("http://localhost:9024/sdk/consensus/get_block_by_number?bcag_name=test&block_number=" + strconv.Itoa(blockNumber))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	jsonData := map[string]interface{}{}
	json.Unmarshal(body, &jsonData)
	block := jsonData["Data"].(map[string]interface{})["Block"].(map[string]interface{})
	return block
}

func (a *Analysis) FetchAllBlocks(blockCount float64) {
	for i := 1; i <= int(blockCount); i++ {
		block := a.GetBlock(i)
		a.Blocks = append(a.Blocks, block)
	}
}

func (a *Analysis) CalculateTpsPerBlock() {
	tpsPoints := plotter.XYs{}
	for i := 8; i < len(a.Blocks)-20; i++ {
		trans := a.Blocks[i]["Transactions"].([]interface{})
		thisBlockTime, _ := strconv.Atoi(a.Blocks[i]["Ts"].(string))
		lastBlockTims, _ := strconv.Atoi(a.Blocks[i-1]["Ts"].(string))
		tps := float64(len(trans)) / float64(thisBlockTime-lastBlockTims)
		// fmt.Println(len(trans), thisBlockTime, lastBlockTims)
		// fmt.Println(tps * 1000)
		tpsPoints = append(tpsPoints, plotter.XY{
			X: float64(i),
			Y: tps * 1000,
		})
	}

	var tpsDrawer Drawer = Drawer{}
	tpsDrawer.YLabel = "TPS (transaction per second)"
	tpsDrawer.Build("bussiness_tps", []float64{0, 0, 40, 100})
	tpsDrawer.SetPoint(tpsPoints, "tps")
	tpsDrawer.DoDraw()
}

func (a *Analysis) CalculatePackageEachMiner() {
	miners := map[string]int{}

	for i := 8; i < len(a.Blocks); i++ {
		if miners[a.Blocks[i]["Miner"].(string)] == 0 {
			miners[a.Blocks[i]["Miner"].(string)] = 0
		}
		miners[a.Blocks[i]["Miner"].(string)] = miners[a.Blocks[i]["Miner"].(string)] + 1
	}

	nodeSum := 0.0
	countSum := 0.0
	deviationSum := 0.0
	for _, count := range miners {
		nodeSum = nodeSum + 1
		countSum = countSum + float64(count)
	}

	minersPoints := plotter.XYs{}
	i := 0
	for nid, count := range miners {
		deviationSum = deviationSum + math.Pow(float64(count)-countSum/nodeSum, 2)
		fmt.Println(nid, count)
		minersPoints = append(minersPoints, plotter.XY{
			X: float64(i),
			Y: float64(count),
		})
		i = i + 1
	}
	fmt.Println("count: ", nodeSum, " avg: ", countSum/nodeSum, " deviation:", deviationSum/(nodeSum-1))

	var tpsDrawer Drawer = Drawer{}
	tpsDrawer.YLabel = "Package Times"
	tpsDrawer.Build("miner_package", []float64{0, 0, 15, 300})
	tpsDrawer.SetPoint(minersPoints, "tps")
	tpsDrawer.DoDraw()
}

func main() {
	analysis := Analysis{}
	info := analysis.FetchBlockInfo()
	analysis.FetchAllBlocks(info["BlockCount"].(float64))
	analysis.CalculateTpsPerBlock()
	analysis.CalculatePackageEachMiner()
}

// 一个图
type Drawer struct {
	points   []*Point
	fileName string
	rect     []float64

	YLabel string
}

func (drawer *Drawer) Build(filename string, rect []float64) {
	drawer.fileName = filename
	drawer.rect = make([]float64, 4)
	drawer.rect = rect
}

func (drawer *Drawer) SetPoint(point plotter.XYs, name string) {
	for i := 0; i < len(point); i++ {
		for k := 0; k < len(point); k++ {
			if point[i].X < point[k].X {
				temp := point[i]
				point[i] = point[k]
				point[k] = temp
			}
		}
	}

	var p Point = Point{
		point,
		name,
	}

	drawer.points = append(drawer.points, &p)
}

func (drawer *Drawer) DoDraw() {
	plt := plot.New()
	plt.X.Label.Text = "Block Number"
	plt.Y.Label.Text = drawer.YLabel
	var err error
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = drawer.rect[0], drawer.rect[1], drawer.rect[2], drawer.rect[3]

	params := []interface{}{}
	for i := 0; i < len(drawer.points); i++ {
		params = append(params, drawer.points[i].GetName(), drawer.points[i].GetPoints())
	}

	err = plotutil.AddLines(plt,
		params...,
	)
	if err != nil {
		panic(err)
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, drawer.fileName+".png"); err != nil {
		panic(err)
	}
}

// 一条线
type Point struct {
	points plotter.XYs
	name   string
}

func (point *Point) GetPoints() plotter.XYs {
	return point.points
}

func (point *Point) GetName() string {
	return point.name
}
