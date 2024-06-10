package config

type sharedConfig struct {
	Replicaof string
}

var AppConfig = sharedConfig{}
