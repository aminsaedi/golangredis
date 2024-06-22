package tools

import "math/rand"

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
