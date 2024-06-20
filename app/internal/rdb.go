package internal

import (
	"os"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

func ReadRdbFile() {
	path := c.AppConfig.Dir + "/" + c.AppConfig.Dbfilename
	file, err := os.Open(path)
	if err != nil {
		return
	}
}
