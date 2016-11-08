package main

import (
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

type ipType struct {
	Addr     string
	Port     string
	Ssl      bool
	CreateAt time.Time
	WorkedAt time.Time
}

type linkType struct {
	Host    string
	Ssl     bool
	CheckAt time.Time
}

func initDB() {
	dbase, err := bolt.Open("ips.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	db = dbase
	createBucket([]byte("ips"))
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

func getAllIP() *mapsIP {
	allIP := newMapsIP()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ips"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ip ipType
			ip.decode(v)
			allIP.set(string(k), ip)
		}

		return nil
	})
	return allIP
}

func saveNewIP() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ips"))
		for k, v := range ips.values {
			if v.CreateAt.Sub(startAppTime) > 0 {
				ipBytes, _ := v.encode()
				b.Put([]byte(k), ipBytes)
			}
		}
		return nil
	})
	return err
}
