package pkg

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

func connectToMaster() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// read `empty.rdb` file and send it to the master

	commands := []string{
		internal.ToArray("PING"),
		internal.ToArray("REPLCONF", "listening-port", fmt.Sprint(config.AppConfig.BindingPort)),
		internal.ToArray("REPLCONF", "capa", "psync2"),
		internal.ToArray("PSYNC", "?", "-1"),
	}
	for _, command := range commands {
		conn.Write([]byte(command))
		reply := make([]byte, 1024)
		conn.Read(reply)
	}

}

func PropogateToSlaves(commands []string) {
	for _, slaveAddress := range config.AppConfig.ConnectedSlaves {
		conn, err := net.Dial("tcp", slaveAddress)
		if err != nil {
			fmt.Println("Failed to connect to slave: ", slaveAddress)
			continue
		}
		defer conn.Close()

		for _, command := range commands {
			conn.Write([]byte(command))
			reply := make([]byte, 1024)
			conn.Read(reply)
		}
	}
}