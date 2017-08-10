package main

import (
	"net/url"

	"github.com/go-pg/pg"
)

var db *pg.DB

func initDB() {
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: pass,
		Database: dbname,
	})
}

func getAllProxy() *mapProxy {
	var proxies []Proxy
	err := db.Model(&Proxy{}).Select(&proxies)
	if err != nil {
		errmsg("getAllProxy select", err)
	}
	mProxy := newMapProxy()
	debugmsg("load proxy from db", len(proxies))
	for _, proxy := range proxies {
		proxy.URL, err = url.Parse(proxy.Hostname)
		if err != nil {
			errmsg("getAllProxy url.Parse", err)
		}
		mProxy.set(proxy)
	}
	debugmsg("load proxy", len(mProxy.values))
	return mProxy
}

func saveAllProxy(mProxy *mapProxy) {
	var u, i int64
	for _, proxy := range mP.values {
		if proxy.Update {
			u++
			_ = db.Update(&proxy)
		}
		if proxy.Insert {
			i++
			_ = db.Insert(&proxy)
		}
	}
	debugmsg("update proxy", u)
	debugmsg("insert proxy", i)
}

func getAllLinks() *mapLink {
	var ls []Link
	err := db.Model(&Link{}).Select(&ls)
	if err != nil {
		errmsg("getAllLinks select", err)
	}
	mlink := newMapLink()
	for _, link := range ls {
		mlink.set(link)
	}
	debugmsg("load links", len(mlink.values))
	return mlink
}

func saveAllLinks(ls *mapLink) {
	var (
		u, i int64
	)
	for _, link := range ls.values {
		if link.Insert {
			i++
			_ = db.Insert(&link)
		} else {
			u++
			_ = db.Update(&link)
		}
	}
	debugmsg("update links", u)
	debugmsg("insert links", i)
}

func makeTables() {
	db.ExecOne(`
		CREATE TABLE IF NOT EXISTS proxies (
			hostname  text PRIMARY KEY,
			host      text,
			port      text,
			work      boolean,
			anon      boolean,
			checks    integer,
			create_at timestamptz DEFAULT now(),
			update_at timestamptz,
			response  integer,
			UNIQUE(hostname)
		);
	
		CREATE TABLE IF NOT EXISTS links (
			hostname  text PRIMARY KEY,
			update_at timestamptz DEFAULT now(),
			UNIQUE(hostname)
		);

		INSERT INTO links (hostname) VALUES 
			('https://hidester.com/proxydata/php/data.php?mykey=data&offset=0&limit=1000&orderBy=latest_check&sortOrder=DESC&country=&port=&type=undefined&anonymity=undefined&ping=undefined&gproxy=2'),
			('http://gatherproxy.com/embed/'),
			('http://txt.proxyspy.net/proxy.txt'),
			('http://webanetlabs.net/publ/24'),
			('http://awmproxy.com/freeproxy.php'),
			('http://www.samair.ru/proxy/type-01.htm'),
			('https://www.us-proxy.org/'),
			('http://free-proxy-list.net/'),
			('http://www.proxynova.com/proxy-server-list/'),
			('http://proxyserverlist-24.blogspot.ru/'),
			('http://gatherproxy.com/'),
			('https://hidemy.name/ru/proxy-list/'),
			('https://hidemy.name/en/proxy-list/?type=hs&anon=34#list'),
			('https://free-proxy-list.com'),
			('https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=anonymous&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=2&connect_time%5B%5D=3&up_time=40&search=Search'),
			('http://www.idcloak.com/proxylist/free-proxy-servers-list.html'),
			('https://premproxy.com/list/'),
			('https://proxy-list.org/english/index.php'),
			('https://www.sslproxies.org/')
		ON CONFLICT (hostname) DO NOTHING;
	`)
}
