package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
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
		"master_replid":      config.AppConfig.MasterReplId,
		"master_repl_offset": fmt.Sprint(config.AppConfig.MasterReplOffset),
	}

	if config.AppConfig.Replicaof != "" {
		result["role"] = "slave"
	}

	return ToBulkStringFromMap(result)
}

func Replconf(args ...string) string {
	return ToSimpleString("OK")
}

func Psync(args ...string) string {
	return ToSimpleString("FULLRESYNC " + config.AppConfig.MasterReplId + " " + fmt.Sprint(config.AppConfig.MasterReplOffset))
}
