package internal

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

func Echo(args ...string) string {
	value := args[1]
	return ToBulkString(value)
}

func Ping() string {
	return ToSimpleString("PONG")
}

func Set(args ...string) string {

	type setConfig struct {
		Key        string
		Value      string
		ExpiryType string
		ExpiryIn   string
	}

	config := setConfig{
		Key:        args[1],
		Value:      args[3],
		ExpiryType: GetArayElement(args, 5, ""),
		ExpiryIn:   GetArayElement(args, 7, ""),
	}

	toSet := DataItem{
		value: config.Value,
	}
	if len(config.ExpiryType) == 2 {
		ms, _ := strconv.Atoi(config.ExpiryIn)
		toSet.validTill = time.Now().Add(time.Duration(ms) * time.Millisecond)
	}
	SetStorageItem(config.Key, toSet)
	c.PropogationStatus.Commands = append(c.PropogationStatus.Commands, ToArray("SET", config.Key, config.Value))
	return ToSimpleString("OK")
}

func Get(args ...string) string {
	key := args[1]
	item, ok := GetStorageItem(key)
	if ok {
		return ToBulkString(item.value)
	}
	return ToSimpleError("-1")

}

func Info(selection ...string) string {
	result := map[string]string{
		"role":               "master",
		"master_replid":      c.AppConfig.MasterReplId,
		"master_repl_offset": fmt.Sprint(c.AppConfig.MasterReplOffset),
	}

	if c.AppConfig.Replicaof != "" {
		result["role"] = "slave"
	}

	return ToBulkStringFromMap(result)
}

func Replconf(args ...string) string {
	if args[1] == "ACK" {
		fmt.Println("Here: Incrementing counter", args)
		c.Counter.Start()
		c.Counter.Increment()
	}
	if args[1] == "GETACK" {
		return ToArray("REPLCONF", "ACK", strconv.Itoa(c.PropogationStatus.TransferedBytes))
	}
	return ToSimpleString("OK")
}

func Psync(args ...string) string {
	return ToSimpleString("FULLRESYNC " + c.AppConfig.MasterReplId + " " + fmt.Sprint(c.AppConfig.MasterReplOffset))
}

func RDBFileToString(filePath string) string {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return "$" + fmt.Sprint(len(dat)) + "\r\n" + string(dat)
}

func Wait(args ...string) string {
	var waitTimeInMs, leastFullyPropogatedReplicasCount int

	time.Sleep(time.Duration(50) * time.Millisecond)
	for _, replica := range c.AppConfig.ConnectedReplicas {
		go func(replica net.Conn) {
			replica.Write([]byte(ToArray("REPLCONF", "GETACK", "*")))
		}(replica)
	}

	// wg.Wait()
	if len(args) == 4 {
		waitTimeInMs, _ = strconv.Atoi(args[3])
		leastFullyPropogatedReplicasCount, _ = strconv.Atoi(args[1])
	}
	if int(c.Counter.GetCount()) < leastFullyPropogatedReplicasCount {
		time.Sleep(time.Duration(waitTimeInMs) * time.Millisecond)
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
	fmt.Println("Sending: ", c.Counter.GetCount())
	if c.Counter.GetStarted() {
		return ToSimpleInt(c.Counter.GetCount() - 1)
	}
	return ToSimpleInt(len(c.AppConfig.ConnectedReplicas))
}

func Config(args ...string) string {
	if len(args) != 4 {
		return ToSimpleError("UNKNOWN")
	}
	if args[1] == "GET" && args[3] == "dir" {
		return ToArray("dir", c.AppConfig.Dir)
	}
	if args[1] == "GET" && args[3] == "dbfilename" {
		return ToArray("dbfilename", c.AppConfig.Dbfilename)
	}
	return ToSimpleError("UNKNOWN")
}
