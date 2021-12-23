package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	ConsensusModules "github.com/stevewooo/testin/Modules/Consensus"
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

	// HardCode for test
	Miners []string

	// 共识参数
	Repuations map[string]float64

	// 关联状态
	TaskHackers []*Transaction.TaskHacker

	// 操作对象暂不做存储
}

func (w *WorldStatus) Build(config map[string]string) {
	w.config = config

	// 正常情况下的矿工应该是专家列表。这里hardcode矿工是为了初始化启动网络
	w.Miners = []string{"047204499d849948aaffdec7ce2703f5b3",
		"0492ec813ab9ce7c94e49c84abcb0c7d64",
		"04c52654247aa39be86b5ce356ac7e24f8",
		"043abf9b64da3cf82a6833d827a6a60cb1",
		"0433cd50fa5977da115025e90cf5698c08"}
}

// 以对象形式返回一个世界状态
func (w *WorldStatus) GetWorldStatus() map[string]interface{} {
	return map[string]interface{}{
		"Hackers":     w.Hackers,
		"Enterprises": w.Enterprises,
		"Experts":     w.Experts,
		"Tasks":       w.Tasks,
		"Repuations":  w.Repuations,
	}
}

func (w *WorldStatus) cleanStatus() {
	w.Hackers = []*Transaction.Hacker{}
	w.Enterprises = []*Transaction.Enterprise{}
	w.Experts = []*Transaction.Expert{}
	w.Tasks = []*Transaction.Task{}
	w.TaskHackers = []*Transaction.TaskHacker{}
	w.Repuations = map[string]float64{}
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

	if trans.Type == "ApplyTaskByHacker" {
		w.TaskHackers = append(w.TaskHackers, &trans.TaskHacker)
	}

	// operation类型的交易
	if trans.Type == "AuthorizationHackerToTaskByEnterprise" {
		for i := 0; i < len(w.TaskHackers); i++ {
			if w.TaskHackers[i].TaskID == trans.AuthorizationHackerToTaskByEnterprise.TaskID &&
				w.TaskHackers[i].From == trans.AuthorizationHackerToTaskByEnterprise.HackerID {
				w.TaskHackers[i].IsPermission = "true"
				w.TaskHackers[i].ExpertList = []string{"047204499d849948aaffdec7ce2703f5b3"} // hard code
				w.TaskHackers[i].PermissionInformation = trans.AuthorizationHackerToTaskByEnterprise.PermissionInformation
			}
		}
	}

	if trans.Type == "PublishReportByHacker" {
		for i := 0; i < len(w.TaskHackers); i++ {
			if w.TaskHackers[i].TaskID == trans.TaskHackerReport.TaskID &&
				w.TaskHackers[i].From == trans.TaskHackerReport.From {
				w.TaskHackers[i].ReportPath = trans.TaskHackerReport.ReportPath
				// 每次提交报告都要清空专家评审意见
				w.TaskHackers[i].ExpertReviewReports = []*Transaction.ExpertReviewReport{}
			}
		}
	}

	if trans.Type == "ReviewReportByExpert" {
		for i := 0; i < len(w.TaskHackers); i++ {
			if w.TaskHackers[i].TaskID == trans.ExpertReviewReport.TaskID &&
				w.TaskHackers[i].From == trans.ExpertReviewReport.HackerID {
				w.TaskHackers[i].ExpertReviewReports = append(w.TaskHackers[i].ExpertReviewReports, &trans.ExpertReviewReport)
			}
		}

		// 专家review一次，会对信誉值进行加分
		repuationDataExists := false
		for nodeID, value := range w.Repuations {
			if trans.ExpertReviewReport.From == nodeID {
				repuationDataExists = true
				w.Repuations[nodeID] = value + 1.0 // TODO，目前每次加一分
				break
			}
		}
		if repuationDataExists == false {
			w.Repuations[trans.ExpertReviewReport.From] = 1.0
		}
	}
}

// 按照现有区块，构建世界状态
func (w *WorldStatus) DoBuildStatus() {
	w.cleanStatus()
	for i := 0; i < len(w.blocks); i++ {
		// 每新增一个区块，就要清空区块打包者的Repuation
		for nid, _ := range w.Repuations {
			if w.blocks[i].Miner == nid {
				w.Repuations[nid] = 0
			}
		}
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

// 拉取节点上的缓存打包意向信息
func (w *WorldStatus) GetRemotePackageIntentionCache(term string) ([]*ConsensusModules.PackageIntention, error) {
	localTopBlock := w.GetLocalTopBlock()
	if localTopBlock == nil {
		panic("远程区块不存在，本地缓存区块获取失败")
	}

	localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	newBlockNumber := strconv.Itoa(localTopBlockNumber + 1)
	prefix := "packageIntentionCache-" + newBlockNumber + "-" + term + "-"
	cacheResp, err := sdkApi.GetCacheByPerfix(w.config, prefix)
	if err != nil {
		return nil, errors.New("获取远程打包意向缓存失败")
	}
	if cacheResp.Status != 2000 {
		return nil, errors.New("Sdk报错：" + cacheResp.Message)
	}
	pis := []*ConsensusModules.PackageIntention{}
	for i := 0; i < len(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})); i++ {
		pi := ConsensusModules.PackageIntention{}
		err = pi.LoadFromJSONString(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
		if err != nil {
			fmt.Println(err)
		}

		pis = append(pis, &pi)
	}

	return pis, nil
}

// 拉取节点上的缓存打包意向排行榜
func (w *WorldStatus) GetRemoteIntentionRankCache(term string) ([]*ConsensusModules.IntentionRank, error) {
	localTopBlock := w.GetLocalTopBlock()
	if localTopBlock == nil {
		panic("远程区块不存在，本地缓存区块获取失败")
	}

	localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	newBlockNumber := strconv.Itoa(localTopBlockNumber + 1)
	prefix := "packageIntentionRankCache-" + newBlockNumber + "-" + term + "-"
	cacheResp, err := sdkApi.GetCacheByPerfix(w.config, prefix)
	if err != nil {
		return nil, errors.New("获取远程打包排行榜缓存失败")
	}
	if cacheResp.Status != 2000 {
		return nil, errors.New("Sdk报错：" + cacheResp.Message)
	}
	pis := []*ConsensusModules.IntentionRank{}
	for i := 0; i < len(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})); i++ {
		pi := ConsensusModules.IntentionRank{}
		err = pi.LoadFromJSONString(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
		if err != nil {
			fmt.Println(err)
		}

		pis = append(pis, &pi)
	}

	return pis, nil
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
