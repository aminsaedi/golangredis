package pkg

import (
	"fmt"
	"net"
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

	// go func() {
	// 	reg := regexp.MustCompile(`ACK`)
	// 	scanner := bufio.NewScanner(conn)
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		if reg.MatchString(line) {
	// 			// atomic.AddInt32(&count, 1)
	// 			fmt.Println("Ooon", conn.RemoteAddr().String())
	// 			config.Counter.Increment()
	// 			break
	// 		}
	// 	}
	// }()
	for {
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
