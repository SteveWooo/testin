package Transaction

import (
	"encoding/json"
)

// 存储在区块中的交易对象
type Transaction struct {
	Type  string
	Nonce int
	Hash  string
	// 对象表
	Hacker     Hacker
	Enterprise Enterprise
	Expert     Expert
	Task       Task
	// 关联对象表
	TaskHacker TaskHacker

	// 操作类型对象表
	AuthorizationHackerToTaskByEnterprise AuthorizationHackerToTaskByEnterprise
	TaskHackerReport                      TaskHackerReport
	ExpertReviewReport                    ExpertReviewReport
}

func (t *Transaction) DoInit() {
	// if t.Type == "RegisterHacker" {
	// 	t.Hash = t.Hacker.Hash
	// }
	// if t.Type == "RegisterEnterprise" {
	// 	t.Hash = t.Enterprise.Hash
	// }
	// if t.Type == "RegisterExpert" {
	// 	t.Hash = t.Expert.Hash
	// }
}

func (t *Transaction) LoadFromJSONString(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), t)
	t.DoInit()
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) LoadFromInterface(obj interface{}) error {
	jsonByte, _ := json.Marshal(obj)
	t.LoadFromJSONString(string(jsonByte))
	return nil
}
