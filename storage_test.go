package logpeck

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
	"time"
)

func CleanTestDB(db *DB) {
	err := db.boltdb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(configBucket))
		return err
	})
	if err != nil {
		panic(err)
	}
	db.Close()
}

func TestBoltDB(*testing.T) {
	defer logExecTime(time.Now(), "open_close")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	key, value := "helloBoltDB", "logpeck"

	// test put
	fmt.Printf("put key[%s] value[%s]\n", key, value)
	err = db.put(configBucket, key, value)
	if err != nil {
		panic(err)
	}

	// test get
	value_get, e := db.get(configBucket, key)
	if e != nil {
		panic(err)
	}
	fmt.Printf("value: %s\n", value_get)
	if value_get != value {
		panic(value_get)
	}

	// test scan
	key = "2BoltDB"
	fmt.Printf("put key[%s] value[%s]\n", key, value)
	err = db.put(configBucket, key, value)
	if err != nil {
		panic(err)
	}
	res, s_err := db.scan(configBucket)
	if s_err != nil {
		panic(s_err)
	}
	if len(res) != 2 || res[key] != value {
		for k, v := range res {
			fmt.Printf("k:%s, v:%s\n", k, v)
		}
		panic(fmt.Errorf("result len: %d, value: %s", len(res), res[key]))
	}
}

func TestJson(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"
	action := "add"
	filterExpr := "panic"

	config := PeckTaskConfig{
		Name:       name,
		LogPath:    logPath,
		Action:     action,
		FilterExpr: filterExpr,
	}

	raw, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	//	fmt.Println(string(raw[:]))
	var unma PeckTaskConfig
	err = json.Unmarshal(raw, &unma)
	if err != nil {
		panic(err)
	}
	if unma.Name != name ||
		unma.LogPath != logPath ||
		unma.Action != action ||
		unma.FilterExpr != filterExpr {
		panic(unma)
	}
}
