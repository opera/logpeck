package logpeck

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const configBucket string = "config"
const statBucket string = "stat"

type DB struct {
	boltdb *bolt.DB
}

var db *DB

func GetDBHandler() *DB {
	if db == nil {
		panic("DB not open")
	}
	return db
}

func OpenDB(path string) (err error) {
	boltdb, e := bolt.Open(path, 0600, nil)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Open database error: %s.", e)
		return e
	}
	err = boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(configBucket))
		if err != nil {
			return fmt.Errorf("create bucket(%s): %s", configBucket, err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(statBucket))
		if err != nil {
			return fmt.Errorf("create bucket(%s): %s", statBucket, err)
		}
		return nil
	})
	db = &DB{boltdb: boltdb}
	return nil
}

func (p *DB) Close() error {
	e := p.boltdb.Close()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Close database error: %s.", e)
	}
	return e
}

func (p *DB) makeConfigRawKey(logPath, name string) string {
	return logPath + "#" + name
}

func (p *DB) SaveConfig(config *PeckTaskConfig) error {
	rawKey := p.makeConfigRawKey(config.LogPath, config.Name)
	rawValueByte, err := json.Marshal(config)
	if err != nil {
		return err
	}
	rawValue := string(rawValueByte[:])
	//	fmt.Println(rawKey + string(" ") + rawValue)
	return p.put(configBucket, rawKey, rawValue)
}

func (p *DB) GetConfig(logPath, name string) (*PeckTaskConfig, error) {
	rawKey := p.makeConfigRawKey(logPath, name)
	rawValue := p.get(configBucket, rawKey)
	if len(rawValue) == 0 {
		return nil, errors.New("Task not exist")
	}
	//	fmt.Println(rawKV)
	var result PeckTaskConfig
	err := result.Unmarshal([]byte(rawValue))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *DB) RemoveConfig(logPath, name string) error {
	rawKey := p.makeConfigRawKey(logPath, name)
	err := p.remove(configBucket, rawKey)
	if err != nil {
		return err
	}
	return nil
}

func (p *DB) GetAllConfigs() (configs []PeckTaskConfig, err error) {
	rawKV, err := p.scan(configBucket)
	if err != nil {
		return nil, err
	}
	//	fmt.Println(rawKV)
	for _, v := range rawKV {
		config := &PeckTaskConfig{}
		err = config.Unmarshal([]byte(v))
		if err != nil {
			panic(fmt.Errorf("raw[%s], err[%s]", string(v[:]), err))
		}
		configs = append(configs, *config)
	}
	return
}

func (p *DB) makeStatRawKey(logPath, name string) string {
	return logPath + "#" + name
}

func (p *DB) SaveStat(stat *PeckTaskStat) error {
	rawKey := p.makeStatRawKey(stat.LogPath, stat.Name)
	rawValueByte, err := json.Marshal(stat)
	if err != nil {
		return err
	}
	rawValue := string(rawValueByte[:])
	//	log.Println("[Storage] SaveStat: " + rawKey + string(" ") + rawValue)
	return p.put(statBucket, rawKey, rawValue)
}

func (p *DB) GetStat(logPath, name string) (*PeckTaskStat, error) {
	rawKey := p.makeStatRawKey(logPath, name)
	rawValue := p.get(statBucket, rawKey)
	if len(rawValue) == 0 {
		return nil, errors.New("Task not exist")
	}
	//	fmt.Println("[Storage] GetStat: " + rawKey + string(" ") + rawValue)
	var result PeckTaskStat
	err := json.Unmarshal([]byte(rawValue), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *DB) RemoveStat(logPath, name string) error {
	rawKey := p.makeStatRawKey(logPath, name)
	err := p.remove(statBucket, rawKey)
	if err != nil {
		return err
	}
	return nil
}

func (p *DB) GetAllStats() (stats []PeckTaskStat, err error) {
	rawKV, err := p.scan(statBucket)
	if err != nil {
		return nil, err
	}
	//	fmt.Println(rawKV)
	for _, v := range rawKV {
		stat := &PeckTaskStat{}
		err = json.Unmarshal([]byte(v), stat)
		if err != nil {
			panic(fmt.Errorf("raw[%s], err[%s]", string(v[:]), err))
		}
		stats = append(stats, *stat)
	}
	return
}

func (p *DB) put(bucket string, key string, value string) error {
	err := p.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

func (p *DB) get(bucket string, key string) string {
	var value []byte
	p.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return string(value[:])
}

func (p *DB) remove(bucket string, key string) error {
	err := p.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}

func (p *DB) scan(bucket string) (map[string]string, error) {
	result := make(map[string]string)
	err := p.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.ForEach(func(k, v []byte) error {
			result[string(k[:])] = string(v[:])
			return nil
		})
		return nil
	})
	return result, err
}
