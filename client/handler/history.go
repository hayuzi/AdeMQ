package handler

import (
	"github.com/AdeMQ/datastruct/linear"
)

// CmdHistory 命令历史记录
// 我们使用了环形队列作为底层结构来实现历史记录的存存储
type CmdHistory struct {
	CmdList *linear.RingQueue
}

// NewCmdHistory 返回一个命令历史结构的指针
func NewCmdHistory() *CmdHistory {
	return &CmdHistory{
		CmdList: linear.NewRingQueue(100),
	}
}

func (c *CmdHistory) Push(cmd string) {
	if c.CmdList.IsFull() {
		_, _ = c.CmdList.DeQueue()
	}
	c.CmdList.EnQueue(cmd)
}

func (c *CmdHistory) All() []interface{} {
	return c.CmdList.FetchAllElem()
}
