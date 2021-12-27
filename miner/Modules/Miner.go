package modules

import (
	"errors"
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

// 定期发布打包意向->
// 拉取足够的打包意向，业务声誉排序->
// 发布最高可打包者的nodeID->
// 拉取最高
func (m *Miner) RunProofOfBussinessReputation() {
	time.Sleep(time.Second * 1)
	// 用来控制轮数
	Term := 1
	for {
		time.Sleep(time.Second * 1)

		// 1 发布打包意向
		err := m.PoBR_DoSendIntention(Term)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 2 获取意向列表，进行本地Repuation排行，发布计算出来的第一名声誉节点到远程节点上
		for {
			time.Sleep(time.Second * 2)
			err := m.PoBR_SendRepuationRank(Term)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// 3 获取节点第一名声誉列表，统计获取第一名最多的节点
			// 如果自己是第一名，那就打包
			// 不是的话，就直接进入监听最新区块事件
			for {
				time.Sleep(time.Second * 2)
				// 有可能因为获取到的缓存数据不到2/3，导致计算失败
				rank1MaxMiner, err := m.PoBR_GetMaxReputationMiner(Term)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// 4 如果NodeID是自己，就打包
				// 如果NodeID不是自己，就监听最新区块信息。如果得到最新区块的话，
				if rank1MaxMiner == m.config["nodeID"] {
					fmt.Println("我要打包了:" + m.config["nodeID"])
					m.PoBR_DoPackage(Term)
				} else {
					fmt.Println("我不需要打包:" + m.config["nodeID"])
				}

				// 5 持续拉最新区块，当最新区块更新后，重新开始流程
				m.PoBR_CheckIsFinishedPackage()
				break
			}
			break
		}
	}
}

// 拉取交易，生成打包区块
// @params Term 选主轮次
func (m *Miner) doPackage(Term string) error {
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
			return err
		}

		// TODO检查交易，排除掉无效交易。但这个检查其实在cvm中就有，矿工不一定必须提交该检查

		if len(trans) == 0 {
			return errors.New("远程节点无缓存交易")
		}

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
			"Term":    Term,
		})
		if err != nil {
			return errors.New("请求sdk出错:" + err.Error())
		}
		if callDoPackResp.Status != 2000 {
			return errors.New("sdk报错：" + callDoPackResp.Message)
		}

		// 把自己打包好的区块写入本地缓存
		// m.WorldStatus.AddNewLocalBlock(&block)
		// m.WorldStatus.DoBuildStatus()

		return nil // 返回空代表打包发布完成
	}

	if localBlockNumber < remoteBlockNumber {
		// TODO动态规划优化
		m.WorldStatus.FetchAllBlocks()
		m.WorldStatus.DoBuildStatus()
		return errors.New("本地区块尚未与节点区块同步")
	}

	panic("本地区块比远程区块还高，请确认是否已经完成配置")
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
