package modules

// TODO 抽象选主流程
type PoBR struct {
	config    map[string]string
	Term      string
	BlockList []string // 上一个时间限制内打包失败的打包者，在这一轮会被拉黑
}

// 开始一次投票-排序-打包操作
func (pobr *PoBR) DoPack() {

}
