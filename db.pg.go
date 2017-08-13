package main

// import (
// 	"net/url"
// 	"time"

// 	"github.com/go-pg/pg"
// )

// // Proxy - proxy unit
// type Proxy struct {
// 	Insert   bool          `sql:"-"           json:"-"`
// 	Update   bool          `sql:"-"           json:"-"`
// 	Hostname string        `sql:"hostname,pk" json:"hostname"`
// 	URL      *url.URL      `sql:"-"           json:"-"`
// 	Host     string        `sql:"host"        json:"-"`
// 	Port     string        `sql:"port"        json:"-"`
// 	IsWork   bool          `sql:"work"        json:"-"`
// 	IsAnon   bool          `sql:"anon"        json:"-"`
// 	Checks   int           `sql:"checks"      json:"-"`
// 	CreateAt time.Time     `sql:"create_at"   json:"-"`
// 	UpdateAt time.Time     `sql:"update_at"   json:"-"`
// 	Response time.Duration `sql:"response"    json:"-"`
// }

// func initDB() *pg.DB {
// 	return pg.Connect(&pg.Options{
// 		User:     user,
// 		Password: pass,
// 		Database: dbname,
// 	})
// }

// func getAllProxy(db *pg.DB) *mapProxy {
// 	debugmsg("start getAllProxy")
// 	var proxies []Proxy
// 	err := db.Model(&Proxy{}).Select(&proxies)
// 	if err != nil {
// 		errmsg("getAllProxy select", err)
// 	}
// 	mProxy := newMapProxy()
// 	for _, proxy := range proxies {
// 		proxy.URL, err = url.Parse(proxy.Hostname)
// 		if err != nil {
// 			errmsg("getAllProxy url.Parse", err)
// 		}
// 		mProxy.set(proxy)
// 	}
// 	debugmsg("end getAllProxy, load proxy", len(mProxy.values))
// 	return mProxy
// }

// func get100Proxy(db *pg.DB) *mapProxy {
// 	debugmsg("start getAllProxy")
// 	var proxies []Proxy
// 	err := db.Model(&Proxy{}).Limit(100).Select(&proxies)
// 	if err != nil {
// 		errmsg("getAllProxy select", err)
// 	}
// 	mProxy := newMapProxy()
// 	for _, proxy := range proxies {
// 		proxy.URL, err = url.Parse(proxy.Hostname)
// 		if err != nil {
// 			errmsg("getAllProxy url.Parse", err)
// 		}
// 		mProxy.set(proxy)
// 	}
// 	debugmsg("end getAllProxy, load proxy", len(mProxy.values))
// 	return mProxy
// }

// func saveAllProxy(db *pg.DB, mProxy *mapProxy) {
// 	debugmsg("start saveAllProxy")
// 	var u, i int64
// 	for _, proxy := range mProxy.values {
// 		if proxy.Update {
// 			u++
// 			err := db.Update(&proxy)
// 			if err != nil {
// 				errmsg("saveAllProxy Update", err)
// 			}
// 		}
// 		if proxy.Insert {
// 			i++
// 			err := db.Insert(&proxy)
// 			if err != nil {
// 				errmsg("saveAllLinks Insert", err)
// 			}
// 		}
// 	}
// 	debugmsg("update proxy", u)
// 	debugmsg("insert proxy", i)
// 	debugmsg("end getAllProxy")
// }

// func updateAllProxy(db *pg.DB, mProxy *mapProxy) {
// 	debugmsg("start updateAllProxy")
// 	for _, proxy := range mProxy.values {
// 		err := db.Update(&proxy)
// 		if err != nil {
// 			errmsg("updateAllProxy update", err)
// 		}
// 	}
// 	debugmsg("end updateAllProxy, update proxy", len(mProxy.values))
// }

// func getAllLinks(db *pg.DB) *mapLink {
// 	debugmsg("start getAllLinks")
// 	var ls []Link
// 	err := db.Model(&Link{}).Select(&ls)
// 	if err != nil {
// 		errmsg("getAllLinks select", err)
// 	}
// 	mLink := newMapLink()
// 	for _, link := range ls {
// 		mLink.set(link)
// 	}
// 	debugmsg("end getAllLinks, load links", len(mLink.values))
// 	return mLink
// }

// func saveAllLinks(db *pg.DB, ls *mapLink) {
// 	debugmsg("start saveAllLinks")
// 	var (
// 		u, i int64
// 	)
// 	for _, link := range ls.values {
// 		if link.Insert {
// 			i++
// 			err := db.Insert(&link)
// 			if err != nil {
// 				errmsg("saveAllLinks Insert", err)
// 			}
// 		} else {
// 			u++
// 			err := db.Update(&link)
// 			if err != nil {
// 				errmsg("saveAllLinks Update", err)
// 			}
// 		}
// 	}
// 	debugmsg("update links", u)
// 	debugmsg("insert links", i)
// 	debugmsg("end saveAllLinks")
// }

// // func makeTables() {
// // 	db.ExecOne(`
// // 		CREATE TABLE IF NOT EXISTS proxies (
// // 			hostname  text PRIMARY KEY,
// // 			host      text,
// // 			port      text,
// // 			work      boolean,
// // 			anon      boolean,
// // 			checks    integer,
// // 			create_at timestamptz DEFAULT now(),
// // 			update_at timestamptz,
// // 			response  integer,
// // 			UNIQUE(hostname)
// // 		);

// // 		CREATE TABLE IF NOT EXISTS links (
// // 			hostname  text PRIMARY KEY,
// // 			update_at timestamptz DEFAULT now(),
// // 			UNIQUE(hostname)
// // 		);

// // 		INSERT INTO links (hostname) VALUES
// // 			('https://hidester.com/proxydata/php/data.php?mykey=data&offset=0&limit=1000&orderBy=latest_check&sortOrder=DESC&country=&port=&type=undefined&anonymity=undefined&ping=undefined&gproxy=2'),
// // 			('http://gatherproxy.com/embed/'),
// // 			('http://txt.proxyspy.net/proxy.txt'),
// // 			('http://webanetlabs.net/publ/24'),
// // 			('http://awmproxy.com/freeproxy.php'),
// // 			('http://www.samair.ru/proxy/type-01.htm'),
// // 			('https://www.us-proxy.org/'),
// // 			('http://free-proxy-list.net/'),
// // 			('http://www.proxynova.com/proxy-server-list/'),
// // 			('http://proxyserverlist-24.blogspot.ru/'),
// // 			('http://gatherproxy.com/'),
// // 			('https://hidemy.name/ru/proxy-list/'),
// // 			('https://hidemy.name/en/proxy-list/?type=hs&anon=34#list'),
// // 			('https://free-proxy-list.com'),
// // 			('https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=anonymous&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=2&connect_time%5B%5D=3&up_time=40&search=Search'),
// // 			('http://www.idcloak.com/proxylist/free-proxy-servers-list.html'),
// // 			('https://premproxy.com/list/'),
// // 			('https://proxy-list.org/english/index.php'),
// // 			('https://www.sslproxies.org/')
// // 		ON CONFLICT (hostname) DO NOTHING;
// // 	`)
// // }
