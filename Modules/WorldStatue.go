package modules

import (
	"errors"
	"strconv"

	sdkApi "github.com/stevewooo/testin/Modules/SdkApi"
	"github.com/stevewooo/testin/Modules/Transaction"
)

// 链上数据的实时世界状态
type WorldStatus struct {
	config map[string]string

	blocks []*Block // 链上所有区块

	// 状态
	Hackers []*Transaction.Hacker
	// Experts
	// Enterprises
}

func (w *WorldStatus) Build(config map[string]string) {
	w.config = config
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
		tran.LoadFromJSONString(transResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
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

// 读取交易数据，写入世界状态。
func (w *WorldStatus) parseTransaction(trans *Transaction.Transaction) {
	if trans.Type == "RegisterHacker" {
		w.Hackers = append(w.Hackers, &trans.Hacker)
	}
}

func (w *WorldStatus) DoBuildStatus() {
	for i := 0; i < len(w.blocks); i++ {
		for k := 0; k < len(w.blocks[i].Transactions); k++ {
			w.parseTransaction(w.blocks[i].Transactions[k])
		}
	}
}
