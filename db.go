package main

import (
	"net/url"
	"time"

	"github.com/go-pg/pg"
)

var db *pg.DB

type IP struct {
	Hostname string        `sql:"hostname"  json:"hostname"`
	IsWork   bool          `sql:"work"      json:"-"`
	IsAnon   bool          `sql:"anon"      json:"-"`
	Checks   int           `sql:"checks"    json:"-"`
	CreateAt time.Time     `sql:"create_at" json:"-"`
	UpdateAt time.Time     `sql:"update_at" json:"-"`
	Response time.Duration `sql:"response"  json:"-"`
}

type Link struct {
	Host    string    `sql:"host"`
	CheckAt time.Time `sql:"check_at"`
}

type Proxy struct {
	URL      *url.URL      `json:"url"`
	IsWork   bool          `json:"-"`
	IsAnon   bool          `json:"-"`
	Checks   int           `json:"-"`
	CreateAt time.Time     `json:"-"`
	UpdateAt time.Time     `json:"-"`
	Response time.Duration `json:"-"`
}

func initDB() {
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: pass,
		Database: dbname,
	})
}

func getAllProxy() *mapProxy {
	var i []IP
	err := db.Model(&IP{}).Select(&i)
	if err != nil {
		errmsg("getAllIP select", err)
	}
	mProxy := newMapProxy()
	for _, ip := range i {
		var proxy Proxy
		proxy, err = proxyFromIP(ip)
		if err == nil {
			mProxy.set(proxy)
		}
	}
	return mProxy
}

func existIP(ip IP) bool {
	var result bool
	_, _ = db.Query(&result, "select exists(select 1 from ips where hostname = ?)", ip.Hostname)
	return result
}

func existLink(link Link) bool {
	var result bool
	_, _ = db.Query(&result, "select exists(select 1 from links where host = ?)", link.Host)
	return result
}

func saveAllProxy(mProxy *mapProxy) {
	for _, v := range ips.values {
		ip, err := ipFromProxy(v)
		if err == nil {
			if existIP(ip) {
				_ = db.Update(&ip)
			} else {
				_ = db.Insert(&ip)
			}
		}
	}
}

func getAllLinks() *mapLink {
	var ls []Link
	err := db.Model(&Link{}).Select(&ls)
	if err != nil {
		errmsg("getAllLinks select", err)
	}
	mlink := newMapLink()
	for _, link := range ls {
		mlink.set(link.Host)
	}
	return mlink
}

func saveAllLinks(ls *mapLink) {
	for _, link := range ls.values {
		if existLink(link) {
			_ = db.Update(&link)
		} else {
			_ = db.Insert(&link)
		}
	}
}

func makeTables() {

}
