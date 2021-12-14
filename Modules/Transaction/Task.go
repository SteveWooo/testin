package Transaction

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

// 测试员子对象
type TaskHacker struct {
	HackerID string // 关联到隔壁Hacker表
	TaskID   string // 任务的Hash值
	Ts       string // 最后一次上传报告的日期

	Hash      string
	Signature string

	IsPermission          string              // 是否已经对测试员进行授权
	PermissionInformation string              // 授权信息说明
	ExpertList            []string            // 参与这次评审的专家列表
	ReportPATH            string              // 上传报告的IPFS地址，默认空
	TaskExpertReports     []*TaskExpertReport // 专家报告
	Negotiations          []*TaskNegotiation  // 协商过程
}

// 专家报告子子对象
type TaskExpertReport struct {
	ExpertID string
	TaskID   string
	HackerID string
	Ts       string

	Hash      string
	Signature string

	Score string
	Memo  string
}

// 协商消息子子对象
type TaskNegotiation struct {
}
