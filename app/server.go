package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type DataItem struct {
	value  string
	expiry int
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

func toBulkString(input string) string {
	return "$" + fmt.Sprint(len(input)) + "\r\n" + input + "\r\n"
}

func toSimpleString(input string) string {
	return "+" + input + "\n\r"
}

func echo(value string) string {
	return toBulkString(value)
}

func ping() string {
	return toSimpleString("PONG")
}

func set(key, value string) string {
	storage[key] = DataItem{
		value:  value,
		expiry: -1,
	}
	return toSimpleString("OK")
}

func get(key string) string {
	item, ok := storage[key]
	if ok {
		return item.value
	} else {
		return toSimpleString("Error")
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()
		tokens = append(tokens, text)

		fmt.Println("In loop", tokens)

		if len(tokens) < 3 {
			continue
		}
		var result string
		command := strings.ToUpper(tokens[2])
		switch command {
		case "ECHO":
			if len(tokens) == 5 {
				fmt.Println("Calling echo func", tokens)
				result = echo(tokens[4])
				tokens = make([]string, 0)
			}
		case "PING":
			result = ping()
			tokens = make([]string, 0)
		case "SET":
			result = set(tokens[4], tokens[55])
		case "GET":
			result = get(tokens[4])
		}

		fmt.Printf("Rresult is : %q\n", result)
		conn.Write([]byte(result))
	}
	fmt.Println(tokens)
}
