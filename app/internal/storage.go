package internal

import (
	"time"

	t "github.com/codecrafters-io/redis-starter-go/app/tools"
)

type DataItem struct {
	id        string
	streamId  string
	value     string
	validTill time.Time
}

type StreamItem struct {
	id string
}

func (di *DataItem) isValid() bool {
	if di.validTill.IsZero() {
		return true
	}
	return di.validTill.After(time.Now())
}

var plainStorage = make(map[string]DataItem)
var streamStorage = make(map[string]StreamItem)

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
