package Transaction

import "encoding/json"

// 测试员对象
type Hacker struct {
	From          string
	Name          string
	Resume        string
	Qualification string
	Ts            string

	Hash      string
	Signature string
}

func (obj *Hacker) LoadFromInterface(i interface{}) {
	jsonByte, _ := json.Marshal(i)
	_ = json.Unmarshal(jsonByte, obj)
}
