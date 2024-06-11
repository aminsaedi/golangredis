package internal

import "fmt"

func ToBulkString(input ...string) string {
	// return "$" + fmt.Sprint(len(input)) + "\r\n" + input + "\r\n"
	totalLength := 0
	finalString := ""
	for _, v := range input {
		finalString += v + "\r\n"
		totalLength += len(v) + 2
	}
	fmt.Println("Total length: ", totalLength)
	fmt.Printf("Final string: %q\n", finalString)
	fmt.Printf("Out: %q", "$"+fmt.Sprint(totalLength)+"\r\n"+finalString)
	return "$" + fmt.Sprint(totalLength) + "\r\n" + finalString
}

func ToSimpleString(input string) string {
	return "+" + input + "\r\n"
}

func ToSimpleError(input string) string {
	return "$" + input + "\r\n"
}
