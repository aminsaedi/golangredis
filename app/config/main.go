package config

import (
	"math/rand"
	"net"
)

func generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 40)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type sharedConfig struct {
	BindingPort               int
	Replicaof                 string
	MasterReplId              string
	MasterReplOffset          int
	ConnectedReplicas         []net.Conn
	FullyPropogatedReplicaIds []string
}

type propogationStatus struct {
	Commands        []string
	TransferedBytes int
}

var PropogationStatus = propogationStatus{}

var AppConfig = sharedConfig{
	MasterReplId: generateRandomString(),
}
