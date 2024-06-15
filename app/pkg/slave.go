package pkg

import (
	"fmt"
	"net"
	"time"

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
	for {
		time.Sleep(100 * time.Millisecond)
		for _, command := range config.PropogationStatus.Commands {
			key := command + "__" + conn.RemoteAddr().String()

			if _, ok := propogated[key]; ok {
				continue
			}

			conn.Write([]byte(command))
			conn.Write([]byte(internal.ToArray("REPLCONF", "GETACK", "*")))

			propogated[key] = struct{}{}
		}
		if len(propogated) == len(config.PropogationStatus.Commands) {
			// append the replicaId to config.AppConfig.FullyPropogatedReplicaIds if not already present
			if !funk.ContainsString(config.AppConfig.FullyPropogatedReplicaIds, config.AppConfig.MasterReplId) {
				config.AppConfig.FullyPropogatedReplicaIds = append(config.AppConfig.FullyPropogatedReplicaIds, config.AppConfig.MasterReplId)
			}
		} else {
			// filter the replicaId from config.AppConfig.FullyPropogatedReplicaIds
			config.AppConfig.FullyPropogatedReplicaIds = funk.FilterString(config.AppConfig.FullyPropogatedReplicaIds, func(replicaId string) bool {
				return replicaId != config.AppConfig.MasterReplId
			})

		}
	}
}
