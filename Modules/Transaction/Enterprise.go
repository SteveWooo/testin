package Transaction

import "encoding/json"

// 测试员对象
type Enterprise struct {
	Connection string
	From       string
	LogoPath   string
	Name       string
	SocialCode string
	Ts         string

	Hash      string
	Signature string
}

func (obj *Enterprise) LoadFromInterface(i interface{}) {
	jsonByte, _ := json.Marshal(i)
	_ = json.Unmarshal(jsonByte, obj)
}
