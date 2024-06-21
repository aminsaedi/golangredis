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
	// dataBaseSection = [][]byte{{0xfe, 0x00, 0xfb, 0x01, 0x00, 0x00, 0x09, 0x62, 0x6c, 0x75, 0x62, 0x65, 0x72, 0x72, 0x79, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79}}

	// fe00fb0100000a737472617762657272790662616e616e61
	dataBaseSection = [][]byte{{0xfe, 0x00, 0xfb, 0x01, 0x00, 0x00, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79, 0x06, 0x62, 0x61, 0x6e, 0x61, 0x6e, 0x61}}

	fmt.Println("Meta data")
	for _, meta := range metaData {
		fmt.Println(hex.EncodeToString(meta))
	}
	fmt.Println("Data base section")
	for _, data := range dataBaseSection {
		fmt.Println(hex.EncodeToString(data))
	}
	for index, data := range dataBaseSection[0] {
		fmt.Print(hex.EncodeToString([]byte{data}))
		fmt.Printf(" - %c - i:%d\n", data, index)
	}

	keyValues := make(map[string]string)

	for i := 6; i < len(dataBaseSection[0]); i++ {
		fmt.Println("i", i)
		keySize := int(dataBaseSection[0][i])
		fmt.Println("Key size", keySize)
		key := dataBaseSection[0][i+1 : i+1+keySize]
		fmt.Println("Key", string(key))
		valueSize := int(dataBaseSection[0][i+keySize])
		fmt.Println("Value size", valueSize)
		value := dataBaseSection[0][i+keySize+1 : keySize+valueSize+i+1]
		fmt.Println("Value", string(value))
		i += keySize + valueSize + 1
		keyValues[string(key)] = string(value)

	}

	fmt.Println(keyValues)
}
