package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataItem struct {
	value        string
	validityInMS int64
	createdOn    time.Time `go:"init=now"`
}

func (di *DataItem) isValid() bool {
	return di.createdOn.UnixMilli()+di.validityInMS < time.Now().UnixMilli()
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

func getArayElement[T any](arr []T, index int, defaultValue T) T {
	if index >= 0 && index < len(arr) {
		return arr[index]
	}
	return defaultValue

}

func toBulkString(input string) string {
	return "$" + fmt.Sprint(len(input)) + "\r\n" + input + "\r\n"
}

func toSimpleString(input string) string {
	return "+" + input + "\r\n"
}

func toSimpleError(input string) string {
	return "$" + input + "\r\n"
}

func echo(value string) string {
	return toBulkString(value)
}

func ping() string {
	return toSimpleString("PONG")
}

type SetConfig struct {
	key        string
	value      string
	expiryType string
	expiryIn   string
}

func set(config SetConfig) string {

	toSet := DataItem{
		value: config.value,
	}
	if len(config.expiryType) == 2 {
		ms, _ := strconv.Atoi(config.expiryIn)
		toSet.validityInMS = int64(ms)
		fmt.Println("Setting expiry to ", ms)

	}
	storage[config.key] = toSet
	return toSimpleString("OK")
}

func get(key string) string {
	item, ok := storage[key]
	if ok && item.isValid() {
		return toBulkString(item.value)
	}
	return toSimpleError("-1")

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
					result = echo(tokens[4])
				case "PING":
					result = ping()
				case "SET":
					result = set(SetConfig{
						key:        tokens[4],
						value:      tokens[6],
						expiryType: getArayElement(tokens, 8, ""),
						expiryIn:   getArayElement(tokens, 10, ""),
					})
				case "GET":
					result = get(tokens[4])
				}

				conn.Write([]byte(result))

				// reset tokens
				tokens = make([]string, 0)
			}
		}

	}
	fmt.Println(tokens)
}
