package internal

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"slices"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

var rdbMagicNumber = []byte{0x52, 0x45, 0x44, 0x49, 0x53}

func ReadRdbFile() {
	path := c.AppConfig.Dir + "/" + c.AppConfig.Dbfilename
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	if !reflect.DeepEqual(data[:len(rdbMagicNumber)], rdbMagicNumber) {
		return
	}

	metaData := make([][]byte, 0)

	for i := len(rdbMagicNumber); i < len(data); i++ {
		if data[i] == 0xfa {
			ending := slices.Index(data[i+1:], 0xfa)
			if ending == -1 {
				ending = slices.Index(data[i+1:], 0xfe)
			}
			metaData = append(metaData, data[i:i+ending])
			i += ending
		}
	}

	dataBaseSection := make([][]byte, 0)

	for i := len(rdbMagicNumber); i < len(data); i++ {
		if data[i] == 0xfe {
			ending := slices.Index(data[i+1:], 0xfe)
			if ending == -1 {
				ending = slices.Index(data[i+1:], 0xff)
			}
			dataBaseSection = append(dataBaseSection, data[i:i+ending])
			i += ending
		}
	}

	fmt.Println("Meta data")
	for _, meta := range metaData {
		fmt.Println(hex.EncodeToString(meta))
	}
	fmt.Println("Data base section")
	for _, data := range dataBaseSection {
		fmt.Println(hex.EncodeToString(data))
	}
}
