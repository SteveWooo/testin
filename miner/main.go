package main

import (
	Modules "github.com/stevewooo/testin/Modules"
	Transaction "github.com/stevewooo/testin/Modules/Transaction"
	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
	argvs "github.com/stevewooo/testin/Modules/Utils/argvs"
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
		// "privateKey":   "8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6",
	}

	// 通过启动参数输入privateKey:
	var argv argvs.Argv
	argv.Load()
	if argv.Get("privateKey") == "" {
		panic("privateKey不能为空")
	}
	config["privateKey"] = argv.Get("privateKey")

	config["nodeID"] = Sign.GetPublicKey(config["privateKey"])[0:34]
	miner := MinerModules.Miner{}
	miner.Build(config)

	// go miner.Run()
	go miner.RunProofOfBussinessReputation()

	c := make(chan bool, 1)
	<-c
}
