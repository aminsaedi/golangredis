package internal

import (
	"fmt"
	"os"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

func ReadRdbFile() {
	path := c.AppConfig.Dir + "/" + c.AppConfig.Dbfilename
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fmt.Println("Reading RDB file: ", path)
	fmt.Println(data[0])
}
