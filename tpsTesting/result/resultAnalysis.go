package main

import (
	modules "github.com/stevewooo/testin/tpsTesting/Modules"
)

func main() {
	analysisPoBR := modules.Analysis{}
	analysisPoBR.LoadAllBlocks("./block_PoBR_15per_allTrans.json")
	// analysisPoBR.LoadAllBlocks("./block_PoBR_10per_allTrans.json")
	analysisPoBR.CalculateDatasPerBlockAndSet()

	analysisPBFT := modules.Analysis{}
	analysisPBFT.LoadAllBlocks("./block_PBFT_15per_allTrans.json")
	// analysisPBFT.LoadAllBlocks("./block_PBFT_10per_allTrans.json")
	analysisPBFT.CalculateDatasPerBlockAndSet()

	tpsDrawer := modules.Drawer{}
	tpsDrawer.YLabel = "Tps (每秒交易完成数)"
	tpsDrawer.Build("img-Tps between PBFT and PoBR", []float64{0, 0, 20, 250})
	tpsDrawer.SetPoint(analysisPoBR.TpsPoints, "PoBR 算法 Tps")
	tpsDrawer.SetPoint(analysisPBFT.TpsPoints, "PBFT 算法 Tps")
	tpsDrawer.DoDraw()

	bcpbDrawer := modules.Drawer{}
	bcpbDrawer.YLabel = "Bcpb (每区块通信字节消耗)"
	bcpbDrawer.Build("img-Bcpb between PBFT and PoBR", []float64{0, 0, 6500, 250})
	bcpbDrawer.SetPoint(analysisPoBR.ByteUsed, "PoBR 算法 Bcpb")
	bcpbDrawer.SetPoint(analysisPBFT.ByteUsed, "PBFT 算法 Bcpb")
	bcpbDrawer.DoDraw()

	bcpsDrawer := modules.Drawer{}
	bcpsDrawer.YLabel = "Bcps (每秒通信字节消耗)"
	bcpsDrawer.Build("img-Bcps between PBFT and PoBR", []float64{0, 0, 2, 250})
	bcpsDrawer.SetPoint(analysisPoBR.BytePerSecond, "PoBR 算法 Bcps")
	bcpsDrawer.SetPoint(analysisPBFT.BytePerSecond, "PBFT 算法 Bcps")
	bcpsDrawer.DoDraw()
}
