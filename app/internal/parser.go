package internal

import "fmt"

func ToBulkString(input ...string) string {
	// return "$" + fmt.Sprint(len(input)) + "\r\n" + input + "\r\n"
	totalLength := 0
	finalString := ""
	for _, v := range input {
		totalLength += len(v)
		finalString += v + "\r\n"
	}
	return "$" + fmt.Sprint(totalLength) + "\r\n" + finalString
}

func ToSimpleString(input string) string {
	return "+" + input + "\r\n"
}

func ToSimpleError(input string) string {
	return "$" + input + "\r\n"
}
