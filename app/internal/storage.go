package internal

import (
	"errors"
	"fmt"
	"sort"
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

	fmt.Println("entryIds__:", si.entryIds, "___", entryId)
	isValid, err := isEntryIdsValid(append(si.entryIds, entryId))

	if !isValid {
		return false, err
	}

	si.entryIds = append(si.entryIds, entryId)

	for _, item := range dataItems {
		item.entryId = entryId
		item.streamId = si.key
		item.id = si.key + "_" + entryId

		SetStorageItem(item.id, item)
	}
	return true, nil
}

func (si *StreamItem) getSequenceNumbersByTime(timestamp string) []int {
	seqNumbers := make([]int, 0)
	for _, entryId := range si.entryIds {
		parts := strings.Split(entryId, "-")
		if parts[0] == timestamp {
			seqNum, _ := strconv.Atoi(parts[1])
			seqNumbers = append(seqNumbers, seqNum)
		}
	}
	// sort the sequence numbers in ascending order
	sort.Ints(seqNumbers)
	return seqNumbers
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

func isEntryIdsValid(entryIds []string) (bool, error) {

	timestampMap := make(map[string][]int)

	// Iterate over the array to split and group by timestamp
	for _, item := range entryIds {
		if item == "0-0" {
			return false, errors.New("ERR The ID specified in XADD must be greater than 0-0")
		}
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

		timestampMap[timestamp] = append(timestampMap[timestamp], seqNum)
	}

	fmt.Println("entryIds", entryIds)
	fmt.Println("timestampMap", timestampMap)

	// Check if timestamp groups are incremental
	for i := 1; i < len(entryIds); i++ {
		prevTime := strings.Split(entryIds[i-1], "-")[0]
		currTime := strings.Split(entryIds[i], "-")[0]

		if currTime < prevTime {
			return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}

	}

	// Check if sequence numbers in each group are incremental
	for _, seqNums := range timestampMap {
		for i := 1; i < len(seqNums); i++ {
			if seqNums[i] <= seqNums[i-1] {
				return false, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
			}
		}
	}

	return true, nil
}

func AddEntryToStream(stream StreamItem, entryId string, dataItems []DataItem) (ok bool, err error) {
	// si.addEntry(entryId, dataItems)
	ok, err = stream.addEntry(entryId, dataItems)
	streamStorage[stream.key] = stream
	return ok, err
}

func NormalizeEntryId(entryId string, stream StreamItem) (updatedEntryId string) {

	updatedEntryId = entryId

	if entryId == "*" {
		return strconv.Itoa(int(time.Now().UnixMilli())) + "-0"
	}

	splited := strings.Split(entryId, "-")
	timePart := splited[0]
	seqPart := splited[1]

	if seqPart == "*" {
		allSequences := stream.getSequenceNumbersByTime(timePart)
		if len(allSequences) == 0 {
			if timePart == "0" {
				updatedEntryId = "0-1"
			} else {
				updatedEntryId = timePart + "-0"
			}
		} else {
			updatedEntryId = timePart + "-" + strconv.Itoa(allSequences[len(allSequences)-1]+1)
		}
	}

	return updatedEntryId
}
