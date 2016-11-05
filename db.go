package main

import (
	"fmt"

	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

type ip struct {
	addr     string
	port     string
	ssl      bool
	createAt time.Time
	workedAt time.Time
}

type link struct {
	host      string
	ssl       bool
	lastCheck time.Time
}

func initDB() {
	dbase, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		panic(err)
	}
	db = dbase
	createBucket([]byte("ips"))
	createBucket([]byte("links"))
}

func createBucket(b []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func get(bucket, key []byte) ([]byte, error) {
	var value []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(key)
		if v != nil {
			value = append(value, b.Get(key)...)
		}
		return nil
	})

	return value, err
}

func put(bucket, key, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put(key, value)
	})

	return err
}

func del(bucket []byte, key interface{}) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		k, err := toBytes(key)
		if err != nil {
			return fmt.Errorf("invalid key:%v", err)
		}
		return b.Delete(k)
	})

	return err
}
