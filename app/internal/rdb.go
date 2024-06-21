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
	byteData, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if !reflect.DeepEqual(byteData[:len(rdbMagicNumber)], rdbMagicNumber) {
		return
	}

	metaData := make([][]byte, 0)

	for i := len(rdbMagicNumber); i < len(byteData); i++ {
		if byteData[i] == 0xfa {
			ending := slices.Index(byteData[i+1:], 0xfa)
			if ending == -1 {
				ending = slices.Index(byteData[i+1:], 0xfe)
			}
			metaData = append(metaData, byteData[i:i+ending])
			i += ending
		}
	}

	dataBaseSection := make([][]byte, 0)

	for i := len(rdbMagicNumber); i < len(byteData); i++ {
		if byteData[i] == 0xfe {
			ending := slices.Index(byteData[i+1:], 0xfe)
			if ending == -1 {
				ending = slices.Index(byteData[i+1:], 0xff)
			}
			dataBaseSection = append(dataBaseSection, byteData[i:i+ending])
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

	keyValues := make(map[string]string)

	for i := 5; i < len(dataBaseSection[0]); i++ {
		keyStartIndex := slices.Index(dataBaseSection[0][i:], 0x0a)
		keyEndIndex := slices.Index(dataBaseSection[0][i+keyStartIndex+1:], 0x09)
		if keyStartIndex == -1 || keyEndIndex == -1 {
			break
		}
		key := dataBaseSection[0][i+keyStartIndex+1 : i+keyStartIndex+1+keyEndIndex]

		valueStartIndex := slices.Index(dataBaseSection[0][keyEndIndex:], 0x09)

		value := dataBaseSection[0][valueStartIndex+1+keyEndIndex:]

		keyValues[string(key)] = string(value)

		i += keyStartIndex + keyEndIndex + 1
	}

	fmt.Println(keyValues)

}
