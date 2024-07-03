package internal

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	c "github.com/codecrafters-io/redis-starter-go/app/config"
	t "github.com/codecrafters-io/redis-starter-go/app/tools"
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
		ExpiryType: t.GetArayElement(args, 5, ""),
		ExpiryIn:   t.GetArayElement(args, 7, ""),
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
	return ToBulkString("")

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

func Replconf(conn net.Conn, args ...string) string {
	if args[1] == "ACK" {
		c.UniqeCounter.Start()
		c.UniqeCounter.Increment(conn.RemoteAddr().String())
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
	if int(c.UniqeCounter.GetCount()) < leastFullyPropogatedReplicasCount {
		time.Sleep(time.Duration(waitTimeInMs) * time.Millisecond)
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
	fmt.Println("Sending: ", c.UniqeCounter.GetCount())
	if c.UniqeCounter.GetStarted() {
		return ToSimpleInt(c.UniqeCounter.GetCount())
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

func Keys(args ...string) string {
	keys := GetAllKeys()
	return ToArray(keys...)
}

func Type(args ...string) string {
	key := args[1]
	if IsStorageKeyValid(key) {
		return ToSimpleString("string")
	} else if IsStreamKeyValid(key) {
		return ToSimpleString("stream")
	}
	return ToSimpleString("none")
}

var addedItems = 0

func Xadd(args ...string) string {
	streamKey := args[1]
	entryId := args[3]
	fields := args[5:]

	dataItems := make([]DataItem, 0)

	for i := 0; i < len(fields); i += 3 {
		dataItems = append(dataItems, DataItem{
			key:   fields[i],
			value: fields[i+2],
		})
	}

	stream := GetOrCreateStream(streamKey)

	entryId = NormalizeEntryId(entryId, stream)

	ok, err := AddEntryToStream(stream, entryId, dataItems)
	if !ok {
		return ToSimpleError(err.Error())
	}
	addedItems++
	return ToBulkString(entryId)
}

func Xrange(args ...string) string {
	streamKey := args[1]
	start := args[3]
	end := args[5]

	stream := GetOrCreateStream(streamKey)

	firstIndex := slices.Index(stream.entryIds, start)
	lastIndex := slices.Index(stream.entryIds, end)

	if start == "-" {
		firstIndex = 0
	}
	if end == "+" {
		lastIndex = len(stream.entryIds) - 1
	}

	entryIds := stream.entryIds[firstIndex : lastIndex+1]

	result := make([]string, 0)

	for _, entryId := range entryIds {
		temp := make([]string, 0)
		temp = append(temp, entryId)

		temp2 := make([]string, 0)
		item, ok := GetStorageItem(streamKey + "_" + entryId)
		if ok {
			temp2 = append(temp2, item.key)
			temp2 = append(temp2, item.value)
		}

		temp2Str := ToArray(temp2...)
		temp = append(temp, temp2Str)
		result = append(result, ToArray(temp...))

	}

	finalResult := ToArray(result...)

	return finalResult
}

func Xread(args ...string) string {
	fmt.Println("XREAD", args)

	if strings.ToUpper(args[1]) == "BLOCK" {
		blockTime, _ := strconv.Atoi(args[3])
		timer := time.NewTimer(time.Duration(blockTime) * time.Millisecond)
		addedItems = 0
	loop:
		for {
			select {
			case <-timer.C:
				break loop
			default:
				if addedItems > 0 {
					break loop
				}
			}
		}
		args = args[4:]
	}

	getValue := func(streamKey string, entryId string) string {
		stream := GetOrCreateStream(streamKey)

		index := slices.Index(stream.entryIds, entryId)

		if index == -1 {
			index = 0
		}

		result := make([]string, 0)

		result = append(result, streamKey)

		entryIds := stream.entryIds[index:]

		for _, entryId := range entryIds {
			temp := make([]string, 0)
			temp = append(temp, entryId)

			temp2 := make([]string, 0)
			item, ok := GetStorageItem(streamKey + "_" + entryId)
			if ok {
				temp2 = append(temp2, item.key)
				temp2 = append(temp2, item.value)
			}

			temp2Str := ToArray(temp2...)
			temp = append(temp, temp2Str)
			result = append(result, ToArray(ToArray(temp...)))

		}

		resultStr := ToArray(result...)

		return resultStr
	}

	result := make([]string, 0)

	if len(args) == 6 {
		streamKey := args[3]
		entryId := args[5]

		result = append(result, getValue(streamKey, entryId))
	} else if len(args) == 10 {
		streamKey1 := args[3]
		entryId1 := args[7]

		streamKey2 := args[5]
		entryId2 := args[9]

		result = append(result, getValue(streamKey1, entryId1))
		result = append(result, getValue(streamKey2, entryId2))
	}

	fmt.Printf("Result: %q\n", result)

	finalResult := ToArray(result...)

	return finalResult
}
