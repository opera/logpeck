package logpeck

import (
	"fmt"
	"testing"
	"time"
)

func logExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	fmt.Printf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}

func TestBoltDB(*testing.T) {
	defer logExecTime(time.Now(), "open_close")
	db, err := OpenDB(".test_bolt.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("put key[hello] value[logpeck]\n")
	err = db.put(configBucket, "hello", "logpeck")
	if err != nil {
		panic(err)
	}

	value, e := db.get(configBucket, "hello")
	if e != nil {
		panic(err)
	}
	fmt.Printf("value: %s\n", value)
	if value != "logpeck" {
		panic(value)
	}
}
