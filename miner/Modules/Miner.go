package modules

import (
	"fmt"
	"strconv"
	"time"

	Modules "github.com/stevewooo/testin/Modules"
	"github.com/stevewooo/testin/Modules/SdkApi"
)

type Miner struct {
	config      map[string]string
	WorldStatus *Modules.WorldStatus
}

func (m *Miner) Build(config map[string]string) {
	m.config = config
	// 构造世界状态
	worldStatus := Modules.WorldStatus{}
	worldStatus.Build(config)
	worldStatus.FetchAllBlocks()
	worldStatus.DoBuildStatus()

	m.WorldStatus = &worldStatus
}

func (m *Miner) Run() {
	for {
		time.Sleep(time.Second * 3)
		// 判断自己是否打包节点
		if m.config["nodeID"] != "047204499d849948aaffdec7ce2703f5b3" {
			continue
		}

		// 1 检查最新区块是否和本地一致
		remoteTopBlock := m.WorldStatus.GetRemoteTopBlock()
		remoteBlockNumber, _ := strconv.Atoi(remoteTopBlock.Number)
		localTopBlock := m.WorldStatus.GetLocalTopBlock()
		localBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
		// 本地区块为节点上共识的区块时，拉取缓存交易，尝试进行打包
		if localBlockNumber == remoteBlockNumber {
			trans, err := m.WorldStatus.GetRemoteTransactionCache()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// TODO检查交易，排除掉无效交易。但这个检查其实在cvm中就有，矿工不一定必须提交该检查

			if len(trans) == 0 {
				continue
			}
			fmt.Println("packaing")
			// 进行打包
			newBlockNumber := remoteBlockNumber + 1

			block := Modules.Block{}
			block.PreviousHash = remoteTopBlock.Hash
			block.Number = strconv.Itoa(newBlockNumber)
			block.Miner = m.config["nodeID"]
			block.Ts = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			block.Transactions = trans
			block.Sign(m.config["privateKey"])

			callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
				"MC_Call": "DoPackage",
				"Block":   block,
			})
			if err != nil {
				fmt.Println("请求sdk出错")
				continue
			}
			if callDoPackResp.Status != 2000 {
				fmt.Println("sdk报错：" + callDoPackResp.Message)
				continue
			}

			// 把自己打包好的区块写入本地缓存
			m.WorldStatus.AddNewLocalBlock(&block)
			m.WorldStatus.DoBuildStatus()

			continue
		}

		if localBlockNumber < remoteBlockNumber {
			// TODO动态规划优化
			m.WorldStatus.FetchAllBlocks()
			m.WorldStatus.DoBuildStatus()
			continue
		}

	}
}
