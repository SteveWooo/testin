package Transaction

import "encoding/json"

// 测试员对象
type Task struct {
	Budget                string
	From                  string // 发布人的NodeID，不是企业id
	MaxAuthorizationCount string
	Name                  string
	Require               string
	Resume                string
	Ts                    string

	Hash      string
	Signature string

	IsPublic    string
	TaskHackers []*TaskHacker // 测试实例列表（测试员列表）
}

func (obj *Task) LoadFromInterface(i interface{}) {
	jsonByte, _ := json.Marshal(i)
	_ = json.Unmarshal(jsonByte, obj)
}

// 测试员子对象
type TaskHacker struct {
	From   string // 关联到隔壁Hacker表
	TaskID string // 任务的Hash值
	Ts     string // 最后一次上传报告的日期

	Hash      string
	Signature string

	IsPermission          string   // 是否已经对测试员进行授权
	PermissionInformation string   // 授权信息说明
	ExpertList            []string // 参与这次评审的专家列表
	ReportPath            string   // 上传报告的IPFS地址，默认空

	ExpertReviewReports []*ExpertReviewReport // 专家报告
	Negotiations        []*TaskNegotiation    // 协商过程

	Hacker *Hacker // 关联信息
}

func (obj *TaskHacker) LoadFromInterface(i interface{}) {
	jsonByte, _ := json.Marshal(i)
	_ = json.Unmarshal(jsonByte, obj)
}

// 专家报告子子对象
type ExpertReviewReport struct {
	From     string // 关联到专家表
	TaskID   string
	HackerID string
	Score    string // 评分 0到100
	Memo     string // 评语
	Ts       string

	Hash      string
	Signature string
}

// 协商消息子子对象
type TaskNegotiation struct {
}

// 操作类型交易对象
type AuthorizationHackerToTaskByEnterprise struct {
	From                  string // 企业id
	HackerID              string // 测试员ID
	PermissionInformation string // 授权信息
	TaskID                string // 任务id
	Ts                    string

	Hash      string
	Signature string
}

type TaskHackerReport struct {
	From       string // 测试员id
	TaskID     string
	ReportPath string
	Ts         string

	Hash      string
	Signature string
}
