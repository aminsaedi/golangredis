package pkg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/k0kubun/pp"
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

	scanner := bufio.NewScanner(conn)
	isConnectionFromSlave := false

	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()

		// ignore binary data
		text = strings.Map(func(r rune) rune {
			if r < 32 || r > 126 { // non-printable characters
				return -1
			}
			return r
		}, text)
		r := regexp.MustCompile(`.{5,}\*([0-9]+)$`)
		matches := r.FindStringSubmatch(text)
		if len(matches) > 0 {
			fmt.Println("Matched", text, matches[1])
			text = "*" + matches[1]
		}

		tokens = append(tokens, text)

		pp.Print(tokens)

		if len(tokens) > 0 && strings.HasPrefix(tokens[0], "*") {
			requiredItems, _ := strconv.Atoi(tokens[0][1:])
			requiredItems = requiredItems*2 + 1
			if len(tokens) == requiredItems {
				// run command

				var result string
				command := strings.ToUpper(tokens[2])
				switch command {
				case "ECHO":
					result = internal.Echo(tokens[4])
				case "PING":
					result = internal.Ping()
				case "SET":
					result = internal.Set(internal.SetConfig{
						Key:        tokens[4],
						Value:      tokens[6],
						ExpiryType: internal.GetArayElement(tokens, 8, ""),
						ExpiryIn:   internal.GetArayElement(tokens, 10, ""),
					})
				case "GET":
					result = internal.Get(tokens[4])
				case "INFO":
					result = internal.Info(tokens[3:]...)
				case "REPLCONF":
					result = internal.Replconf(tokens[3:]...)
				case "PSYNC":
					result = internal.Psync(tokens[3:]...)
					conn.Write([]byte(result))
					result = internal.RDBFileToString("empty.rdb")
					isConnectionFromSlave = true
					// print connected slave address and port
					fmt.Println("Connected slave: ", conn.RemoteAddr().String())
				}

				conn.Write([]byte(result))

				// reset tokens
				tokens = make([]string, 0)

			}

		}
		if isConnectionFromSlave {
			// conn.Write([]byte("amin"))
			go PropogateToSlaves(conn)
			fmt.Println("Breaking loop")
			break
		}

	}

}
