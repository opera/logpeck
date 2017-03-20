package logpeck

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const configBucket string = "config"
const statBucket string = "stat"

type DB struct {
	bolt *bolt.DB
}

func OpenDB(path string) (db *DB, err error) {
	bolt, e := bolt.Open("path", 0600, nil)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Open database error: %s.", e)
		return nil, e
	}
	err = bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(configBucket))
		if err != nil {
			return fmt.Errorf("create bucket(%s): %s", configBucket, err)
		}
		b, err = tx.CreateBucketIfNotExists([]byte(statBucket))
		if err != nil {
			return fmt.Errorf("create bucket(%s): %s", statBucket, err)
		}
		return nil
	})
	db = &DB{bolt: bolt}
	return db, nil
}

func (p *DB) Close() error {
	e := p.Close()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Close database error: %s.", e)
	}
	return e
}

func saveKeyValue(bucket string, key string, value string) bool {

	return true
}
