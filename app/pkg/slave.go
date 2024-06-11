package pkg

import (
	"fmt"
	"net"
	"strings"

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
		internal.RDBFileToString("empty.rdb"),
	}
	for _, command := range commands {
		conn.Write([]byte(command))
		if !strings.Contains(command, "PSYNC") {
			reply := make([]byte, 1024)
			conn.Read(reply)
			fmt.Println("Waiting for reply: ", string(reply))
		} else {
			fmt.Println("Skipping reply")
		}
	}
}
