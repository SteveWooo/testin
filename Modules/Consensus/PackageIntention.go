package consensus

import (
	"encoding/json"

	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
)

// 打包意向数据包
type PackageIntention struct {
	From      string
	Intention string
	Term      string
	Ts        string

	Hash      string
	Signature string
}

func (pi *PackageIntention) Sign(privateKey string) error {
	source := "DoPackageIntention" + pi.From + pi.Intention + pi.Term + pi.Ts
	pi.Hash = Sign.Hash(source)
	var err error
	pi.Signature, err = Sign.Sign(pi.Hash, privateKey)
	return err
}

func (pi *PackageIntention) LoadFromJSONString(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), pi)
	if err != nil {
		return err
	}
	return nil
}

func (pi *PackageIntention) LoadFromInterface(obj interface{}) error {
	jsonByte, _ := json.Marshal(obj)
	pi.LoadFromJSONString(string(jsonByte))
	return nil
}

func (pi *PackageIntention) PackToJSONString() string {
	jsonByte, _ := json.Marshal(pi)
	return string(jsonByte)
}
