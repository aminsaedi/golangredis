package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	internal "github.com/codecrafters-io/redis-starter-go/app/internal"
)

type DataItem struct {
	value     string
	validTill time.Time
}

func (di *DataItem) isValid() bool {
	if di.validTill.IsZero() {
		return true
	}
	return di.validTill.After(time.Now())
}

var storage = make(map[string]DataItem)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()
		tokens = append(tokens, text)

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
				}

				conn.Write([]byte(result))

				// reset tokens
				tokens = make([]string, 0)
			}
		}

	}
	fmt.Println(tokens)
}
