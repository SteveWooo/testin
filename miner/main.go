package main

import (
	Modules "github.com/stevewooo/testin/Modules"
	Transaction "github.com/stevewooo/testin/Modules/Transaction"
	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
	MinerModules "github.com/stevewooo/testin/miner/Modules"
)

func main() {
	var trans Modules.ConsensusStorageData
	trans = new(Transaction.Transaction)
	trans.LoadFromJSONString("{}")

	config := map[string]string{
		"bcagName": "test",
		// "sdkRpcServer": "http://192.168.10.45:9024",
		"sdkRpcServer": "http://127.0.0.1:9024",
		"privateKey":   "8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6",
	}

	config["nodeID"] = Sign.GetPublicKey(config["privateKey"])[0:34]
	miner := MinerModules.Miner{}
	miner.Build(config)

	go miner.Run()

	c := make(chan bool, 1)
	<-c
}
