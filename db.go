package main

import (
	"fmt"

	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func get(bucket []byte, key interface{}) ([]byte, error) {
	var value []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		k, err := toBytes(key)
		if err != nil {
			return fmt.Errorf("invalid key:%v", err)
		}
		v := b.Get(k)
		if v != nil {
			value = append(value, b.Get(k)...)
		}
		return nil
	})

	return value, err
}

func put(bucket []byte, key, value interface{}) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		k, err := toBytes(key)
		if err != nil {
			return fmt.Errorf("invalid key:%v", err)
		}
		v, err := toBytes(value)
		if err != nil {
			return fmt.Errorf("invalid value:%v", err)
		}
		return b.Put(k, v)
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

func toBytes(data interface{}) ([]byte, error) {
	var (
		v   []byte
		err error
	)

	switch val := data.(type) {
	case string:
		v = []byte(val)
	case []byte:
		v = val
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v = []byte(fmt.Sprintf("%d", val))
	case float64, float32:
		v = []byte(fmt.Sprintf("%f", val))
	case fmt.Stringer:
		v = []byte(val.String())
	default:
		err = fmt.Errorf("non supported types")
	}
	return v, err
}

func isOld(u string) bool {
	v, err := get([]byte("url"), u)
	if err != nil {
		return true
	}
	var t time.Time
	err = t.UnmarshalBinary(v)
	if err != nil {
		return true
	}
	tn := time.Now()
	return tn.Sub(t) > time.Duration(15)*time.Second
}

func setTimeURL(u string) {
	tn := time.Now()
	v, err := tn.MarshalBinary()
	if err != nil {
		return
	}
	put([]byte("url"), u, v)
	return
}
