package internal

import (
	"bufio"
	"strconv"
	"strings"
)

func GetArayElement[T any](arr []T, index int, defaultValue T) T {
	if index >= 0 && index < len(arr) {
		return arr[index]
	}
	return defaultValue

}

func ParseCommand(scanner *bufio.Scanner) (command string, args []string) {
	var tokens []string
	for scanner.Scan() {
		text := scanner.Text()
		tokens = append(tokens, text)

		if len(tokens) > 0 && strings.HasPrefix(tokens[0], "*") {
			requiredItems, _ := strconv.Atoi(tokens[0][1:])
			requiredItems = requiredItems*2 + 1
			if len(tokens) == requiredItems {
				command = strings.ToUpper(tokens[2])
				args = tokens[3:]
				break
			}
		}
	}
	return
}
