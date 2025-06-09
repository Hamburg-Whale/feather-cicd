package main

import (
	"feather/config"
	"feather/handler"
	"feather/repository"
	"feather/service"
	"flag"

	_ "github.com/go-sql-driver/mysql"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", "localhost:8080", "port set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)
	if r, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		n := handler.NewServer(service.NewService(r), *port)
		if err := n.StartServer(); err != nil {
			panic(err)
		}
	}
}
