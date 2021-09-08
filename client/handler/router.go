package handler

import (
	"github.com/AdeMQ/client/handler/commands"
)

func initHelpInfo() map[string]string {
	// TODO 新命令都需要注入进来
	// 注入帮助信息（请按照字典顺序处理）
	cmdHelp := make(map[string]string)
	cmdHelp[commands.ConstHelp] = commands.DescHelp()
	cmdHelp[commands.ConstHistory] = commands.DescHistory()
	cmdHelp[commands.ConstPing] = commands.DescPing()
	return cmdHelp
}

func initHandlers() map[string]HandleFunc {
	// 所有新增的数据结构要通过此处注入进来（请按照字典顺序处理）
	cmdDict := make(map[string]HandleFunc)
	cmdDict[commands.ConstHelp] = commands.Help
	cmdDict[commands.ConstHistory] = commands.History
	cmdDict[commands.ConstPing] = commands.Ping
	return cmdDict
}
