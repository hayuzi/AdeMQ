package commands

import (
	"context"
	"github.com/AdeMQ/client/remote"
)

func DescPing() string {
	return `
ping:
    命令介绍:    向远程服务器发送连接消息
    命令格式:    ping
    命令参数:    无`
}

func Ping(ctx context.Context, params ...string) interface{} {
	// TODO 远程命令考虑整体封装，统一处理
	server := ctx.Value(ConstRemote)
	srv, ok := server.(*remote.Remote)
	if !ok {
		return "Error: remote error"
	}
	data, err := remote.FormatRequest(ConstPing, params)
	if err != nil {
		return "Error: cmd format error"
	}
	err = srv.SendMsgToRequestChan(data)
	if err != nil {
		return err
	}
	result, err := srv.GetResponseFromChan()
	if err != nil {
		return err.Error()
	}
	return string(result)
}
