package config

import (
	"fmt"
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
	ids       []string
	isStarted bool
}

func (c *CounterType) Increment(id string) {
	c.mu.Lock()
	c.ids = append(c.ids, id)
	c.mu.Unlock()
}

func (c *CounterType) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	// return number of unique ids
	unique := make(map[string]bool)
	for _, id := range c.ids {
		unique[id] = true
	}
	fmt.Println("ids", c.ids, "unique", unique)
	return len(unique)
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
