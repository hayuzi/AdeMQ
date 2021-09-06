package handler

// CmdHistory 命令历史记录
// 采用队列这样一个数据结构，可以先进先出，淘汰达到上限的旧命令
//   理论上来说，命令是非定长的，TODO 我们可以选择单向或者双向链表的数据结构
//   但是此处为了简便，也可以选择使用 go语言的slice或者数组，但是我们得自己做一个环形队列的实现
type CmdHistory struct {
	CmdList []string
}

// NewCmdHistory 返回一个命令历史结构的指针
func NewCmdHistory(len int) *CmdHistory {
	return &CmdHistory{
		CmdList: make([]string, len),
	}
}

func (c *CmdHistory) push() {

}

func (c *CmdHistory) pop() {

}
