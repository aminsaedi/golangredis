package internal

import "time"

type DataItem struct {
	value     string
	validTill time.Time
}

func (di *DataItem) isValid() bool {
	if di.validTill.IsZero() {
		return true
	}
	return di.validTill.After(time.Now())
}

var storage = make(map[string]DataItem)

func GetStorageItem(key string) (DataItem, bool) {
	item, ok := storage[key]
	if ok && item.isValid() {
		return item, true
	}
	return DataItem{}, false
}

func SetStorageItem(key string, item DataItem) {
	storage[key] = item
}

func GetAllKeys() []string {
	keys := make([]string, 0)
	for key, item := range storage {
		if item.isValid() {
			keys = append(keys, key)
		}
	}
	return keys
}
