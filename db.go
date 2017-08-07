package main

import (
	"time"

	"github.com/go-pg/pg"
)

var db *pg.DB

type IP struct {
	ID       int64         `sql:"id"        json:"-"`
	Address  string        `sql:"address"   json:"address"`
	Port     string        `sql:"port"      json:"port"`
	Ssl      bool          `sql:"ssl"       json:"ssl"`
	IsWork   bool          `sql:"work"      json:"-"`
	IsAnon   bool          `sql:"anon"      json:"anon"`
	Checks   int64         `sql:"checks"    json:"-"`
	CreateAt time.Time     `sql:"create_at" json:"-"`
	UpdateAt time.Time     `sql:"update_at" json:"-"`
	Response time.Duration `sql:"response"  json:"-"`
}

type Link struct {
	ID      int64     `sql:"id"`
	Host    string    `sql:"host"`
	CheckAt time.Time `sql:"check_at"`
}

func initDB() {
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: pass,
		Database: dbname,
	})
}

func getAllIP() ([]IP, error) {
	var ips []IP
	err := db.Model(&IP{}).Select(&ips)
	if err != nil {
		errmsg("getAllIP select", err)
	}
	return ips, err
}

// func saveNewIP() error {
// 	err := db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("ips"))
// 		for k, v := range ips.values {
// 			if v.Addr != "" && v.Port != "" && v.CreateAt.Sub(startAppTime) > 0 {
// 				ipBytes, _ := v.encode()
// 				err := b.Put([]byte(k), ipBytes)
// 				if err != nil {
// 					errmsg("saveNewIP b.Put", err)
// 				}
// 			}
// 		}
// 		return nil
// 	})
// 	return err
// }

// func saveAllIP() error {
// 	err := db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("ips"))
// 		for k, v := range ips.values {
// 			if v.Addr != "" && v.Port != "" {
// 				ipBytes, _ := v.encode()
// 				err := b.Put([]byte(k), ipBytes)
// 				if err != nil {
// 					errmsg("saveAllIP b.Put", err)
// 				}
// 			}
// 		}
// 		return nil
// 	})
// 	return err
// }

func getAllLinks() ([]Link, error) {
	var links []Link
	err := db.Model(&Link{}).Select(&links)
	if err != nil {
		errmsg("getAllLinks select", err)
	}
	return links, err
}

// func saveLinks() error {
// 	err := db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte("links"))
// 		for k, v := range links.values {
// 			if v.Host != "" {
// 				linkBytes, _ := v.encode()
// 				err := b.Put([]byte(k), linkBytes)
// 				if err != nil {
// 					errmsg("saveLinks b.Put", err)
// 				}
// 			}
// 		}
// 		return nil
// 	})
// 	return err
// }
