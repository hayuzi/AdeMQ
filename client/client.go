package main

import "github.com/AdeMQ/client/wincmd"

func main() {
	cli := wincmd.NewWinClient()
	cli.Run()
}
