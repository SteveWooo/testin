package modules

import "github.com/stevewooo/testin/Modules/Transaction"

type PBFT_Preprepare struct {
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

// 从远程返回回来的Interface中，构建出一个preprepare包
func (p *PBFT_Preprepare) LoadFromInterface(obj interface{}) {

}

// 签名
func (p *PBFT_Preprepare) Sign(privateKey string) {

}
