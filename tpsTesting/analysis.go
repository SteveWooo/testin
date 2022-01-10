package main

import (
	modules "github.com/stevewooo/testin/tpsTesting/Modules"
)

func main() {
	analysis := modules.Analysis{}
	info := analysis.FetchBlockInfo()

	err := analysis.LoadAllBlocks("./block.json")
	if err != nil {
		analysis.FetchAllBlocks(info["BlockCount"].(float64))
		analysis.SaveAllBlocks()
	}

	analysis.CalculateTpsPerBlock()
	analysis.CalculatePackageEachMiner()
}
