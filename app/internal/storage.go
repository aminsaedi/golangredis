package internal

import (
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
	id string
}

var plainStorage = make(map[string]DataItem)
var streamStorage = make(map[string]StreamItem)

func (di *DataItem) isValid() bool {
	if di.validTill.IsZero() {
		return true
	}
	return di.validTill.After(time.Now())
}

func (si *StreamItem) addEntry(entryId string, dateItems []DataItem) {

	for _, item := range dateItems {
		item.entryId = entryId
		item.streamId = si.id
		item.id = si.id + "_" + entryId

		SetStorageItem(item.id, item)
	}
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

func GetOrCreateStream(streamId string) StreamItem {
	item, ok := streamStorage[streamId]
	if !ok {
		if streamId == "" {
			streamId = t.GenerateRandomString()
		}
		item = StreamItem{id: streamId}
		streamStorage[streamId] = item
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
