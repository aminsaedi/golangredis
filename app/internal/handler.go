package internal

import (
	"fmt"
	"os"
	"strconv"
	"sync"
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

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *Counter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func Wait(args ...string) string {
	fmt.Println("Wait: ", args)
	var waitTimeInMs, leastFullyPropogatedReplicasCount int

	// var wg sync.WaitGroup
	// ch := make(chan string, 50)
	// var mu sync.Mutex
	// var count int32

	// counter := &Counter{}

	// time.Sleep(time.Duration(50) * time.Millisecond)
	// syncStatus := make(map[string]bool)
	// for _, replica := range c.AppConfig.ConnectedReplicas {
	// 	// wg.Add(1)
	// 	syncStatus[replica.RemoteAddr().String()] = false
	// 	go func(replica net.Conn) {
	// 		// defer wg.Done()
	// 		replica.Write([]byte(ToArray("REPLCONF", "GETACK", "*")))
	// 		scanner := bufio.NewScanner(replica)
	// 		reg := regexp.MustCompile(`ACK`)
	// 		for scanner.Scan() {
	// 			line := scanner.Text()
	// 			fmt.Println("Got: ", line)
	// 			if reg.MatchString(line) {
	// 				// atomic.AddInt32(&count, 1)
	// 				fmt.Println("Adding one for ", replica.RemoteAddr().String())
	// 				mu.Lock()
	// 				count++
	// 				mu.Unlock()
	// 				syncStatus[replica.RemoteAddr().String()] = true
	// 				counter.Increment()
	// 				break
	// 			}
	// 		}

	// 	}(replica)
	// }

	// wg.Wait()
	if len(args) == 4 {
		waitTimeInMs, _ = strconv.Atoi(args[3])
		leastFullyPropogatedReplicasCount, _ = strconv.Atoi(args[1])
	}
	if int(config.Counter.GetCount()) < leastFullyPropogatedReplicasCount {
		time.Sleep(time.Duration(waitTimeInMs) * time.Millisecond)
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
	fmt.Println("Sending: ", config.Counter.GetCount())
	// pp.Print(syncStatus)
	return ToSimpleInt(int(config.Counter.GetCount()))
	// return ToSimpleInt(100)
}
