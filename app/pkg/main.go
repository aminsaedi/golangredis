package pkg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

type StartConfig struct {
	Port      int
	Replicaof string
}

func StartServer(config StartConfig) {

	c.AppConfig.Replicaof = config.Replicaof
	c.AppConfig.BindingPort = config.Port

	if c.AppConfig.Replicaof != "" {
		go connectToMaster()
	}

	listenser, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.Port))
	if err != nil {
		fmt.Println("Failed to bind to port")
		os.Exit(1)
	}

	fmt.Println("Listening on port ", config.Port)

	for {
		conn, err := listenser.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleRequestAsMaster(conn)
	}
}

func HandleRequestAsMaster(conn net.Conn) {
	// defer conn.Close()
	for {
		scanner := bufio.NewScanner(conn)
		isConnectionFromSlave := false

		command, args := internal.ParseCommand(scanner)

		var result string
		switch command {
		case "ECHO":
			result = internal.Echo(args[0])
		case "PING":
			result = internal.Ping()
		case "SET":
			fmt.Println("SET command")
			// result = internal.Set(internal.SetConfig{
			// 	Key:        tokens[4],
			// 	Value:      tokens[6],
			// 	ExpiryType: internal.GetArayElement(tokens, 8, ""),
			// 	ExpiryIn:   internal.GetArayElement(tokens, 10, ""),
			// })
		case "GET":
			result = internal.Get(args[1])
		case "INFO":
			result = internal.Info(args[3:]...)
		case "REPLCONF":
			result = internal.Replconf(args[3:]...)
		case "PSYNC":
			result = internal.Psync(args[3:]...)
			conn.Write([]byte(result))
			result = internal.RDBFileToString("empty.rdb")
			isConnectionFromSlave = true
			// print connected slave address and port
			fmt.Println("Connected slave: ", conn.RemoteAddr().String())
		}

		conn.Write([]byte(result))

		if isConnectionFromSlave {
			// conn.Write([]byte("amin"))
			go PropogateToSlaves(conn)
			fmt.Println("Breaking loop")
			break
		}
	}

}
