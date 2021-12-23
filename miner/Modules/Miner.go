package modules

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	Modules "github.com/stevewooo/testin/Modules"
	ConsensusModules "github.com/stevewooo/testin/Modules/Consensus"
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

// 定期发布打包意向->拉取足够的打包意向，业务声誉排序->发布最高可打包者的nodeID
// ->拉取最高
func (m *Miner) RunProofOfBussinessReputation() {
	time.Sleep(time.Second * 1)
	for {
		time.Sleep(time.Second * 3)

		// 1 发布打包意向
		packageIntention := ConsensusModules.PackageIntention{}
		packageIntention.From = m.config["nodeID"]
		packageIntention.Intention = "true"
		packageIntention.Term = "1"
		packageIntention.Ts = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		packageIntention.Sign(m.config["privateKey"])

		callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
			"MC_Call":          "DoPackageIntention",
			"PackageIntention": packageIntention,
		})
		if err != nil {
			fmt.Println("请求sdk出错")
			continue
		}
		if callDoPackResp.Status != 2000 {
			fmt.Println("sdk报错：" + callDoPackResp.Message)
			continue
		}

		// 2 监听所有人发布的打包意向
		for {
			time.Sleep(time.Second * 1)
			pis, err := m.WorldStatus.GetRemotePackageIntentionCache(packageIntention.Term)
			if err != nil {
				fmt.Println(err)
				continue
			}

			ledgePisCount := 0
			for i := 0; i < len(pis); i++ {
				for k := 0; k < len(m.WorldStatus.Miners); k++ {
					if pis[i].From == m.WorldStatus.Miners[k] {
						ledgePisCount = ledgePisCount + 1
						break
					}
				}
			}

			// 计算看是否达到2/3参与者
			minCandicateCount := math.Floor(float64(len(m.WorldStatus.Miners)) * 2 / 3)
			if ledgePisCount < int(minCandicateCount) {
				// 规定次数内，没达到2/3意向者的话，就重新进入投票环节（说明出现平票情况，处理速度无限快的情况下，出现的概率为 pow(1/3, n)
				continue
			}

			// 3 进入排行步骤
			// 从worldStatus中查询每个参与者的声誉分数
			ranker := []map[string]interface{}{} // 用于暂时排序
			for i := 0; i < len(pis); i++ {
				// 不参与的就不要拿进来排序
				if pis[i].Intention != "true" {
					continue
				}
				for nid, score := range m.WorldStatus.Repuations {
					if pis[i].From == nid {
						ranker = append(ranker, map[string]interface{}{
							"nodeID":    nid,
							"repuation": score,
						})
					}
				}
			}
			// 对ranker进行排序
			for i := 0; i < len(ranker); i++ {
				for k := i + 1; k < len(ranker); k++ {
					if ranker[i]["repuation"].(float64) < ranker[k]["repuation"].(float64) {
						temp := ranker[i]
						ranker[i] = ranker[k]
						ranker[k] = temp
					}
				}
			}

			// 取出所有第一名
			maxRepuationMiner := []string{}
			for i := 0; i < len(ranker); i++ {
				if ranker[i]["repuation"].(float64) == ranker[0]["repuation"].(float64) {
					maxRepuationMiner = append(maxRepuationMiner, ranker[i]["nodeID"].(string))
					continue
				}
				break
			}

			// 4 发表排名信息，调用cvm脚本
			intentionRank := ConsensusModules.IntentionRank{}
			intentionRank.Rank_1 = strings.Join(maxRepuationMiner, ",")
			intentionRank.From = m.config["nodeID"]
			intentionRank.Term = "1"
			intentionRank.Ts = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			intentionRank.Sign(m.config["privateKey"])

			callShareRankResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
				"MC_Call":       "ShareIntentionRank",
				"IntentionRank": intentionRank,
			})
			if err != nil {
				fmt.Println("请求sdk出错")
				continue
			}
			if callShareRankResp.Status != 2000 {
				fmt.Println("sdk报错：" + callShareRankResp.Message)
				continue
			}

			// 5 监听所有人的排名情况
			for {
				time.Sleep(time.Second * 1)
				pir, err := m.WorldStatus.GetRemoteIntentionRankCache(packageIntention.Term)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// 计算看是否达到2/3参与者
				ledgePirCount := 0
				for i := 0; i < len(pir); i++ {
					for k := 0; k < len(m.WorldStatus.Miners); k++ {
						if pir[i].From == m.WorldStatus.Miners[k] {
							ledgePirCount = ledgePirCount + 1
							break
						}
					}
				}
				// 计算看是否达到2/3参与者
				minIntentionRankCandicateCount := math.Floor(float64(len(m.WorldStatus.Miners)) * 2 / 3)
				if ledgePirCount < int(minIntentionRankCandicateCount) {
					// 规定次数内，没达到2/3意向者的话，就重新进入投票环节（说明出现平票情况，处理速度无限快的情况下，出现的概率为 pow(1/3, n)
					continue
				}

				// 6 获取排名第一的票数最多的nodeID
				rank1List := map[string]int{} // NodeID获得第一名的次数
				for i := 0; i < len(pir); i++ {
					topRanks := strings.Split(pir[i].Rank_1, ",")
					for r := 0; r < len(topRanks); r++ {
						if rank1List[topRanks[r]] == 0 {
							rank1List[topRanks[r]] = 0
						}
						rank1List[topRanks[r]] = rank1List[topRanks[r]] + 1
					}
				}
				rank1MaxCount := 0 // 获得第一名的最大次数
				for _, count := range rank1List {
					if count > rank1MaxCount {
						rank1MaxCount = count
					}
				}
				// 最大的认可次数未达到2/3总数
				if float64(rank1MaxCount) < minIntentionRankCandicateCount {
					// 多次未达到，说明平票了
					continue
				}

				rank1MaxMiners := []string{} // 获得最大次数的Miner，即可参与最终排名
				for nid, count := range rank1List {
					if count == rank1MaxCount {
						rank1MaxMiners = append(rank1MaxMiners, nid)
					}
				}

				// 对平票者进行排序🌟，这里后续要设计一个排序算法，用上一个区块的哈希作为
				// 哈希环的起点，然后对两者进行比较
				topBlock := m.WorldStatus.GetLocalTopBlock()
				if topBlock == nil {
					panic("本地最高区块获取失败")
				}
				hashRangeStart := topBlock.Hash[0:8] // 哈希环起点（区块哈希是64位的，取前32位即可。NodeID是34位的，取后32位即可
				hashRangeStartNum, _ := strconv.ParseUint(hashRangeStart[:8], 16, 32)
				for i := 0; i < len(rank1MaxMiners); i++ {
					for k := i + 1; k < len(rank1MaxMiners); k++ {
						nodeIDNumForI, _ := strconv.ParseUint(rank1MaxMiners[i][2:10], 16, 32)
						nodeIDNumForK, _ := strconv.ParseUint(rank1MaxMiners[k][2:10], 16, 32)
						nodeIDNumForI = nodeIDNumForI - hashRangeStartNum
						nodeIDNumForK = nodeIDNumForK - hashRangeStartNum

						// "ffffffff" - 自己，回退一个环
						if nodeIDNumForI < 0 {
							nodeIDNumForI = 4294967295 + nodeIDNumForI
						}
						if nodeIDNumForK < 0 {
							nodeIDNumForK = 4294967295 + nodeIDNumForK
						}

						// 越大越靠前
						if nodeIDNumForI < nodeIDNumForK {
							temp := rank1MaxMiners[i]
							rank1MaxMiners[i] = rank1MaxMiners[k]
							rank1MaxMiners[k] = temp
						}
					}
				}

				// 7 如果NodeID是自己，就打包
				// 7 如果NodeID不是自己，就监听最新区块信息。如果得到最新区块的话，
				if rank1MaxMiners[0] == m.config["nodeID"] {
					fmt.Println("我要打包了:" + m.config["nodeID"])
				} else {
					fmt.Println("我不需要打包:" + m.config["nodeID"])
				}

				// 8 持续拉最新区块，当最新区块更新后，重新开始流程

			}
		}
	}
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

			// fmt.Println(trans[0])

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
