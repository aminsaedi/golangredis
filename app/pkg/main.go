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
	i "github.com/codecrafters-io/redis-starter-go/app/internal"
)

type StartConfig struct {
	Port       int
	Replicaof  string
	Dir        string
	Dbfilename string
}

func StartServer(config StartConfig) {

	c.AppConfig.Replicaof = config.Replicaof
	c.AppConfig.BindingPort = config.Port
	c.AppConfig.Dir = config.Dir
	c.AppConfig.Dbfilename = config.Dbfilename

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

		go HandleRequestAsMaster(conn, true)
	}
}

func HandleRequestAsMaster(conn net.Conn, shouldSendResponse bool) {

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
			// fmt.Println("Matched", text, matches[1])
			tokens = make([]string, 0)
			text = "*" + matches[1]
		}

		tokens = append(tokens, text)

		// pp.Print(tokens)

		if len(tokens) > 0 && len(tokens[0]) > 1 && strings.HasPrefix(tokens[0], "*") {
			requiredItems, _ := strconv.Atoi(tokens[0][1:])
			requiredItems = requiredItems*2 + 1
			if len(tokens) == requiredItems {

				var result string
				command := strings.ToUpper(tokens[2])
				// fmt.Printf("Is slave: %v Command: %v Args: %v\n", !shouldSendResponse, command, tokens[3:])
				switch command {
				case "ECHO":
					result = i.Echo(tokens[3:]...)
				case "PING":
					result = i.Ping()
				case "SET":
					result = i.Set(tokens[3:]...)
				case "GET":
					result = i.Get(tokens[3:]...)
				case "INFO":
					result = i.Info(tokens[3:]...)
				case "REPLCONF":
					result = i.Replconf(tokens[3:]...)
				case "PSYNC":
					result = i.Psync(tokens[3:]...)
					conn.Write([]byte(result))
					result = i.RDBFileToString("empty.rdb")
					isConnectionFromSlave = true
				case "WAIT":
					result = i.Wait(tokens[3:]...)
				case "CONFIG":
					result = i.Config(tokens[3:]...)
				}

				totalTokenLength := len(strings.Join(tokens, "")) + (len(tokens) * 2)
				c.PropogationStatus.TransferedBytes += totalTokenLength
				if shouldSendResponse || command == "REPLCONF" {
					conn.Write([]byte(result))
				}

				// reset tokens
				tokens = make([]string, 0)

			}

		}
		if isConnectionFromSlave {
			go PropogateToSlaves(conn)
			break
		}

	}

}
