package feService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FeService struct {
	config map[string]string
}

func (feService *FeService) Build() {
	feService.config = map[string]string{}
	feService.config["httpPort"] = "10001"
	feService.config["sdkRpcServer"] = "http://127.0.0.1:9024"
	feService.config["privateKey"] = "8e1e5e540a07954e07a840d89eeed064b58ec16346b118ca6ad25831211f2ad6"
}

func (feService *FeService) Run() {

	http.HandleFunc("/api/hacker/register", feService.ApiHackerRegister)

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

func (feService *FeService) ApiHackerRegister(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	req.ParseForm()
	decoder := json.NewDecoder(req.Body)
	var reqParams map[string]interface{}
	decoder.Decode(&reqParams)

	fmt.Println(reqParams)

	// 设置共识脚本调用的函数
	reqParams["MC_Call"] = "RegisterHacker"

	submitParams, _ := json.Marshal(reqParams)
	submitBody := map[string]interface{}{
		"Params":   string(submitParams),
		"BcagName": "test",
	}
	bodyJSON, _ := json.Marshal(submitBody)
	reader := bytes.NewReader(bodyJSON)

	// 向SDK发送请求
	request, err := http.NewRequest("POST", feService.config["sdkRpcServer"]+"/sdk/consensus/call_trans", reader)
	if err != nil {
		fmt.Fprintf(res, "200")
		return
	}
	defer request.Body.Close()
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(res, "200")
		return
	}
	_, err = ioutil.ReadAll(resp.Body)

	fmt.Fprintf(res, "200")
}
