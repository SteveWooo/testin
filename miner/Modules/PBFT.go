package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/stevewooo/testin/Modules/SdkApi"
	"github.com/stevewooo/testin/Modules/Transaction"
	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
)

func (m *Miner) PBFT_GetMinCandicateCount() float64 {
	return math.Floor(float64(len(m.WorldStatus.Miners))*2/3) + 1
}

// 构建Preprepare包，并发送
func (m *Miner) PBFT_DoPreprePare() error {
	// 检查区块编号，拉交易
	remoteTopBlock := m.WorldStatus.GetRemoteTopBlock()
	remoteBlockNumber, _ := strconv.Atoi(remoteTopBlock.Number)
	localTopBlock := m.WorldStatus.GetLocalTopBlock()
	localBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	if localBlockNumber == remoteBlockNumber { // 区块合法，拉缓存交易
		trans, err := m.WorldStatus.GetRemoteTransactionCache()
		if err != nil {
			return err
		}

		if len(trans) == 0 {
			return errors.New("无交易")
		}

		// 构建，发送preprepare
		prepreParePack := PBFT_PrePrepare{}
		prepreParePack.SetTransactions(trans)
		prepreParePack.Number = strconv.Itoa(remoteBlockNumber + 1)
		prepreParePack.PreviousHash = remoteTopBlock.Hash
		prepreParePack.Miner = m.config["nodeID"]
		prepreParePack.MinerSign(m.config["privateKey"])

		callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
			"MC_Call":        "DoPBFTPreprepare",
			"PrepreParePack": prepreParePack,
		})
		if err != nil {
			return err
		}
		if callDoPackResp.Status != 2000 {
			return errors.New("sdk报错：" + callDoPackResp.Message)
		}

		return nil
	}

	return errors.New("区块不同步")
}

// 拉取远程缓存区块
func (m *Miner) PBFT_GetRemotePrePreparePackCache() ([]*PBFT_PrePrepare, error) {
	localTopBlock := m.WorldStatus.GetLocalTopBlock()
	if localTopBlock == nil {
		panic("远程区块不存在，本地缓存区块获取失败")
	}

	localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	newBlockNumber := strconv.Itoa(localTopBlockNumber + 1)
	prefix := "PBFTPrepreParePack-" + newBlockNumber + "-"
	cacheResp, err := SdkApi.GetCacheByPerfix(m.config, prefix)
	if err != nil {
		return nil, errors.New("获取远程打包排行榜缓存失败")
	}
	if cacheResp.Status != 2000 {
		return nil, errors.New("Sdk报错：" + cacheResp.Message)
	}

	prePreparePackCache := []*PBFT_PrePrepare{}
	for i := 0; i < len(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})); i++ {
		pc := PBFT_PrePrepare{}
		err = pc.LoadFromJSONString(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
		if err != nil {
			fmt.Println(err)
		}

		prePreparePackCache = append(prePreparePackCache, &pc)
	}

	return prePreparePackCache, nil
}

// preprepare包
type PBFT_PrePrepare struct {
	Hash         string
	PreviousHash string
	Number       string
	Miner        string // 这一轮的矿工
	MerkleRoot   string

	MinerSignature string // 矿工的签名

	Transactions []*Transaction.Transaction

	// From      string // 见证者
	// Signature string // 见证者的签名
}

// 载入交易，构建merkle根
func (p *PBFT_PrePrepare) SetTransactions(trans []*Transaction.Transaction) {
	p.Transactions = trans
	// 计算MerkleRoot
	source := ""
	for i := 0; i < len(p.Transactions); i++ {
		source = source + p.Transactions[i].Hash
	}
	p.MerkleRoot = Sign.Hash(source)
}

// 从远程返回回来的Interface中，构建出一个preprepare包
func (p *PBFT_PrePrepare) LoadFromJSONString(str string) error {
	err := json.Unmarshal([]byte(str), p)
	if err != nil {
		return err
	}
	return nil
}

// 签名
func (p *PBFT_PrePrepare) MinerSign(privateKey string) {
	source := p.MerkleRoot + p.Miner + p.Number + p.PreviousHash
	p.Hash = Sign.Hash(source)
	p.MinerSignature, _ = Sign.Sign(p.Hash, privateKey)
}

// 检查签名
func (p *PBFT_PrePrepare) CheckMinerSign() error {
	recoverNodeID, err := Sign.Recover(p.MinerSignature, p.Hash)
	if err != nil {
		return err
	}
	if recoverNodeID[0:34] != p.Miner {
		return errors.New("签名校验失败")
	}
	return nil
}

// Prepare阶段：
func (m *Miner) PBFT_DoPrepare(prePreparePack *PBFT_PrePrepare) error {
	// 构建，发送preprepare
	preparePack := PBFT_Prepare{}
	preparePack.BuildFromPreprePare(prePreparePack)
	preparePack.From = m.config["nodeID"]
	preparePack.Sign(m.config["privateKey"])

	callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
		"MC_Call":     "DoPBFTPrepare",
		"PreparePack": preparePack,
	})
	if err != nil {
		return err
	}
	if callDoPackResp.Status != 2000 {
		return errors.New("sdk报错：" + callDoPackResp.Message)
	}

	return nil
}

