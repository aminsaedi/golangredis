package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"slices"
	"time"

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

	// metaData := make([][]byte, 0)

	// for i := len(rdbMagicNumber); i < len(byteData); i++ {
	// 	if byteData[i] == 0xfa {
	// 		ending := slices.Index(byteData[i+1:], 0xfa)
	// 		if ending == -1 {
	// 			ending = slices.Index(byteData[i+1:], 0xfe)
	// 		}
	// 		metaData = append(metaData, byteData[i:i+ending])
	// 		i += ending
	// 	}
	// }

	dataBaseSection := make([][]byte, 0)

	for i := len(rdbMagicNumber); i < len(byteData); i++ {
		if byteData[i] == 0xfe {
			ending := slices.Index(byteData[i+1:], 0xfe)
			if ending == -1 {
				ending = slices.Index(byteData[i+1:], 0xff)
			}
			dataBaseSection = append(dataBaseSection, byteData[i:i+ending+1])
			i += ending
		}
	}

	// fmt.Println("Meta data")
	// for _, meta := range metaData {
	// 	fmt.Println(hex.EncodeToString(meta))
	// }
	fmt.Println("Data base section")
	for _, data := range dataBaseSection {
		fmt.Println(hex.EncodeToString(data))
	}

	for index, data := range dataBaseSection[0] {
		fmt.Print(hex.EncodeToString([]byte{data}))
		fmt.Printf(" - %c - i:%d\n", data, index)
	}

	// totalSize := int(dataBaseSection[0][3])
	withExpiryItemSize := int(dataBaseSection[0][4])

	if withExpiryItemSize == 0 {
		i := 5
		for {

			keySize := int(dataBaseSection[0][i+1])
			fmt.Println("Key size", keySize, " Key index", i+1)
			valueSize := int(dataBaseSection[0][i+2+keySize])
			fmt.Println("Value size", valueSize, " Value index", i+2+keySize)

			key := dataBaseSection[0][i+2 : i+2+keySize]
			value := dataBaseSection[0][i+3+keySize : i+3+keySize+valueSize]

			fmt.Print("Key: ", string(key))
			fmt.Println(" -- Value:", string(value))

			SetStorageItem(string(key), DataItem{value: string(value)})

			// find index of next 0x00
			nextNull := slices.Index(dataBaseSection[0][i+3+keySize+valueSize:], 0x00)
			if nextNull == -1 {
				break
			}
			i += 3 + keySize + valueSize + nextNull
			if i >= len(dataBaseSection[0]) {
				break
			}

		}
	} else {
		i := 5
		for {

			fmt.Println("index", i)

			firstByte := dataBaseSection[0][i]

			var validTill time.Time
			// expiry in Unix time milliseconds
			if firstByte == 0xfc {
				// next 8 bytes are expiry time
				expiryTime := dataBaseSection[0][i+1 : i+9]
				expiry := binary.LittleEndian.Uint64(expiryTime)
				validTill = time.Unix(0, int64(expiry)*int64(time.Millisecond))
				i += 9
			}
			keySize := int(dataBaseSection[0][i+1])
			fmt.Println("Key size", keySize, " Key index", i+1)
			valueSize := int(dataBaseSection[0][i+2+keySize])
			fmt.Println("Value size", valueSize, " Value index", i+2+keySize)

			key := dataBaseSection[0][i+2 : i+2+keySize]
			value := dataBaseSection[0][i+3+keySize : i+3+keySize+valueSize]

			fmt.Print("Key: ", string(key))
			fmt.Println(" -- Value:", string(value))

			SetStorageItem(string(key), DataItem{value: string(value), validTill: validTill})

			// find index of next 0x00
			nextNull := slices.Index(dataBaseSection[0][i+3+keySize+valueSize:], 0x00)
			if nextNull == -1 {
				break
			}
			i += 2 + keySize + valueSize + nextNull
			if i >= len(dataBaseSection[0]) {
				break
			}

		}
	}

}
