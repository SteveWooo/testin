package modules

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	ConsensusModules "github.com/stevewooo/testin/Modules/Consensus"
	"github.com/stevewooo/testin/Modules/SdkApi"
)

// 发布本节点的打包意向
func (m *Miner) PoBR_DoSendIntention(Term int) error {
	packageIntention := ConsensusModules.PackageIntention{}
	packageIntention.From = m.config["nodeID"]
	packageIntention.Intention = "true"
	packageIntention.Term = strconv.Itoa(Term)
	packageIntention.Ts = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	packageIntention.Sign(m.config["privateKey"])

	callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
		"MC_Call":          "DoPackageIntention",
		"PackageIntention": packageIntention,
	})
	if err != nil {
		return err
	}
	if callDoPackResp.Status != 2000 {
		return errors.New("sdk报错：" + callDoPackResp.Message)
	}

	return nil
}

// 获取节点上的所有Intention信息，计算排名，然后向节点发布本地计算出来的所有第一名的信息
func (m *Miner) PoBR_SendRepuationRank(Term int) error {
	pis, err := m.WorldStatus.GetRemotePackageIntentionCache(strconv.Itoa(Term))
	if err != nil {
		return err
	}

	// 参与联盟共识的节点必须是miner列表之中的节点才行
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
	minCandicateCount := m.PoBR_GetMinCandicateCount()
	if ledgePisCount < int(minCandicateCount) {
		// 规定次数内，没达到2/3意向者的话，就重新进入投票环节（说明出现平票情况，处理速度无限快的情况下，出现的概率为 pow(1/3, n)
		return errors.New("获取到的intention信息数据不足2/3")
	}

	// 把本节点获取到的数据进行排行
	// 从worldStatus中查询每个参与者的声誉分数
	ranker := []map[string]interface{}{} // 用于记录每个nid的声誉
	for i := 0; i < len(pis); i++ {
		// 不参与的就不要拿进来排序
		if pis[i].Intention != "true" {
			continue
		}
		for nid, score := range m.WorldStatus.Reputations {
			if pis[i].From == nid {
				ranker = append(ranker, map[string]interface{}{
					"nodeID":    nid,
					"repuation": score,
				})
			}
		}
	}
	// 对节点进行声誉排序
	for i := 0; i < len(ranker); i++ {
		for k := i + 1; k < len(ranker); k++ {
			if ranker[i]["repuation"].(float64) < ranker[k]["repuation"].(float64) {
				temp := ranker[i]
				ranker[i] = ranker[k]
				ranker[k] = temp
			}
		}
	}

	// 取出所有第一名(因为第一名会出现并列情况)
	maxRepuationMiner := []string{}
	for i := 0; i < len(ranker); i++ {
		if ranker[i]["repuation"].(float64) == ranker[0]["repuation"].(float64) {
			maxRepuationMiner = append(maxRepuationMiner, ranker[i]["nodeID"].(string))
			continue
		}
		break
	}

	// 发表排名信息，调用cvm脚本
	intentionRank := ConsensusModules.IntentionRank{}
	intentionRank.Rank_1 = strings.Join(maxRepuationMiner, ",")
	intentionRank.From = m.config["nodeID"]
	intentionRank.Term = strconv.Itoa(Term)
	intentionRank.Ts = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	intentionRank.Sign(m.config["privateKey"])

	callShareRankResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
		"MC_Call":       "ShareIntentionRank",
		"IntentionRank": intentionRank,
	})
	if err != nil {
		return errors.New("请求sdk出错")
	}
	if callShareRankResp.Status != 2000 {
		return errors.New("sdk报错：" + callShareRankResp.Message)
	}
	return nil
}

