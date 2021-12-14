package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	sdkApi "github.com/stevewooo/testin/Modules/SdkApi"
	"github.com/stevewooo/testin/Modules/Transaction"
)

// 链上数据的实时世界状态
type WorldStatus struct {
	config map[string]string

	blocks []*Block // 链上所有区块

	// 状态
	Hackers     []*Transaction.Hacker
	Enterprises []*Transaction.Enterprise
	Experts     []*Transaction.Expert
	Tasks       []*Transaction.Task
}

func (w *WorldStatus) Build(config map[string]string) {
	w.config = config
}

// 以对象形式返回一个世界状态
func (w *WorldStatus) GetWorldStatus() map[string]interface{} {
	return map[string]interface{}{
		"Hackers":     w.Hackers,
		"Enterprises": w.Enterprises,
		"Experts":     w.Experts,
		"Tasks":       w.Tasks,
	}
}

func (w *WorldStatus) cleanStatus() {
	w.Hackers = []*Transaction.Hacker{}
	w.Enterprises = []*Transaction.Enterprise{}
	w.Experts = []*Transaction.Expert{}
	w.Tasks = []*Transaction.Task{}
}

// 读取交易数据，写入世界状态。
func (w *WorldStatus) parseTransaction(trans *Transaction.Transaction) {
	if trans.Type == "RegisterHacker" {
		w.Hackers = append(w.Hackers, &trans.Hacker)
	}

	if trans.Type == "RegisterEnterprise" {
		w.Enterprises = append(w.Enterprises, &trans.Enterprise)
	}

	if trans.Type == "RegisterExpert" {
		w.Experts = append(w.Experts, &trans.Expert)
	}

	if trans.Type == "PublishTaskByEnterprise" {
		w.Tasks = append(w.Tasks, &trans.Task)
	}
}

// 按照现有区块，构建世界状态
func (w *WorldStatus) DoBuildStatus() {
	w.cleanStatus()
	for i := 0; i < len(w.blocks); i++ {
		for k := 0; k < len(w.blocks[i].Transactions); k++ {
			w.parseTransaction(w.blocks[i].Transactions[k])
		}
	}
}

func (w *WorldStatus) GetLocalTopBlock() *Block {
	return w.blocks[len(w.blocks)-1]
}

func (w *WorldStatus) AddNewLocalBlock(block *Block) {
	w.blocks = append(w.blocks, block)
}

func (w *WorldStatus) GetRemoteTopBlock() *Block {
	topBlockRes, err := sdkApi.GetTopBlock(w.config)
	if err != nil {
		panic("sdk数据获取失败")
	}
	if topBlockRes.Status != 2000 {
		panic("sdk报错：" + topBlockRes.Message)
	}
	var topBlock Block = Block{}
	topBlock.LoadFromInterface(topBlockRes.Data.(map[string]interface{})["TopBlock"])
	return &topBlock
}

// 拉取缓存数据，用于构建区块
func (w *WorldStatus) GetRemoteTransactionCache() ([]*Transaction.Transaction, error) {
	localTopBlock := w.GetLocalTopBlock()
	if localTopBlock == nil {
		panic("远程区块不存在，本地缓存区块获取失败")
	}

	localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	newBlockNumber := strconv.Itoa(localTopBlockNumber + 1)
	prefix := "transCache-" + newBlockNumber + "-"
	transResp, err := sdkApi.GetCacheByPerfix(w.config, prefix)
	if err != nil {
		return nil, errors.New("获取远程区块缓存失败")
	}
	if transResp.Status != 2000 {
		return nil, errors.New("Sdk报错：" + transResp.Message)
	}

	// 构建缓存交易数据
	trans := []*Transaction.Transaction{}
	for i := 0; i < len(transResp.Data.(map[string]interface{})["Caches"].([]interface{})); i++ {
		tran := Transaction.Transaction{}
		err = tran.LoadFromJSONString(transResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
		if err != nil {
			fmt.Println(err)
		}

		// TODO通过世界状态检查交易
		trans = append(trans, &tran)
	}
	return trans, nil
}

// 拉取线上所有区块。因为不是每次构建都需要全量拉取，所以要独立封装出来
func (w *WorldStatus) FetchAllBlocks() {
	// 1 获取顶部区块
	topBlock := w.GetRemoteTopBlock()

	// 2 根据顶部区块，获取所有区块
	allBlockRes, err := sdkApi.GetBlockByRange(w.config, "1", topBlock.Number)
	if err != nil {
		panic("sdk数据获取失败")
	}
	if allBlockRes.Status != 2000 {
		panic("sdk报错：" + allBlockRes.Message)
	}
	// 清理本地区块缓存
	w.blocks = []*Block{}
	for i := 0; i < len(allBlockRes.Data.(map[string]interface{})["Blocks"].([]interface{})); i++ {
		var block Block = Block{}
		block.LoadFromInterface(allBlockRes.Data.(map[string]interface{})["Blocks"].([]interface{})[i])
		w.blocks = append(w.blocks, &block)
	}
}

// 把世界状态打包成JSON格式数据
func (w *WorldStatus) PackToJSONString() string {
	result := w.GetWorldStatus()
	jsonByte, _ := json.Marshal(result)
	return string(jsonByte)
}
