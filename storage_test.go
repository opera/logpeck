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

const kTestDBPath string = ".unittest.db"

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
	value_get := db.get(configBucket, key)
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

	// test remove
	log.Printf("remove key[%s] value[%s]\n", key)
	err = db.remove(configBucket, key)
	if err != nil {
		panic(err)
	}
	value_get = db.get(configBucket, key)
	if value_get != "" {
		panic(value_get)
	}
}

func TestJson(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"
	filterExpr := "panic"

	config := PeckTaskConfig{
		Name:       name,
		LogPath:    logPath,
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
		unma.FilterExpr != filterExpr {
		panic(unma)
	}
}

func TestConfigsAccess(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"
	filterExpr := "panic"
	ESConfig := ElasticSearchConfig{
		Hosts: []string{"127.0.0.1:9200"},
		Index: "test",
		Type:  "testType",
	}
	esconfig := SenderConfig{
		Name:   "ElasticsearchConfig",
		Config: ESConfig,
	}
	config := PeckTaskConfig{
		Name:         name,
		LogPath:      logPath,
		FilterExpr:   filterExpr,
		SenderConfig: esconfig,
	}

	defer LogExecTime(time.Now(), "config access")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveConfig
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

	// Test GetConfig
	config_get_tmp := &PeckTaskConfig{
		Name:    name + "-0",
		LogPath: logPath + "-0",
	}
	config_get, e := db.GetConfig(config_get_tmp.LogPath, config_get_tmp.Name)
	if e != nil {
		panic(e)
	}
	if config_get.Name != config_get_tmp.Name ||
		config_get.LogPath != config_get_tmp.LogPath {
		panic(config_get)
	}

	// Test GetAllConfigs
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

	// Test RemoveConfig
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			config.Name = fmt.Sprintf("%s-%d", name, j)
			config.LogPath = fmt.Sprintf("%s-%d", logPath, i)
			err = db.RemoveConfig(config.LogPath, config.Name)
			if err != nil {
				panic(fmt.Errorf("i[%d] j[%d] err[%s]", i, j, err))
			}
		}
	}

	configs, c_err = db.GetAllConfigs()
	if len(configs) != 0 {
		panic(len(configs))
	}
}

func TestStatsAccess(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"

	stat := PeckTaskStat{
		Name:    name,
		LogPath: logPath,
		Stop:    true,
	}

	defer LogExecTime(time.Now(), "stats access")
	err := OpenDB(kTestDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveStat
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			stat.Name = fmt.Sprintf("%s-%d", name, j)
			stat.LogPath = fmt.Sprintf("%s-%d", logPath, i)
			err = db.SaveStat(&stat)
			if err != nil {
				panic(fmt.Errorf("i[%d] j[%d] err[%s]", i, j, err))
			}
		}
	}

	// Test GetStat
	stat_get_tmp := &PeckTaskStat{
		Name:    name + "-0",
		LogPath: logPath + "-0",
	}
	stat_get, e := db.GetStat(stat_get_tmp.LogPath, stat_get_tmp.Name)
	if e != nil {
		panic(e)
	}
	if stat_get.Name != stat_get_tmp.Name ||
		stat_get.LogPath != stat_get_tmp.LogPath {
		panic(stat_get)
	}

	// Test GetAllStats
	stats, c_err := db.GetAllStats()
	if c_err != nil {
		panic(c_err)
	}
	if len(stats) != 9 {
		panic(len(stats))
	}

	for _, stat := range stats {
		if !strings.Contains(stat.Name, name) ||
			!strings.Contains(stat.LogPath, logPath) {
			panic(stats)
		}
	}

	// Test RemoveStat
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			stat.Name = fmt.Sprintf("%s-%d", name, j)
			stat.LogPath = fmt.Sprintf("%s-%d", logPath, i)
			err = db.RemoveStat(stat.LogPath, stat.Name)
			if err != nil {
				panic(fmt.Errorf("i[%d] j[%d] err[%s]", i, j, err))
			}
		}
	}

	stats, c_err = db.GetAllStats()
	if len(stats) != 0 {
		panic(len(stats))
	}
}
