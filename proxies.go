package main

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/serbe/pool"
)

// Proxy - proxy unit
type Proxy struct {
	Insert   bool          `sql:"-"         json:"-"`
	Update   bool          `sql:"-"         json:"-"`
	Hostname string        `sql:"hostname"  json:"hostname"`
	URL      *url.URL      `sql:"-"         json:"-"`
	Host     string        `sql:"host"      json:"-"`
	Port     string        `sql:"port"      json:"-"`
	IsWork   bool          `sql:"work"      json:"-"`
	IsAnon   bool          `sql:"anon"      json:"-"`
	Checks   int           `sql:"checks"    json:"-"`
	CreateAt time.Time     `sql:"create_at" json:"-"`
	UpdateAt time.Time     `sql:"update_at" json:"-"`
	Response time.Duration `sql:"response"  json:"-"`
}

type mapProxy struct {
	m      sync.RWMutex
	values map[string]Proxy
}

func newMapProxy() *mapProxy {
	return &mapProxy{values: make(map[string]Proxy)}
}

func (mProxy *mapProxy) get(hostname string) Proxy {
	mProxy.m.RLock()
	proxy := mProxy.values[hostname]
	mProxy.m.RUnlock()
	return proxy
}

func (mProxy *mapProxy) set(proxy Proxy) {
	mProxy.m.Lock()
	mProxy.values[proxy.Hostname] = proxy
	mProxy.m.Unlock()
}

func newProxy(host, port string, ssl bool) (Proxy, error) {
	var (
		proxy  Proxy
		schema string
	)
	if ssl {
		schema = "https://"
	} else {
		schema = "http://"
	}
	hostname := schema + host + ":" + port
	URL, err := url.Parse(hostname)
	if err != nil {
		return proxy, err
	}
	proxy.Hostname = hostname
	proxy.Insert = true
	proxy.URL = URL
	proxy.Host = host
	proxy.Port = port
	proxy.CreateAt = time.Now()
	return proxy, err
}

func setProxy(host string, port string, ssl bool) {
	proxy, err := newProxy(host, port, false)
	if err == nil {
		if !mP.existProxy(proxy.Hostname) {
			numIPs++
			mP.set(proxy)
		}
	}
}

func (mProxy *mapProxy) existProxy(hostname string) bool {
	mProxy.m.RLock()
	_, ok := mProxy.values[hostname]
	mProxy.m.RUnlock()
	return ok
}

func taskToProxy(task pool.Task) Proxy {
	proxy := Proxy{
		URL:      task.Proxy,
		UpdateAt: time.Now(),
	}
	if task.Error == nil {
		strBody := string(task.Body)
		if reRemoteIP.Match(task.Body) && !strings.Contains(strBody, myIP) {
			proxy.IsWork = true
			proxy.Checks = 0
			if strings.Count(strBody, "<p>") == 1 {
				proxy.IsAnon = true
			}
			return proxy
		}
	}
	proxy.Checks++
	return proxy
}
