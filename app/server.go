package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

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

func echo() string {
	return "hey"
}

func ping() string {
	return "+PONG\r\n"
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()
		tokens = append(tokens, text)

		fmt.Println("In loop", tokens)

		if len(tokens) < 2 {
			continue
		}
		var result string
		command := strings.ToUpper(tokens[1])
		switch command {
		case "ECHO":
			if len(tokens) == 4 {
				fmt.Println("Calling echo func", tokens)
				result = echo()
				tokens = make([]string, 0)
			}
		case "PING":
			result = ping()
			tokens = make([]string, 0)
		}
		conn.Write([]byte(result))
	}
	fmt.Println(tokens)
}
