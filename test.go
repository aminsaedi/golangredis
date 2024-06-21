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

	// fe00fb01000009626c756562657272790a73747261776265727279
	dataBaseSection = [][]byte{{0xfe, 0x00, 0xfb, 0x01, 0x00, 0x00, 0x09, 0x62, 0x6c, 0x75, 0x62, 0x65, 0x72, 0x72, 0x79, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79}}

	fmt.Println("Meta data")
	for _, meta := range metaData {
		fmt.Println(hex.EncodeToString(meta))
	}
	fmt.Println("Data base section")
	for _, data := range dataBaseSection {
		fmt.Println(hex.EncodeToString(data))
	}
	for index, data := range dataBaseSection[0][6:] {
		fmt.Print(hex.EncodeToString([]byte{data}))
		fmt.Printf(" - %c - i:%d\n", data, index)
	}

	keyValues := make(map[string]string)

	for i := 7; i < len(dataBaseSection[0]); i++ {
		spliterIndex := slices.Index(dataBaseSection[0][i:], 0x0a)
		if spliterIndex == -1 {
			break
		}
		key := dataBaseSection[0][i : i+spliterIndex]
		value := dataBaseSection[0][i+spliterIndex+1:]

		fmt.Println("Key", string(key), "Value", string(value))

		i += spliterIndex

	}

	fmt.Println(keyValues)
}
