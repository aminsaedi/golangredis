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

	// fe00fb0400000a73747261776265727279066f72616e6765000970696e656170706c650970696e656170706c65000662616e616e61056d616e676f00056170706c650a73747261776265727279
	// dataBaseSection = [][]byte{{0xfe, 0x00, 0xfb, 0x04, 0x00, 0x00, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79, 0x06, 0x6f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x00, 0x09, 0x70, 0x69, 0x6e, 0x65, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x09, 0x70, 0x69, 0x6e, 0x65, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x00, 0x06, 0x62, 0x61, 0x6e, 0x61, 0x6e, 0x61, 0x05, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x00, 0x05, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x0a, 0x73, 0x74, 0x72, 0x61, 0x77, 0x62, 0x65, 0x72, 0x72, 0x79}}

	fakeStringInput := "fe00fb0404fc000c288ac701000000056170706c6509626c75656265727279fc000c288ac70100000004706561720662616e616e61fc009cef127e0100000009726173706265727279056772617065fc000c288ac701000000066f72616e6765066f72616e6765"

	decoded, _ := hex.DecodeString(fakeStringInput)
	dataBaseSection = [][]byte{decoded}

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

	// totalSize := int(dataBaseSection[0][3])
	withExpiryItemSize := int(dataBaseSection[0][4])

	keyValues := make(map[string]string)

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

			keyValues[string(key)] = string(value)

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

			firstByte := dataBaseSection[0][i]

			// expiry in Unix time milliseconds
			if firstByte == 0xfc {
				// next 8 bytes are expiry time
				expiryTime := dataBaseSection[0][i+1 : i+9]
				fmt.Println("Expiry time", hex.EncodeToString(expiryTime))
				i += 9
			}

		}
	}

	fmt.Println(keyValues)
}
