package config

import (
	"math/rand"
	"net"
	"sync"
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
	BindingPort       int
	Replicaof         string
	MasterReplId      string
	MasterReplOffset  int
	ConnectedReplicas []net.Conn
	Dir               string
	Dbfilename        string
}

type propogationStatus struct {
	Commands        []string
	TransferedBytes int
}

var PropogationStatus = propogationStatus{}

var AppConfig = sharedConfig{
	MasterReplId: generateRandomString(),
}

type CounterType struct {
	mu        sync.Mutex
	count     int
	isStarted bool
}

func (c *CounterType) Increment() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *CounterType) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func (c *CounterType) Start() {
	c.mu.Lock()
	c.isStarted = true
	c.mu.Unlock()
}

func (c *CounterType) GetStarted() bool {
	return c.isStarted
}

var Counter = &CounterType{}
