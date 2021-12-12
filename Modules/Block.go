package modules

import (
	"encoding/json"

	"github.com/stevewooo/testin/Modules/Transaction"
	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
)

type Block struct {
	Hash string

	PreviousHash string
	Number       string
	MerkleRoot   string
	Miner        string
	Ts           string

	Transactions []*Transaction.Transaction
	Signature    string
}

func (b *Block) Sign(privateKey string) error {
	merkleSource := ""
	for i := 0; i < len(b.Transactions); i++ {
		merkleSource = merkleSource + b.Transactions[i].Hash
	}
	b.MerkleRoot = Sign.Hash(merkleSource)

	source := "Block" + b.MerkleRoot + b.Miner + b.Number + b.PreviousHash + b.Ts
	b.Hash = Sign.Hash(source)
	var err error
	b.Signature, err = Sign.Sign(b.Hash, privateKey)
	return err
}

func (b *Block) LoadFromJSONString(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), b)
	if err != nil {
		return err
	}
	return nil
}

func (b *Block) LoadFromInterface(obj interface{}) error {
	jsonByte, _ := json.Marshal(obj)
	b.LoadFromJSONString(string(jsonByte))
	return nil
}

func (b *Block) PackToJSONString() string {
	jsonByte, _ := json.Marshal(b)
	return string(jsonByte)
}
