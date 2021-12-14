package feService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	Modules "github.com/stevewooo/testin/Modules"
	Transaction "github.com/stevewooo/testin/Modules/Transaction"
	Sign "github.com/stevewooo/testin/Modules/Utils/Sign"
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
	http.HandleFunc("/api/enterprise/get_my_task", feService.GetEnterprisePublishedTask)
	http.HandleFunc("/api/common/get_task_detail", feService.GetTaskDetail)
	http.HandleFunc("/api/common/get_task", feService.GetTask)

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

// 统一Get检查签名入口
// @params.queryParams
func checkGetSign(queryParams url.Values) bool {
	if len(queryParams["signature"]) == 0 ||
		len(queryParams["salt"]) == 0 ||
		len(queryParams["node_id"]) == 0 {
		return false
	}
	if queryParams["signature"][0] == "" ||
		queryParams["salt"][0] == "" ||
		queryParams["node_id"][0] == "" {
		return false
	}
	recoverPK, err := Sign.Recover(queryParams["signature"][0], queryParams["salt"][0])
	if err != nil {
		return false
	}
	if recoverPK[0:34] != queryParams["node_id"][0] {
		return false
	}
	return true
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

func (feService *FeService) GetEnterprisePublishedTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	resp := HttpResponser{}
	resp.Build()

	// 校验签名权限
	queryParams := req.URL.Query()
	if checkGetSign(queryParams) == false {
		resp.Status = 4003
		resp.Message = "参数错误或检查签名失败，接口无权访问"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	nodeID := queryParams["node_id"][0]

	if len(queryParams["page"]) == 0 {
		resp.Status = 4003
		resp.Message = "参数错误: page"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}
	page, err := strconv.Atoi(queryParams["page"][0])
	if err != nil {
		resp.Status = 4003
		resp.Message = "参数错误:" + err.Error()
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	if len(queryParams["item_per_page"]) == 0 {
		resp.Status = 4003
		resp.Message = "参数错误: page"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}
	itemPerPage, err := strconv.Atoi(queryParams["item_per_page"][0])
	if err != nil {
		resp.Status = 4003
		resp.Message = "参数错误:" + err.Error()
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// TODO 优化
	feService.WorldStatus.FetchAllBlocks()
	feService.WorldStatus.DoBuildStatus()

	// 筛选需要的内容
	status := feService.WorldStatus.GetWorldStatus()
	statusTask := status["Tasks"].([]*Transaction.Task)

	tasks := []*Transaction.Task{}
	for i := len(statusTask) - 1; i >= 0; i-- {
		if statusTask[i].From == nodeID {
			tasks = append(tasks, statusTask[i])
		}
	}

	startIndex := (page - 1) * itemPerPage
	endIndex := startIndex + itemPerPage
	if endIndex > len(tasks) {
		endIndex = len(tasks)
	}
	if startIndex > endIndex {
		resp.Status = 4004
		resp.Message = "超页错误"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	resp.Data = map[string]interface{}{
		"Count": len(tasks),
		"Tasks": tasks[startIndex:endIndex],
	}

	res.Write([]byte(resp.PackToJSONString()))
}

func (feService *FeService) GetTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	resp := HttpResponser{}
	resp.Build()

	// 校验签名权限
	queryParams := req.URL.Query()
	if checkGetSign(queryParams) == false {
		resp.Status = 4003
		resp.Message = "参数错误或检查签名失败，接口无权访问"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	if len(queryParams["page"]) == 0 {
		resp.Status = 4003
		resp.Message = "参数错误: page"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}
	page, err := strconv.Atoi(queryParams["page"][0])
	if err != nil {
		resp.Status = 4003
		resp.Message = "参数错误:" + err.Error()
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	if len(queryParams["item_per_page"]) == 0 {
		resp.Status = 4003
		resp.Message = "参数错误: page"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}
	itemPerPage, err := strconv.Atoi(queryParams["item_per_page"][0])
	if err != nil {
		resp.Status = 4003
		resp.Message = "参数错误:" + err.Error()
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// TODO 优化
	feService.WorldStatus.FetchAllBlocks()
	feService.WorldStatus.DoBuildStatus()

	// 筛选需要的内容
	status := feService.WorldStatus.GetWorldStatus()
	statusTask := status["Tasks"].([]*Transaction.Task)

	tasks := []*Transaction.Task{}
	for i := len(statusTask) - 1; i >= 0; i-- {
		tasks = append(tasks, statusTask[i])
	}

	startIndex := (page - 1) * itemPerPage
	endIndex := startIndex + itemPerPage
	if endIndex > len(tasks) {
		endIndex = len(tasks)
	}
	if startIndex > endIndex {
		resp.Status = 4004
		resp.Message = "超页错误"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	resp.Data = map[string]interface{}{
		"Count": len(tasks),
		"Tasks": tasks[startIndex:endIndex],
	}

	res.Write([]byte(resp.PackToJSONString()))
}

func (feService *FeService) GetTaskDetail(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	resp := HttpResponser{}
	resp.Build()

	// 校验签名权限
	queryParams := req.URL.Query()
	if checkGetSign(queryParams) == false {
		resp.Status = 4003
		resp.Message = "参数错误或检查签名失败，接口无权访问"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// 获取用户的信息和权限
	// nodeID := queryParams["node_id"][0]
	if len(queryParams["task_id"]) == 0 || len(queryParams["task_id"][0]) != 64 {
		resp.Status = 4003
		resp.Message = "参数错误: task_id"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// TODO 优化
	feService.WorldStatus.FetchAllBlocks()
	feService.WorldStatus.DoBuildStatus()

	// 查找任务
	var task *Transaction.Task = nil
	for i := 0; i < len(feService.WorldStatus.Tasks); i++ {
		if feService.WorldStatus.Tasks[i].Hash == queryParams["task_id"][0] {
			task = feService.WorldStatus.Tasks[i]
			break
		}
	}
	if task == nil {
		resp.Status = 4004
		resp.Message = "任务查询失败"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// 查询公司
	var taskCreaterEnterprise *Transaction.Enterprise = nil
	for i := 0; i < len(feService.WorldStatus.Enterprises); i++ {
		if feService.WorldStatus.Enterprises[i].From == task.From {
			taskCreaterEnterprise = feService.WorldStatus.Enterprises[i]
		}
	}
	if taskCreaterEnterprise == nil {
		resp.Status = 5000
		resp.Message = "系统查询失败：任务创建企业信息查询失败"
		res.Write([]byte(resp.PackToJSONString()))
		return
	}

	// 查询权限（任务发起人、任务测试员、任务评审人）
	isCreater := false
	isJoinHacker := false
	isPremissionHacker := false
	isTaskExpert := false

	if task.From == queryParams["node_id"][0] {
		isCreater = true
	}

	if task.TaskHackers != nil {
		for i := 0; i < len(task.TaskHackers); i++ {
			if task.TaskHackers[i].HackerID == queryParams["node_id"][0] {
				isJoinHacker = true
				if task.TaskHackers[i].IsPermission == "true" {
					isPremissionHacker = true
				}
				break
			}
		}

		for i := 0; i < len(task.TaskHackers); i++ {
			if task.TaskHackers[i].IsPermission == "false" {
				continue
			}

			for k := 0; k < len(task.TaskHackers[i].ExpertList); k++ {
				if task.TaskHackers[i].ExpertList[k] == queryParams["node_id"][0] {
					isTaskExpert = true
					break
				}
			}
		}
	}

	resp.Data = map[string]interface{}{
		"Task": task,
		"Permission": map[string]interface{}{
			"IsCreater":          isCreater,
			"IsJoinHacker":       isJoinHacker,
			"IsPremissionHacker": isPremissionHacker,
			"IsTaskExpert":       isTaskExpert,
		},
		"TaskCreaterEnterprise": taskCreaterEnterprise,
	}

	res.Write([]byte(resp.PackToJSONString()))
}
