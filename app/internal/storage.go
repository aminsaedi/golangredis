package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	t "github.com/codecrafters-io/redis-starter-go/app/tools"
)

type DataItem struct {
	id        string
	streamId  string
	entryId   string
	key       string
	value     string
	validTill time.Time
}

type StreamItem struct {
	key      string
	entryIds []string
}

var plainStorage = make(map[string]DataItem)
var streamStorage = make(map[string]StreamItem)

func (di *DataItem) isValid() bool {
	if di.validTill.IsZero() {
		return true
	}
	return di.validTill.After(time.Now())
}

func (si *StreamItem) addEntry(entryId string, dataItems []DataItem) (ok bool, err error) {

	// add entryId to the stream entryIds
	si.entryIds = append(si.entryIds, entryId)

	fmt.Println("entryIds__:", si.entryIds, "___", entryId)
	isValid := isEntryIdsValid(si.key)

	if !isValid {
		fmt.Println("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	for _, item := range dataItems {
		item.entryId = entryId
		item.streamId = si.key
		item.id = si.key + "_" + entryId

		SetStorageItem(item.id, item)
	}
	return true, nil
}

func GetStorageItem(key string) (DataItem, bool) {
	item, ok := plainStorage[key]
	if ok && item.isValid() {
		return item, true
	}
	return DataItem{}, false
}

func SetStorageItem(key string, item DataItem) {
	if item.id == "" {
		item.id = t.GenerateRandomString()
	}
	if item.key == "" {
		item.key = key
	}
	plainStorage[key] = item
}

func GetAllKeys() []string {
	keys := make([]string, 0)
	for key, item := range plainStorage {
		if item.isValid() {
			keys = append(keys, key)
		}
	}
	return keys
}

func GetOrCreateStream(streamKey string) StreamItem {
	item, ok := streamStorage[streamKey]
	if !ok {
		if streamKey == "" {
			streamKey = t.GenerateRandomString()
		}
		item = StreamItem{key: streamKey, entryIds: make([]string, 0)}
		streamStorage[streamKey] = item
	}
	return item
}

func IsStorageKeyValid(key string) bool {
	_, ok := plainStorage[key]
	return ok
}

func IsStreamKeyValid(streamKey string) bool {
	_, ok := streamStorage[streamKey]
	return ok
}

func isEntryIdsValid(streamKey string) bool {
	stream, ok := streamStorage[streamKey]
	if !ok {
		return false
	}

	// fmt.Println("Ids:", stream.entryIds)

	timestampMap := make(map[string][]int)

	// Iterate over the array to split and group by timestamp
	for _, item := range stream.entryIds {
		parts := strings.Split(item, "-")
		if len(parts) != 2 {
			// Invalid format, skip this item or handle error
			continue
		}
		timestamp := parts[0]
		seqNum, err := strconv.Atoi(parts[1])
		if err != nil {
			// Invalid sequence number, skip this item or handle error
			continue
		}

		// Add sequence number to the corresponding timestamp group
		timestampMap[timestamp] = append(timestampMap[timestamp], seqNum)
	}

	fmt.Println("timestampMap", timestampMap)

	// Check if sequence numbers in each group are incremental
	for _, seqNums := range timestampMap {
		for i := 1; i < len(seqNums); i++ {
			if seqNums[i] <= seqNums[i-1] {
				return false
			}
		}
	}

	return true
}

func AddEntryToStream(stream StreamItem, entryId string, dataItems []DataItem) (ok bool, err error) {
	// si.addEntry(entryId, dataItems)
	ok, err = stream.addEntry(entryId, dataItems)
	streamStorage[stream.key] = stream
	return ok, err
}
