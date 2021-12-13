package feService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	Modules "github.com/stevewooo/testin/Modules"
)

type FeService struct {
	config      map[string]string
	WorldStatus *Modules.WorldStatus
}

func (feService *FeService) Build() {
	feService.config = map[string]string{}
	feService.config["httpPort"] = "10001"
	feService.config["bcagName"] = "test"
	feService.config["sdkRpcServer"] = "http://127.0.0.1:9024"

	worldStatus := Modules.WorldStatus{}
	worldStatus.Build(feService.config)
	worldStatus.FetchAllBlocks()
	worldStatus.DoBuildStatus()

	feService.WorldStatus = &worldStatus
}

func (feService *FeService) Run() {
	// 注册前端代理接口
	http.HandleFunc("/api/proxy", feService.Proxy)

	// 业务查询接口
	http.HandleFunc("/api/world_status/get", feService.GetWorldStatus)

	// 管理静态文件目录
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Http service listened at : " + feService.config["httpPort"])
	// 启动HTTP服务
	err := http.ListenAndServe("0.0.0.0:"+feService.config["httpPort"], nil)
	if err != nil {
		panic(err)
	}
}

func (feService *FeService) Proxy(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	// 处理参数
	req.ParseForm()
	decoder := json.NewDecoder(req.Body)
	var reqParams map[string]interface{}
	decoder.Decode(&reqParams)

	// 参数需要再次JSON格式化，因为传入cvm中的参数只能有一个，必须是个JSON字符串
	paramsJSONString, _ := json.Marshal(reqParams["Params"])
	submitBody := map[string]interface{}{
		"Params":   string(paramsJSONString),
		"BcagName": "test",
	}

	// 构造向sdk发送请求的requester
	bodyJSON, _ := json.Marshal(submitBody)
	reader := bytes.NewReader(bodyJSON)
	request, err := http.NewRequest("POST", feService.config["sdkRpcServer"]+"/sdk/consensus/call_trans", reader)
	if err != nil {
		res.Write([]byte("400: " + err.Error()))
		return
	}
	defer request.Body.Close()
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		res.Write([]byte("400: " + err.Error()))
		return
	}
	sdkRespBody, err := ioutil.ReadAll(resp.Body)

	res.Write([]byte(sdkRespBody))
}

func (feService *FeService) GetWorldStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	// TODO 优化
	feService.WorldStatus.FetchAllBlocks()
	feService.WorldStatus.DoBuildStatus()

	resp := HttpResponser{}
	resp.Build()
	resp.Data = feService.WorldStatus.GetWorldStatus()

	res.Write([]byte(resp.PackToJSONString()))
}
