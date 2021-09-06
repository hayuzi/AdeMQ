package main

import (
	"flag"
	"github.com/AdeMQ/client/wincmd"
)

func main() {
	flag.Parse()
	cli := wincmd.NewWinClient()
	cli.Run()
}
