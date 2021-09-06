package commands

import (
	"context"
	"fmt"
)

func DescHistory() string {
	return `
history:
    命令格式:    history
    命令参数:    无`
}

func History(ctx context.Context, params ...string) interface{} {
	history := ctx.Value(ConstHistory)
	his, ok := history.([]string)
	if !ok {
		return "Error: history error"
	}
	ret := ""
	for idx, val := range his {
		ret += fmt.Sprintf("%d :   %s\n", idx, val)
	}
	return ret
}
