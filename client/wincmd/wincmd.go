package wincmd

import (
	"bufio"
	"fmt"
	"github.com/AdeMQ/client/handler"
	"github.com/AdeMQ/client/remote"
	"os"
	"strings"
)

// WinClient 交互式命令行客户端
type WinClient struct {
	Reader     *bufio.Reader
	Parser     *handler.Parser
	Dispatcher *handler.Dispatcher
	Remote     *remote.Remote
}

// NewWinClient 创建命令行客户端
func NewWinClient() *WinClient {
	return &WinClient{
		Reader:     bufio.NewReader(os.Stdin),
		Parser:     handler.NewParser(),
		Dispatcher: handler.NewDispatcher(),
		Remote:     remote.NewRemote(),
	}
}

// Run 命令行客户端启动运行
func (wc *WinClient) Run() {
	// 链接到服务端
	wc.Remote.Init()

	// 阻塞读取命令行数据
	for {
		fmt.Print("$ ")
		cmdStr, err := wc.Reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		cmdStr = strings.TrimSuffix(cmdStr, "\n")
		if cmdStr != "" {
			wc.handleCommand(cmdStr)
		}
	}
}

func (wc *WinClient) handleCommand(cmdStr string) {
	// 处理命令行逻辑
	cmd := wc.Parser.Parse(cmdStr)
	ret := wc.Dispatcher.Dispatch(cmd, wc.Remote)
	// 得到结果直接输出到标准输出
	if ret == "" {
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, ret)
}
