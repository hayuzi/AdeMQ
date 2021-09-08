package commands

import (
	"context"
	"fmt"
	"strings"
)

func DescHistory() string {
	return `
history:
    命令介绍:    获取历史命令记录
    命令格式:    history
    命令参数:    无`
}

func History(ctx context.Context, params ...string) interface{} {
	history := ctx.Value(ConstHistory)
	cmdHis, ok := history.([]interface{})
	if !ok {
		return "Error: history error"
	}
	ret := ""
	for idx, val := range cmdHis {
		ret += fmt.Sprintf("%d    %s\n", idx, val)
	}
	return strings.TrimSuffix(ret, "\n")
}
