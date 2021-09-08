package main

import (
	"flag"
	"github.com/AdeMQ/client/wincmd"
)

func main() {
	// remote.address flag.String("address", "127.0.0.1:10601", "远程服务端地址")
	flag.Parse()

	cli := wincmd.NewWinClient()
	cli.Run()
}
