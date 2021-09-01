package main

import (
	"flag"
	"github.com/AdeMQ/conf"
	"github.com/AdeMQ/server/service"
	"log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Println("Hello AdeMQ")
	log.Println("TCP listen address ", conf.Conf.Server.Address)

	// 启动服务
	_ = service.Run(conf.Conf.Server)

}
