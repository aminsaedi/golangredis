package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"slices"
)

var rdbMagicNumber = []byte{0x52, 0x45, 0x44, 0x49, 0x53}

func main() {
	byteData, err := os.ReadFile("empty.rdb")

	if err != nil {
		panic(err)
	}
	// data := hex.EncodeToString(byteData)

	// print file magic number

	// print and split every one byte
	// for i := 0; i < len(data); i += 2 {
	// 	print(data[i : i+2])
	// 	if i+2 < len(data) {
	// 		print(" ")
	// 	}
	// }

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

	dataBaseSection = [][]byte{{0xfe, 0x00, 0xfb, 0x01, 0x00, 0x00, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79, 0x09, 0x72, 0x61, 0x73, 0x70, 0x62, 0x65, 0x72, 0x72}}

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
