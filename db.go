package main

import (
	"database/sql"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

// Proxy - proxy unit
type Proxy struct {
	Insert   bool          `sql:"-"           json:"-"`
	Update   bool          `sql:"-"           json:"-"`
	Hostname string        `sql:"hostname,pk" json:"hostname"`
	URL      *url.URL      `sql:"-"           json:"-"`
	Host     string        `sql:"host"        json:"-"`
	Port     string        `sql:"port"        json:"-"`
	IsWork   bool          `sql:"work"        json:"-"`
	IsAnon   bool          `sql:"anon"        json:"-"`
	Checks   int           `sql:"checks"      json:"-"`
	CreateAt time.Time     `sql:"create_at"   json:"-"`
	UpdateAt time.Time     `sql:"update_at"   json:"-"`
	Response time.Duration `sql:"response"    json:"-"`
}

func initDB() (*sql.DB, error) {
	return sql.Open("postgres", "user="+user+" password="+pass+" dbname="+dbname+" sslmode=disable")
}

func getAllProxy(db *sql.DB) *mapProxy {
	debugmsg("start getAllProxy")
	mProxy := newMapProxy()
	rows, err := db.Query("SELECT * FROM proxies")
	if err != nil {
		errmsg("getAllProxy Query select", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Proxy
		err := rows.Scan(
			&p.Hostname,
			&p.Host,
			&p.Port,
			&p.IsWork,
			&p.IsAnon,
			&p.Checks,
			&p.CreateAt,
			&p.UpdateAt,
			&p.Response,
		)
		if err != nil {
			errmsg("getAllProxy rows.Scan", err)
		}
		p.URL, err = url.Parse(p.Hostname)
		if err != nil {
			errmsg("getAllProxy url.Parse", err)
		}
		mProxy.set(p)
	}
	err = rows.Err()
	if err != nil {
		errmsg("getAllProxy rows.Err", err)
	}
	debugmsg("end getAllProxy, load proxy", len(mProxy.values))
	return mProxy
}

func getOldProxy(db *sql.DB) *mapProxy {
	debugmsg("start getOldProxy")
	mProxy := newMapProxy()
	rows, err := db.Query("SELECT * from proxies WHERE update_at < NOW() - (INTERVAL '3 days') * checks")
	if err != nil {
		errmsg("getOldProxy Query select", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Proxy
		err := rows.Scan(
			&p.Hostname,
			&p.Host,
			&p.Port,
			&p.IsWork,
			&p.IsAnon,
			&p.Checks,
			&p.CreateAt,
			&p.UpdateAt,
			&p.Response,
		)
		if err != nil {
			errmsg("getOldProxy rows.Scan", err)
		}
		p.URL, err = url.Parse(p.Hostname)
		if err != nil {
			errmsg("getOldProxy url.Parse", err)
		}
		mProxy.set(p)
	}
	err = rows.Err()
	if err != nil {
		errmsg("getOldProxy rows.Err", err)
	}
	debugmsg("end get100Proxy, load proxy", len(mProxy.values))
	return mProxy
}

func saveAllProxy(db *sql.DB, mProxy *mapProxy) {
	debugmsg("start saveAllProxy")
	var u, i int64
	for _, p := range mProxy.values {
		if p.Update {
			u++
			_, err := db.Exec(`
				UPDATE proxies SET
					host       = $2,
					port       = $3,
					work       = $4,
					anon       = $5,
					checks     = $6,
					create_at  = $7,
					update_at  = $8,
					response   = $9
				WHERE
					hostname = $1
			`,
				&p.Hostname,
				&p.Host,
				&p.Port,
				&p.IsWork,
				&p.IsAnon,
				&p.Checks,
				&p.CreateAt,
				&p.UpdateAt,
				&p.Response,
			)
			if err != nil {
				errmsg("saveAllProxy Update", err)
			}
		}
		if p.Insert {
			i++
			_, err := db.Exec(`
				INSERT INTO proxies (
					hostname,
					host,    
					port,    
					work,  
					anon,  
					checks,  
					create_at,
					update_at,
					response
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7,
					$8,
					$9
				)
			`,
				&p.Hostname,
				&p.Host,
				&p.Port,
				&p.IsWork,
				&p.IsAnon,
				&p.Checks,
				&p.CreateAt,
				&p.UpdateAt,
				&p.Response,
			)
			if err != nil {
				errmsg("saveAllLinks Insert", err)
			}
		}
	}
	debugmsg("update proxy", u)
	debugmsg("insert proxy", i)
	debugmsg("end getAllProxy")
}

func updateAllProxy(db *sql.DB, mProxy *mapProxy) {
	debugmsg("start updateAllProxy")
	for _, p := range mProxy.values {
		_, err := db.Exec(`
			UPDATE proxies SET
				host       = $2,
				port       = $3,
				work       = $4,
				anon       = $5,
				checks     = $6,
				create_at  = $7,
				update_at  = $8,
				response   = $9
			WHERE
				hostname = $1
		`,
			&p.Hostname,
			&p.Host,
			&p.Port,
			&p.IsWork,
			&p.IsAnon,
			&p.Checks,
			&p.CreateAt,
			&p.UpdateAt,
			&p.Response,
		)
		if err != nil {
			errmsg("updateAllProxy update", err)
		}
	}
	debugmsg("end updateAllProxy, update proxy", len(mProxy.values))
}

func getAllLinks(db *sql.DB) *mapLink {
	debugmsg("start getAllLinks")
	mLink := newMapLink()
	rows, err := db.Query("SELECT * FROM links")
	if err != nil {
		errmsg("getAllLinks Query select", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Link
		err := rows.Scan(
			&p.Hostname,
			&p.UpdateAt,
		)
		if err != nil {
			errmsg("getAllLinks rows.Scan", err)
		}
		mLink.set(p)
	}
	err = rows.Err()
	if err != nil {
		errmsg("getAllProxy rows.Err", err)
	}
	debugmsg("end getAllLinks, load links", len(mLink.values))
	return mLink
}

func saveAllLinks(db *sql.DB, mL *mapLink) {
	debugmsg("start saveAllLinks")
	var (
		u, i int64
	)
	for _, l := range mL.values {
		if l.Insert {
			i++
			_, err := db.Exec(`
				INSERT INTO links (
					hostname,
					update_at
				) VALUES (
					$1,
					$2
				)
			`,
				&l.Hostname,
				&l.UpdateAt,
			)
			if err != nil {
				errmsg("saveAllLinks Insert", err)
			}
		} else {
			u++
			_, err := db.Exec(`
				UPDATE links SET
					update_at = $2
				WHERE
					hostname = $1
			`,
				&l.Hostname,
				&l.UpdateAt,
			)
			if err != nil {
				errmsg("saveAllLinks Update", err)
			}
		}
	}
	debugmsg("update links", u)
	debugmsg("insert links", i)
	debugmsg("end saveAllLinks")
}

// func makeTables() {
// 	db.ExecOne(`
// 		CREATE TABLE IF NOT EXISTS proxies (
// 			hostname  text PRIMARY KEY,
// 			host      text,
// 			port      text,
// 			work      boolean,
// 			anon      boolean,
// 			checks    integer,
// 			create_at timestamptz DEFAULT now(),
// 			update_at timestamptz,
// 			response  integer,
// 			UNIQUE(hostname)
// 		);

// 		CREATE TABLE IF NOT EXISTS links (
// 			hostname  text PRIMARY KEY,
// 			update_at timestamptz DEFAULT now(),
// 			UNIQUE(hostname)
// 		);

// 		INSERT INTO links (hostname) VALUES
// 			('https://hidester.com/proxydata/php/data.php?mykey=data&offset=0&limit=1000&orderBy=latest_check&sortOrder=DESC&country=&port=&type=undefined&anonymity=undefined&ping=undefined&gproxy=2'),
// 			('http://gatherproxy.com/embed/'),
// 			('http://txt.proxyspy.net/proxy.txt'),
// 			('http://webanetlabs.net/publ/24'),
// 			('http://awmproxy.com/freeproxy.php'),
// 			('http://www.samair.ru/proxy/type-01.htm'),
// 			('https://www.us-proxy.org/'),
// 			('http://free-proxy-list.net/'),
// 			('http://www.proxynova.com/proxy-server-list/'),
// 			('http://proxyserverlist-24.blogspot.ru/'),
// 			('http://gatherproxy.com/'),
// 			('https://hidemy.name/ru/proxy-list/'),
// 			('https://hidemy.name/en/proxy-list/?type=hs&anon=34#list'),
// 			('https://free-proxy-list.com'),
// 			('https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=anonymous&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=2&connect_time%5B%5D=3&up_time=40&search=Search'),
// 			('http://www.idcloak.com/proxylist/free-proxy-servers-list.html'),
// 			('https://premproxy.com/list/'),
// 			('https://proxy-list.org/english/index.php'),
// 			('https://www.sslproxies.org/')
// 		ON CONFLICT (hostname) DO NOTHING;
// 	`)
// }
