package Transaction

import "encoding/json"

// 存储在区块中的交易对象
type Transaction struct {
	Type   string
	Nonce  int
	Hash   string
	Hacker Hacker
}

// 对刚构造好的交易数据进行初始化处理，包括赋予哈希
func (t *Transaction) doInit() {
	if t.Type == "RegisterHacker" {
		t.Hash = t.Hacker.Hash
	}
}

func (t *Transaction) LoadFromJSONString(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), t)
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
