package feService

import "encoding/json"

type HttpResponser struct {
	Status  int
	Message string
	Data    map[string]interface{}
}

func (resp *HttpResponser) Build() {
	resp.Status = 2000
	resp.Message = "ok"
	resp.Data = map[string]interface{}{}
}

func (resp *HttpResponser) PackToJSONString() string {
	jsonByte, _ := json.Marshal(resp)
	return string(jsonByte)
}
