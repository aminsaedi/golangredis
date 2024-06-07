package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
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

		buff := make([]byte, 50)
		c := bufio.NewReader(conn)

		for {
			size, err := c.Read(buff)
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
				break
			}
			fmt.Println("Size:", string(buff[:size]))

			_, err = io.ReadFull(c, buff[:size])
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
				break
			}
			fmt.Println("Data:", string(buff[:size]))
		}

	}

}
