package logpeck

import (
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

func (p *DB) put(bucket string, key string, value string) error {
	err := p.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

func (p *DB) get(bucket string, key string) (string, error) {
	var value []byte
	err := p.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return string(value[:]), err
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
