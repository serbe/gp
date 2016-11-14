package main

import (
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

type ipType struct {
	Addr        string
	Port        string
	Ssl         bool
	isWork      bool
	isAnon      bool
	ProxyChecks int
	CreateAt    time.Time
	LastCheck   time.Time
	Response    time.Duration
}

type linkType struct {
	Host    string
	CheckAt time.Time
}

func initDB() {
	dbase, err := bolt.Open("gp.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
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

func getAllLinks() *mapsLink {
	allLinks := newMapsLink()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("links"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var link linkType
			link.decode(v)
			allLinks.set(string(k))
		}

		return nil
	})
	return allLinks
}

func saveLinks() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("links"))
		for k, v := range links.values {
			linkBytes, _ := v.encode()
			b.Put([]byte(k), linkBytes)
		}
		return nil
	})
	return err
}
