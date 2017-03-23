package logpeck

import (
	"fmt"
	"testing"
	"time"
)

func TestBoltDB(*testing.T) {
	defer logExecTime(time.Now(), "open_close")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDB()
	defer db.Close()

	key, value := "helloBoltDB", "logpeck"

	fmt.Printf("put key[%s] value[%s]\n", key, value)
	err = db.put(configBucket, key, value)
	if err != nil {
		panic(err)
	}

	value_get, e := db.get(configBucket, key)
	if e != nil {
		panic(err)
	}
	fmt.Printf("value: %s\n", value_get)
	if value_get != value {
		panic(value_get)
	}
}
