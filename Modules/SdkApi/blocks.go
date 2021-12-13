package SdkApi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Response struct {
	Status  int
	Message string
	Data    interface{}
}

func (response *Response) BuildFromString(res string) error {
	err := json.Unmarshal([]byte(res), response)
	if err != nil {
		return err
	}
	return nil
}

func (response *Response) PackToJSONString() string {
	jsonByte, _ := json.Marshal(response)
	return string(jsonByte)
}

func doGet(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	response := Response{}
	err = response.BuildFromString(string(body))
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func doPost(url string, params map[string]interface{}) (*Response, error) {
	// 提交的参数主体
	bodyJSON, _ := json.Marshal(params)
	reader := bytes.NewReader(bodyJSON)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	// 设置header
	request.Header.Set("Content-Type", "application/json")

	// 创建一个socket并发送数据
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	cnfResponse := Response{}
	err = cnfResponse.BuildFromString(string(body))
	if err != nil {
		return nil, err
	}

	return &cnfResponse, nil
}

func GetTopBlock(config map[string]string) (*Response, error) {
	return doGet(config["sdkRpcServer"] + "/sdk/consensus/get_info_by_name?bcag_name=" + config["bcagName"])
}

func GetBlock(config map[string]string, number int) (*Response, error) {
	return doGet(config["sdkRpcServer"] + "/sdk/consensus/get_block_by_number?bcag_name=" + config["bcagName"] + "&block_number=" + strconv.Itoa(number))
}

func GetBlockByRange(config map[string]string, start string, end string) (*Response, error) {
	return doGet(config["sdkRpcServer"] + "/sdk/consensus/get_block_by_range?bcag_name=" + config["bcagName"] + "&start=" + start + "&end=" + end)
}

// @params.BcagName 算法名称
// @params.Params 参数JSON字符串
func CallTrans(config map[string]string, params map[string]interface{}) (*Response, error) {
	submitParams := map[string]interface{}{
		"BcagName": config["bcagName"],
	}
	paramsJSONByte, _ := json.Marshal(params)
	submitParams["Params"] = string(paramsJSONByte)

	return doPost(config["sdkRpcServer"]+"/sdk/consensus/call_trans", submitParams)
}
