package internal

import "fmt"

func ToBulkString(input ...string) string {
	// return "$" + fmt.Sprint(len(input)) + "\r\n" + input + "\r\n"
	totalLength := 0
	finalString := ""
	for index, v := range input {
		finalString += v + "\r\n"
		if index == 0 {
			totalLength += len(v)
		} else {
			totalLength += len(v) + 2
		}
	}
	fmt.Println("Total length: ", totalLength)
	fmt.Printf("Final string: %q\n", finalString)
	fmt.Printf("Out: %q", "$"+fmt.Sprint(totalLength)+"\r\n"+finalString)
	return "$" + fmt.Sprint(totalLength) + "\r\n" + finalString
}

func ToBulkStringFromMap(input map[string]string) string {
	// key values should be separated by :
	var finalString string
	for k, v := range input {
		finalString += k + ":" + v + "\r\n"
	}
	return ToBulkString(finalString)
}
func ToSimpleString(input string) string {
	return "+" + input + "\r\n"
}

func ToSimpleError(input string) string {
	return "$" + input + "\r\n"
}

func ToArray(input ...string) string {
	finalString := ""
	for _, v := range input {
		finalString += "$" + fmt.Sprint(len(v)) + "\r\n" + v + "\r\n"
	}
	return "*" + fmt.Sprint(len(input)) + "\r\n" + finalString
}

func ToSimpleInt(input int) string {
	return ":" + fmt.Sprint(input) + "\r\n"
}
