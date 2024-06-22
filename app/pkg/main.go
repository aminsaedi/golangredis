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
	i.ReadRdbFile()

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

func getCommand(tokens []string) string {
	if len(tokens) > 0 && len(tokens[0]) > 1 && strings.HasPrefix(tokens[0], "*") {
		requiredItems, _ := strconv.Atoi(tokens[0][1:])
		requiredItems = requiredItems*2 + 1
		if len(tokens) == requiredItems {

			return strings.ToUpper(tokens[2])
		}
	}
	return ""
}

func sanitizeText(text string, tokens []string) (sanitizedText string, sanitizedTokens []string) {
	sanitizedTokens = tokens
	sanitizedText = strings.Map(func(r rune) rune {
		if r < 32 || r > 126 { // non-printable characters
			return -1
		}
		return r
	}, text)
	r := regexp.MustCompile(`.{5,}\*([0-9]+)$`)
	matches := r.FindStringSubmatch(text)
	if len(matches) > 0 {
		// fmt.Println("Matched", text, matches[1])
		sanitizedTokens = make([]string, 0)
		sanitizedText = "*" + matches[1]
	}
	return sanitizedText, sanitizedTokens
}

func updateTransferedBytes(tokens []string) {
	totalTokenLength := len(strings.Join(tokens, "")) + (len(tokens) * 2)
	c.PropogationStatus.TransferedBytes += totalTokenLength
}

func HandleRequestAsMaster(conn net.Conn, shouldWriteResult bool) {

	// defer conn.Close()

	scanner := bufio.NewScanner(conn)
	isConnectionFromSlave := false

	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()

		text, tokens = sanitizeText(text, tokens)

		// ignore binary data

		tokens = append(tokens, text)

		// pp.Print(tokens)

		command := getCommand(tokens)
		if command == "" {
			continue
		}

		var result string
		args := tokens[3:]
		switch command {
		case "ECHO":
			result = i.Echo(args...)
		case "PING":
			result = i.Ping()
		case "SET":
			result = i.Set(args...)
		case "GET":
			result = i.Get(args...)
		case "INFO":
			result = i.Info(args...)
		case "REPLCONF":
			result = i.Replconf(conn, args...)
		case "PSYNC":
			result = i.Psync(args...)
			conn.Write([]byte(result))
			result = i.RDBFileToString("empty.rdb")
			isConnectionFromSlave = true
		case "WAIT":
			result = i.Wait(args...)
		case "CONFIG":
			result = i.Config(args...)
		case "KEYS":
			result = i.Keys(args...)
		case "TYPE":
			result = i.Type(args...)
		}

		if shouldWriteResult || (command == "REPLCONF" && args[1] == "GETACK") {
			conn.Write([]byte(result))
		}

		updateTransferedBytes(tokens)
		tokens = make([]string, 0)

		if isConnectionFromSlave {
			go PropogateToSlaves(conn)
			break
		}

	}

}
