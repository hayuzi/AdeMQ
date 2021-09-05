package wincmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// WinClient 交互式命令行客户端
type WinClient struct {
	Reader *bufio.Reader
}

// NewWinClient 创建命令行客户端
func NewWinClient() *WinClient {
	return &WinClient{
		Reader: bufio.NewReader(os.Stdin),
	}
}

// Run 命令行客户端启动运行
func (wc *WinClient) Run() {
	for {
		fmt.Print("$ ")
		cmdStr, err := wc.Reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		cmdStr = strings.TrimSuffix(cmdStr, "\n")
		if cmdStr != "" {
			_, _ = fmt.Fprintln(os.Stdout, cmdStr)
		}
	}
}

func handleCommand(cmdStr string) {
	// 处理命令行逻辑
}