// 从节点上获取最新的排行缓存数据。计算出排名最高矿工
func (m *Miner) PoBR_GetMaxReputationMiner(Term int) (string, error) {
	pir, err := m.WorldStatus.GetRemoteIntentionRankCache(strconv.Itoa(Term))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 计算IntentionRank数据是否达到2/3参与者
	ledgePirCount := 0
	for i := 0; i < len(pir); i++ {
		for k := 0; k < len(m.WorldStatus.Miners); k++ {
			if pir[i].From == m.WorldStatus.Miners[k] {
				ledgePirCount = ledgePirCount + 1
				break
			}
		}
	}
	minIntentionRankCandicateCount := m.PoBR_GetMinCandicateCount()
	if ledgePirCount < int(minIntentionRankCandicateCount) {
		return "", errors.New("有效共识参与者信誉排行数据不足")
	}

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
		// 规定次数内，没达到2/3意向者的话，就重新进入投票环节（说明出现平票情况，处理速度无限快的情况下，出现的概率为 pow(1/3, n)
		return "", errors.New("最大认可数未达到2/3")
	}

	rank1MaxMiners := []string{} // 获得最大次数的Miner，即可参与最终排名
	for nid, count := range rank1List {
		if count == rank1MaxCount {
			rank1MaxMiners = append(rank1MaxMiners, nid)
		}
	}

	// 对平票者进行排序🌟，用上一个区块的哈希作为
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

			// 废置，这样不公平
			// nodeIDNumForI = nodeIDNumForI - hashRangeStartNum
			// nodeIDNumForK = nodeIDNumForK - hashRangeStartNum
			// // "ffffffff" - 自己，回退一个环
			// if nodeIDNumForI < 0 {
			// 	nodeIDNumForI = 4294967295 + nodeIDNumForI
			// }
			// if nodeIDNumForK < 0 {
			// 	nodeIDNumForK = 4294967295 + nodeIDNumForK
			// }

			// 越大越靠前
			if nodeIDNumForI < nodeIDNumForK {
				temp := rank1MaxMiners[i]
				rank1MaxMiners[i] = rank1MaxMiners[k]
				rank1MaxMiners[k] = temp
			}
		}
	}
	// 哈希环起点求余即可
	randomMinerIndex := int(hashRangeStartNum) % len(rank1MaxMiners)

	return rank1MaxMiners[randomMinerIndex], nil
}

// PoBR算法的打包流程：持续获取节点上的交易缓存，有的话，就打包下来。
func (m *Miner) PoBR_DoPackage(Term int) error {
	for {
		// 打包到成功为止，一般就是节点没交易缓存，就打包不了
		time.Sleep(time.Second * 2)
		packageErr := m.doPackage(strconv.Itoa(Term))
		if packageErr == nil {
			break
		}
		fmt.Println(packageErr)
	}

	return nil
}

// 持续向节点获取最新区块，检查是否已经完成打包
func (m *Miner) PoBR_CheckIsFinishedPackage() error {
	localTopBlock := m.WorldStatus.GetLocalTopBlock() // 先获取打包前的本地最高区块
	for {
		time.Sleep(1 * time.Second)
		remoteTopBlock := m.WorldStatus.GetRemoteTopBlock()
		localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
		remoteTopBlockNumber, _ := strconv.Atoi(remoteTopBlock.Number)
		if remoteTopBlockNumber > localTopBlockNumber {
			fmt.Println("完成打包，区块编号: " + remoteTopBlock.Number)

			// 直接拉最新worldStatus即可
			m.WorldStatus.AddNewLocalBlock(remoteTopBlock)
			// m.WorldStatus.DoBuildStatus()
			m.WorldStatus.FetchWorldStatus()
			break
		}

		if remoteTopBlockNumber == localTopBlockNumber {
			fmt.Println("远程节点共识打包中...")
			continue
		}

		panic("远程节点区块编号比本地低，请确认节点配置")
	}

	return nil
}

// 获取最少需要达到的投票数量
func (m *Miner) PoBR_GetMinCandicateCount() float64 {
	return math.Floor(float64(len(m.WorldStatus.Miners))*2/3) + 1
}
