package pkg

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"time"

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
		reply := make([]byte, 32)
		conn.Read(reply)
	}
	HandleRequestAsMaster(conn, false)

}

func PropogateToSlaves(conn net.Conn) {
	propogated := make(map[string]struct{})
	slaveId := fmt.Sprint(rand.Int())
	// fmt.Println("Propogating to slave: ", slaveId)
	for {

		go func() {
			buff := make([]byte, 64)
			conn.Read(buff)
			// if buff includes ACK command then log "ALLLLLLL"
			if regexp.MustCompile(`ACK`).Match(buff) {
				fmt.Println("ALLLLLLL")
				config.AppConfig.FullyPropogatedReplicaIds = append(config.AppConfig.FullyPropogatedReplicaIds, slaveId)
			}
		}()

		for _, command := range config.PropogationStatus.Commands {
			key := command + "__" + conn.RemoteAddr().String()

			if _, ok := propogated[key]; ok {
				continue
			}

			conn.Write([]byte(command))

			propogated[key] = struct{}{}
		}

		time.Sleep(50 * time.Millisecond)
	}
}