// 拉取远程缓存区块
func (m *Miner) PBFT_GetRemotePreparePackCache() ([]*PBFT_Prepare, error) {
	localTopBlock := m.WorldStatus.GetLocalTopBlock()
	if localTopBlock == nil {
		panic("远程区块不存在，本地缓存区块获取失败")
	}

	localTopBlockNumber, _ := strconv.Atoi(localTopBlock.Number)
	newBlockNumber := strconv.Itoa(localTopBlockNumber + 1)
	prefix := "PBFTPreparePack-" + newBlockNumber + "-"
	cacheResp, err := SdkApi.GetCacheByPerfix(m.config, prefix)
	if err != nil {
		return nil, errors.New("获取远程打包排行榜缓存失败")
	}
	if cacheResp.Status != 2000 {
		return nil, errors.New("Sdk报错：" + cacheResp.Message)
	}

	preparePackCache := []*PBFT_Prepare{}
	for i := 0; i < len(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})); i++ {
		pc := PBFT_Prepare{}
		err = pc.LoadFromJSONString(cacheResp.Data.(map[string]interface{})["Caches"].([]interface{})[i].(string))
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = pc.CheckSign()
		if err != nil {
			fmt.Println(err)
			continue
		}

		preparePackCache = append(preparePackCache, &pc)
	}

	return preparePackCache, nil
}

// prepare包
type PBFT_Prepare struct {
	Hash         string
	PreviousHash string
	Number       string
	Miner        string // 这一轮的矿工
	MerkleRoot   string

	MinerSignature string // 矿工的签名

	Transactions []*Transaction.Transaction

	From      string // 见证者
	Signature string // 见证者的签名
}

func (p *PBFT_Prepare) BuildFromPreprePare(preprePare *PBFT_PrePrepare) {
	p.Hash = preprePare.Hash
	p.PreviousHash = preprePare.PreviousHash
	p.Number = preprePare.Number
	p.Miner = preprePare.Miner
	p.MerkleRoot = preprePare.MerkleRoot
	p.Transactions = preprePare.Transactions
	p.MinerSignature = preprePare.MinerSignature
}

// 从远程返回回来的json中，构建出一个preprepare包
func (p *PBFT_Prepare) LoadFromJSONString(str string) error {
	err := json.Unmarshal([]byte(str), p)
	if err != nil {
		return err
	}
	return nil
}

// 签名
func (p *PBFT_Prepare) Sign(privateKey string) {
	source := p.MerkleRoot + p.Miner + p.Number + p.PreviousHash
	p.Hash = Sign.Hash(source)
	p.Signature, _ = Sign.Sign(p.Hash, privateKey)
}

// 检查签名
func (p *PBFT_Prepare) CheckMinerSign() error {
	recoverNodeID, err := Sign.Recover(p.MinerSignature, p.Hash)
	if err != nil {
		return err
	}
	if recoverNodeID[0:34] != p.Miner {
		return errors.New("签名校验失败")
	}
	return nil
}

func (p *PBFT_Prepare) CheckSign() error {
	recoverNodeID, err := Sign.Recover(p.Signature, p.Hash)
	if err != nil {
		return err
	}
	if recoverNodeID[0:34] != p.From {
		return errors.New("签名校验失败")
	}
	return nil
}

// Commit阶段

func (m *Miner) PBFT_DoCommit(prepare *PBFT_Prepare) error {
	commitPack := PBFT_Commit{}
	commitPack.BuildFromPrepare(prepare)
	commitPack.From = m.config["nodeID"]
	commitPack.Sign(m.config["privateKey"])

	callDoPackResp, err := SdkApi.CallTrans(m.config, map[string]interface{}{
		"MC_Call":    "DoPBFTCommit",
		"CommitPack": commitPack,
	})
	if err != nil {
		return err
	}
	if callDoPackResp.Status != 2000 {
		return errors.New("sdk报错：" + callDoPackResp.Message)
	}
	return nil
}

type PBFT_Commit struct {
	Hash         string
	PreviousHash string
	Number       string
	Miner        string // 这一轮的矿工
	MerkleRoot   string

	MinerSignature string // 矿工的签名

	Transactions []*Transaction.Transaction

	From      string // 见证者
	Signature string // 见证者的签名
}

func (p *PBFT_Commit) BuildFromPrepare(prepare *PBFT_Prepare) {
	p.Hash = prepare.Hash
	p.PreviousHash = prepare.PreviousHash
	p.Number = prepare.Number
	p.Miner = prepare.Miner
	p.MerkleRoot = prepare.MerkleRoot
	p.Transactions = prepare.Transactions
	p.MinerSignature = prepare.MinerSignature
}

// 从远程返回回来的json中，构建出一个preprepare包
func (p *PBFT_Commit) LoadFromJSONString(str string) error {
	err := json.Unmarshal([]byte(str), p)
	if err != nil {
		return err
	}
	return nil
}

// 签名
func (p *PBFT_Commit) Sign(privateKey string) {
	source := p.MerkleRoot + p.Miner + p.Number + p.PreviousHash
	p.Hash = Sign.Hash(source)
	p.Signature, _ = Sign.Sign(p.Hash, privateKey)
}

// 检查签名
func (p *PBFT_Commit) CheckMinerSign() error {
	recoverNodeID, err := Sign.Recover(p.MinerSignature, p.Hash)
	if err != nil {
		return err
	}
	if recoverNodeID[0:34] != p.Miner {
		return errors.New("签名校验失败")
	}
	return nil
}

func (p *PBFT_Commit) CheckSign() error {
	recoverNodeID, err := Sign.Recover(p.Signature, p.Hash)
	if err != nil {
		return err
	}
	if recoverNodeID[0:34] != p.From {
		return errors.New("签名校验失败")
	}
	return nil
}

// 持续向节点获取最新区块，检查是否已经完成打包
func (m *Miner) PBFT_CheckIsFinishedPackage() error {
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
