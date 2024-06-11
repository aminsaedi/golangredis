package internal

import (
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
	if config.AppConfig.Replicaof != "" {
		return ToBulkString("role:slave")
	}
	// return ToBulkString("role:master\r\nmaster_replid:" + config.AppConfig.MasterReplId + "\r\nmaster_repl_offset:0")
	return ToBulkString(
		"role:master",
		"master_replid:"+config.AppConfig.MasterReplId,
		"master_repl_offset:0",
	)
}
