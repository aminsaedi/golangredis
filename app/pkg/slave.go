package pkg

import "net"

func connectToMaster() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// send PING
	// send REPLCONF listening-port port
	// send REPLCONF ip-address ip
	// send SYNC
	// read RDB

	// send PING
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))

}
