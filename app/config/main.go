package config

import "math/rand"

func generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 40)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type sharedConfig struct {
	Replicaof    string
	MasterReplId string
}

var AppConfig = sharedConfig{
	MasterReplId: generateRandomString(),
}
