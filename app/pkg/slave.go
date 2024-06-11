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
	// send PING
	// send REPLCONF listening-port port
	// send REPLCONF ip-address ip
	// send SYNC
	// read RDB

	// send PING
	conn.Write([]byte(internal.ToArray("PING")))
	conn.Write([]byte(internal.ToArray("REPLCONF", "listening-port", fmt.Sprint(config.AppConfig.BindingPort))))
	conn.Write([]byte(internal.ToArray("REPLCONF", "capa", "psync2")))

}
