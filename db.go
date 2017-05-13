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
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ips"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ip ipType
			err := ip.decode(v)
			if err != nil {
				return err
			}
			if ip.Addr != "" && ip.Port != "" {
				allIP.set(string(k), ip)
			}
		}

		return nil
	})
	if err != nil {
		errmsg("getAllIP db.View", err)
	}
	return allIP
}

func saveNewIP() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ips"))
		for k, v := range ips.values {
			if v.Addr != "" && v.Port != "" && v.CreateAt.Sub(startAppTime) > 0 {
				ipBytes, _ := v.encode()
				err := b.Put([]byte(k), ipBytes)
				if err != nil {
					errmsg("saveNewIP b.Put", err)
				}
			}
		}
		return nil
	})
	return err
}

func saveAllIP() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ips"))
		for k, v := range ips.values {
			if v.Addr != "" && v.Port != "" {
				ipBytes, _ := v.encode()
				err := b.Put([]byte(k), ipBytes)
				if err != nil {
					errmsg("saveAllIP b.Put", err)
				}
			}
		}
		return nil
	})
	return err
}

func getAllLinks() *mapsLink {
	allLinks := newMapsLink()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("links"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var link linkType
			err := link.decode(v)
			if err != nil {
				errmsg("getAllLinks link.decode", err)
			}
			if link.Host != "" {
				allLinks.set(string(k))
			}
		}
		return nil
	})
	if err != nil {
		errmsg("getAllLinks db.View", err)
	}
	return allLinks
}

func saveLinks() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("links"))
		for k, v := range links.values {
			if v.Host != "" {
				linkBytes, _ := v.encode()
				err := b.Put([]byte(k), linkBytes)
				if err != nil {
					errmsg("saveLinks b.Put", err)
				}
			}
		}
		return nil
	})
	return err
}
