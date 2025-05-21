package main

import (
	"feather/config"
	"feather/network"
	"feather/service"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", "localhost:8080", "port set")

func main() {
	flag.Parse()

	config.NewConfig(*pathFlag)
	n := network.NewServer(service.NewService(), *port)
	if err := n.StartServer(); err != nil {
		panic(err)
	}
}
