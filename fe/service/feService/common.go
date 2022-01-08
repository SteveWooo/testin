package feService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (service *FeService) ProxyToWebIndex(res http.ResponseWriter, req *http.Request) {
	indexFile, _ := ioutil.ReadFile("./web/dist/index.html")
	res.Write(indexFile)
}

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
