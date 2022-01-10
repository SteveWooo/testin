package main

import (
	modules "github.com/stevewooo/testin/tpsTesting/Modules"
)

func main() {
	analysisPoBR := modules.Analysis{}
	analysisPoBR.LoadAllBlocks("./block_PoBR_10per.json")
	analysisPoBR.CalculateDatasPerBlockAndSet()

	analysisPBFT := modules.Analysis{}
	analysisPBFT.LoadAllBlocks("./block_PBFT_10per.json")
	analysisPBFT.CalculateDatasPerBlockAndSet()

	tpsDrawer := modules.Drawer{}
	tpsDrawer.YLabel = "Tps (Transaction per second)"
	tpsDrawer.Build("img-Tps between PBFT and PoBR", []float64{0, 0, 15, 250})
	tpsDrawer.SetPoint(analysisPoBR.TpsPoints, "Tps of PoBR")
	tpsDrawer.SetPoint(analysisPBFT.TpsPoints, "Tps of PBFT")
	tpsDrawer.DoDraw()

	bcpbDrawer := modules.Drawer{}
	bcpbDrawer.YLabel = "Bcpb (Byte cost per block)"
	bcpbDrawer.Build("img-Bcpb between PBFT and PoBR", []float64{0, 0, 6500, 250})
	bcpbDrawer.SetPoint(analysisPoBR.ByteUsed, "Bcpb of PoBR")
	bcpbDrawer.SetPoint(analysisPBFT.ByteUsed, "Bcpb of PBFT")
	bcpbDrawer.DoDraw()

	bcpsDrawer := modules.Drawer{}
	bcpsDrawer.YLabel = "Bcps (Byte cost per second)"
	bcpsDrawer.Build("img-Bcps between PBFT and PoBR", []float64{0, 0, 2, 250})
	bcpsDrawer.SetPoint(analysisPoBR.BytePerSecond, "Bcps of PoBR")
	bcpsDrawer.SetPoint(analysisPBFT.BytePerSecond, "Bcps of PBFT")
	bcpsDrawer.DoDraw()
}
