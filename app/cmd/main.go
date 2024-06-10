package cmd

import (
	"flag"
	"regexp"

	"github.com/codecrafters-io/redis-starter-go/app/pkg"
)

func Execute() {
	var port int
	var replicaof string
	flag.IntVar(&port, "port", 6379, "Port to listen on")
	flag.StringVar(&replicaof, "replicaof", "", "Replicate to another server")
	flag.Parse()

	// check replicaof pattern with regex
	if replicaof != "" {
		pattern := regexp.MustCompile(`^([a-zA-Z0-9]+) ([0-9]+)$`)
		if !pattern.MatchString(replicaof) {
			panic("Invalid replicaof pattern")
		}
	}

	startConfig := pkg.StartConfig{
		Port:      port,
		Replicaof: replicaof,
	}
	pkg.StartServer(startConfig)
}
