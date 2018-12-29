package logpeck

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

const testDBPath = ".unittest.db"

func CleanTestDB(db *DB) {
	err := db.boltdb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(configBucket))
		return err
	})
	if err != nil {
		panic(err)
	}
	err = db.boltdb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(statBucket))
		return err
	})
	if err != nil {
		panic(err)
	}
	db.Close()
}

func TestBoltDBAccess(*testing.T) {
	defer LogExecTime(time.Now(), "database access")
	err := OpenDB(testDBPath)
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
	valueGet := db.get(configBucket, key)
	log.Printf("value: %s\n", valueGet)
	if valueGet != value {
		panic(valueGet)
	}

	// test scan
	key = "2BoltDB"
	log.Printf("put key[%s] value[%s]\n", key, value)
	err = db.put(configBucket, key, value)
	if err != nil {
		panic(err)
	}
	res, err := db.scan(configBucket)
	if err != nil {
		panic(err)
	}
	if len(res) != 2 || res[key] != value {
		for k, v := range res {
			log.Printf("k:%s, v:%s\n", k, v)
		}
		panic(fmt.Errorf("result len: %d, value: %s", len(res), res[key]))
	}

	// test remove
	log.Printf("remove key[%s]\n", key)
	err = db.remove(configBucket, key)
	if err != nil {
		panic(err)
	}
	valueGet = db.get(configBucket, key)
	if valueGet != "" {
		panic(valueGet)
	}
}

func TestJson(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"
	filterExpr := "panic"

	config := PeckTaskConfig{
		Name:     name,
		LogPath:  logPath,
		Keywords: filterExpr,
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
		unma.Keywords != filterExpr {
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
		Name:   "ElasticSearch",
		Config: ESConfig,
	}
	config := PeckTaskConfig{
		Name:     name,
		LogPath:  logPath,
		Keywords: filterExpr,
		Sender:   esconfig,
	}

	defer LogExecTime(time.Now(), "config access")
	err := OpenDB(testDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveConfig
	for i := 0; i < 10; i++ {
		config.Name = fmt.Sprintf("%s-%d", name, i)
		config.LogPath = fmt.Sprintf("%s-%d", logPath, i)
		err = db.SaveConfig(&config)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	// Test GetConfig
	configGetTmp := &PeckTaskConfig{
		Name:    name + "-0",
		LogPath: logPath + "-0",
	}
	configGet, e := db.GetConfig(configGetTmp.Name)
	if e != nil {
		panic(e)
	}
	if configGet.Name != configGetTmp.Name ||
		configGet.LogPath != configGetTmp.LogPath {
		fmt.Printf("%s vs %s, %s vs %s\n", configGet.Name, configGetTmp.Name, configGet.LogPath, configGetTmp.LogPath)
		panic(configGet)
	}

	// Test GetAllConfigs
	configs, err := db.GetAllConfigs()
	if err != nil {
		panic(err)
	}
	if len(configs) != 10 {
		panic(len(configs))
	}

	for _, config := range configs {
		if !strings.Contains(config.Name, name) ||
			!strings.Contains(config.LogPath, logPath) {
			panic(configs)
		}
	}

	// Test RemoveConfig
	for i := 0; i < 10; i++ {
		config.Name = fmt.Sprintf("%s-%d", name, i)
		config.LogPath = fmt.Sprintf("%s-%d", logPath, i)
		err = db.RemoveConfig(config.Name)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	configs, err = db.GetAllConfigs()
	if len(configs) != 0 {
		panic(len(configs))
	}
}

func TestStatsAccess(*testing.T) {
	name := "test_peck_task"

	stat := PeckTaskStat{
		Name: name,
		Stop: true,
	}

	defer LogExecTime(time.Now(), "stats access")
	err := OpenDB(testDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveStat
	for i := 0; i < 10; i++ {
		stat.Name = fmt.Sprintf("%s-%d", name, i)
		err = db.SaveStat(&stat)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	// Test GetStat
	statGetTmp := &PeckTaskStat{
		Name: name + "-0",
	}
	statGet, e := db.GetStat(statGetTmp.Name)
	if e != nil {
		panic(e)
	}
	if statGet.Name != statGetTmp.Name {
		panic(statGet)
	}

	// Test GetAllStats
	stats, err := db.GetAllStats()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", stats)
	if len(stats) != 10 {
		panic(len(stats))
	}

	for _, stat := range stats {
		if !strings.Contains(stat.Name, name) {
			panic(stats)
		}
	}

	// Test RemoveStat
	for i := 0; i < 10; i++ {
		stat.Name = fmt.Sprintf("%s-%d", name, i)
		err = db.RemoveStat(stat.Name)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	stats, err = db.GetAllStats()
	if len(stats) != 0 {
		panic(len(stats))
	}
}

func TestConfigCompat(*testing.T) {
	name := "test_peck_task"
	logPath := "./test.log"

	config := PeckTaskConfig{
		Name:    name,
		LogPath: logPath,
	}

	defer LogExecTime(time.Now(), "config access")
	err := OpenDB(testDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveConfig
	for i := 0; i < 10; i++ {
		config.Name = fmt.Sprintf("%s#%d", name, i)
		config.LogPath = fmt.Sprintf("%s-%d", logPath, i)
		err = db.SaveConfig(&config)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	// Test GetAllConfigs
	configs, err := db.GetAllConfigs()
	if err != nil {
		panic(err)
	}
	if len(configs) != 10 {
		panic(len(configs))
	}

	// Test GetConfig
	getTmp := &PeckTaskConfig{
		Name:    name + "#0",
		LogPath: logPath + "-0",
	}
	configGet, e := db.GetConfig("0")
	if e != nil {
		panic(e)
	}
	if configGet.Name != getTmp.Name ||
		configGet.LogPath != getTmp.LogPath {
		fmt.Printf("%s %s %s %s\n", configGet.Name, getTmp.Name, configGet.LogPath, getTmp.LogPath)
		panic(configGet)
	}

}
func TestStatCompat(*testing.T) {
	name := "test_peck_task"

	stat := PeckTaskStat{
		Name: name,
		Stop: true,
	}

	defer LogExecTime(time.Now(), "stats access")
	err := OpenDB(testDBPath)
	if err != nil {
		panic(err)
	}
	db := GetDBHandler()
	defer CleanTestDB(db)

	// Test SaveStat
	for i := 0; i < 10; i++ {
		stat.Name = fmt.Sprintf("%s#%d", name, i)
		err = db.SaveStat(&stat)
		if err != nil {
			panic(fmt.Errorf("i[%d] err[%s]", i, err))
		}
	}

	// Test GetAllStats
	stats, err := db.GetAllStats()
	if err != nil {
		panic(err)
	}
	if len(stats) != 10 {
		panic(len(stats))
	}

	// Test GetStat
	statGetTmp := &PeckTaskStat{
		Name: name + "#0",
	}
	statGet, e := db.GetStat("0")
	if e != nil {
		panic(e)
	}
	if statGet.Name != statGetTmp.Name {
		fmt.Printf("%s %s \n", statGet.Name, statGetTmp.Name)
		panic(statGet)
	}

}
