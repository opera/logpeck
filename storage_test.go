package logpeck

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strings"
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

func TestBoltDBAccess(*testing.T) {
	defer LogExecTime(time.Now(), "database access")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	key, value := "helloBoltDB", "logpeck"

	// test put
	log.Printf("put key[%s] value[%s]\n", key, value)
	err = db.put(configBucket, key, value)
	if err != nil {
		panic(err)
	}

	// test get
	value_get, e := db.get(configBucket, key)
	if e != nil {
		panic(err)
	}
	log.Printf("value: %s\n", value_get)
	if value_get != value {
		panic(value_get)
	}

	// test scan
	key = "2BoltDB"
	log.Printf("put key[%s] value[%s]\n", key, value)
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
			log.Printf("k:%s, v:%s\n", k, v)
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

func TestConfigsAccess(*testing.T) {
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

	defer LogExecTime(time.Now(), "config access")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			config.Name = fmt.Sprintf("%s-%d", name, j)
			config.LogPath = fmt.Sprintf("%s-%d", logPath, i)
			err = db.SaveConfig(&config)
			if err != nil {
				panic(fmt.Errorf("i[%d] j[%d] err[%s]", i, j, err))
			}
		}
	}

	configs, c_err := db.GetAllConfigs()
	if c_err != nil {
		panic(c_err)
	}
	if len(configs) != 9 {
		panic(len(configs))
	}

	for _, config := range configs {
		if !strings.Contains(config.Name, name) ||
			!strings.Contains(config.LogPath, logPath) {
			panic(configs)
		}
	}

}
