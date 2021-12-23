package consensus

import (
	"encoding/json"

	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
)

// 排名
type IntentionRank struct {
	From   string
	Rank_1 string // 排名第一的NodeID，用逗号分隔
	Term   string
	Ts     string

	Hash      string
	Signature string
}

func (pi *IntentionRank) Sign(privateKey string) error {
	source := "ShareIntentionRank" + pi.From + pi.Rank_1 + pi.Term + pi.Ts
	pi.Hash = Sign.Hash(source)
	var err error
	pi.Signature, err = Sign.Sign(pi.Hash, privateKey)
	return err
}

func (pi *IntentionRank) LoadFromJSONString(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), pi)
	if err != nil {
		return err
	}
	return nil
}

func (pi *IntentionRank) LoadFromInterface(obj interface{}) error {
	jsonByte, _ := json.Marshal(obj)
	pi.LoadFromJSONString(string(jsonByte))
	return nil
}

func (pi *IntentionRank) PackToJSONString() string {
	jsonByte, _ := json.Marshal(pi)
	return string(jsonByte)
}
