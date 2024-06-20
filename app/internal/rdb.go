package internal

import (
	"encoding/hex"
	"os"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

func ReadRdbFile() {
	path := c.AppConfig.Dir + "/" + c.AppConfig.Dbfilename
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	d := hex.Dumper(os.Stdout)
	d.Write(data)
	d.Close()
}
