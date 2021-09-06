package handler

import (
	"context"
	"github.com/AdeMQ/client/handler/commands"
	"github.com/AdeMQ/client/remote"
)

type HandleFunc func(context context.Context, params ...string) interface{}

type Dispatcher struct {
	HelpInfo map[string]string
	History  []string
	Handlers map[string]HandleFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		HelpInfo: initHelpInfo(),
		History:  make([]string, 0),
		Handlers: initHandlers(),
	}
}

func (d *Dispatcher) Dispatch(cmd *ParsedCmd, remote *remote.Remote) interface{} {
	// 拦截空数据
	if cmd.Cmd == "" {
		return ""
	}
	// TODO 历史记录计入到数据结构中
	d.History = append(d.History, cmd.Cmd)

	// 正式执行函数
	fn, ok := d.Handlers[cmd.Cmd]
	if !ok {
		return "命令不存在"
	}

	ctx := context.WithValue(context.Background(), "remote", remote)
	ctx = d.injectSpecialContext(ctx, cmd.Cmd)
	return fn(ctx, cmd.Params...)
}

// 注入特殊命令需要的特殊上下文信息
func (d *Dispatcher) injectSpecialContext(ctx context.Context, cmd string) context.Context {
	switch cmd {
	case commands.ConstHelp:
		ctx = context.WithValue(ctx, commands.ConstHelp, d.HelpInfo)
	case commands.ConstHistory:
		ctx = context.WithValue(ctx, commands.ConstHistory, d.History)
	}
	return ctx
}
