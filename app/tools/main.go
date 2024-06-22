package tools

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

func GetArayElement[T any](arr []T, index int, defaultValue T) T {
	if index >= 0 && index < len(arr) {
		return arr[index]
	}
	return defaultValue

}

func GenerateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 40)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetTimeStampFromEntryId(entryId string) (int64, error) {
	splited := strings.Split(entryId, "_")
	if len(splited) == 2 {
		return strconv.ParseInt(splited[1], 10, 64)
	}
	return 0, errors.New("invalid entryId")
}
