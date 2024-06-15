package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	c "github.com/codecrafters-io/redis-starter-go/app/config"
)

func Echo(value string) string {
	return ToBulkString(value)
}

func Ping() string {
	return ToSimpleString("PONG")
}

type SetConfig struct {
	Key        string
	Value      string
	ExpiryType string
	ExpiryIn   string
}

func Set(config SetConfig) string {

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

func Get(key string) string {
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
	fmt.Println("Replconf: ", args)
	if args[1] == "ACK" {
		fmt.Println("ALLLLLLL")
	}
	if args[1] == "GETACK" {
		return ToArray("REPLCONF", "ACK", strconv.Itoa(config.PropogationStatus.TransferedBytes))
	}
	return ToSimpleString("OK")
}

func Psync(args ...string) string {
	return ToSimpleString("FULLRESYNC " + c.AppConfig.MasterReplId + " " + fmt.Sprint(config.AppConfig.MasterReplOffset))
}

func RDBFileToString(filePath string) string {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return "$" + fmt.Sprint(len(dat)) + "\r\n" + string(dat)
}

func Wait(args ...string) string {
	fmt.Println("Wait: ", args)
	var waitTimeInMs, leastFullyPropogatedReplicasCount int
	// leastFullyPropogatedReplicasCount = c.AppConfig.ConnectedReplicasCount

	time.Sleep(time.Duration(100) * time.Millisecond)
	for _, replica := range c.AppConfig.ConnectedReplicas {
		go func() {
			replica.Write([]byte(ToArray("REPLCONF", "GETACK", "*")))
			// buff := make([]byte, 64)
			// replica.Read(buff)
			// fmt.Println("Wait: ", string(buff))
		}()
	}

	if len(args) == 4 {
		waitTimeInMs, _ = strconv.Atoi(args[3])
		leastFullyPropogatedReplicasCount, _ = strconv.Atoi(args[1])
	}
	if len(c.AppConfig.FullyPropogatedReplicaIds) < leastFullyPropogatedReplicasCount {
		time.Sleep(time.Duration(waitTimeInMs) * time.Millisecond)
	}
	return ToSimpleInt(len(c.AppConfig.FullyPropogatedReplicaIds))
	// return ToSimpleInt(100)
}
