package pkg

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/k0kubun/pp"
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

func PropogateToSlaves() {
	// to prevent propogation of same command multiple times, we will keep track of commands that have been propogated
	propogated := make(map[string]struct{})
	for {
		time.Sleep(1 * time.Second)
		for _, conn := range config.PropogationStatus.ConnectedSlaves {
			// fmt.Println("Propogating to slave: ", conn.RemoteAddr().String())
			for index, command := range config.PropogationStatus.Commands {
				key := fmt.Sprint(index) + "__" + conn.RemoteAddr().String()
				// fmt.Println("Already propogated: ", propogated)

				if _, ok := propogated[key]; ok {

					continue
				}

				// fmt.Printf("Propogating command: %q\n", command)
				conn.Write([]byte(command))
				propogated[key] = struct{}{}
				pp.Println(propogated)
			}

		}
	}
}
