package cmd

import (
	"flag"

	"github.com/codecrafters-io/redis-starter-go/app/pkg"
)

func Execute() {
	var port int
	flag.IntVar(&port, "port", 6379, "Port to listen on")
	flag.Parse()

	startConfig := pkg.StartConfig{
		Port: port,
	}
	pkg.StartServer(startConfig)
}
