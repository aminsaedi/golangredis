package pkg

import (
	"fmt"
	"net"
	"time"

	"math/rand"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/thoas/go-funk"
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
		// fmt.Printf("Propogated: %v\tTotal: %v\n", len(propogated), len(config.PropogationStatus.Commands))
		if len(propogated) == len(config.PropogationStatus.Commands) {
			// append the replicaId to config.AppConfig.FullyPropogatedReplicaIds if not already present
			if !funk.ContainsString(config.AppConfig.FullyPropogatedReplicaIds, slaveId) {
				// fmt.Println("Fully propogated to slave: ", slaveId)
				config.AppConfig.FullyPropogatedReplicaIds = append(config.AppConfig.FullyPropogatedReplicaIds, slaveId)
			}
		} else {
			// filter the replicaId from config.AppConfig.FullyPropogatedReplicaIds
			config.AppConfig.FullyPropogatedReplicaIds = funk.FilterString(config.AppConfig.FullyPropogatedReplicaIds, func(replicaId string) bool {
				return replicaId != slaveId
			})
		}
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
