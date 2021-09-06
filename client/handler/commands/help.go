package commands

import (
	"context"
	"fmt"
	"sort"
)

func DescHelp() string {
	return `
help:
    命令格式:    help [cmd]
    命令参数:    [cmd] 可选参数: 为你想要获取帮助信息的命令`
}

func Help(ctx context.Context, params ...string) interface{} {
	help := ctx.Value(ConstHelp)
	helpInfo, ok := help.(map[string]string)
	if !ok {
		return "Error: helpInfo error"
	}
	// 如果给定了第二个参数并符合某个命令，直接给单个命令的帮助信息
	if len(params) >= 1 {
		p1 := params[0]
		if cmd, ok := helpInfo[p1]; ok {
			return cmd + "\n"
		}
	}
	var names []string
	for name := range helpInfo {
		names = append(names, name)
	}
	sort.Strings(names)
	ret := ""
	for _, name := range names {
		ret += fmt.Sprintf("%s", helpInfo[name])

	}
	return ret + "\n"
}
