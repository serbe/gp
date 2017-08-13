package main

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/serbe/pool"
)

type mapProxy struct {
	m      sync.RWMutex
	values map[string]Proxy
}

func newMapProxy() *mapProxy {
	return &mapProxy{values: make(map[string]Proxy)}
}

func (mProxy *mapProxy) set(proxy Proxy) {
	mProxy.m.Lock()
	mProxy.values[proxy.Hostname] = proxy
	mProxy.m.Unlock()
}

func (mProxy *mapProxy) get(hostname string) (Proxy, bool) {
	mProxy.m.Lock()
	proxy, ok := mProxy.values[hostname]
	mProxy.m.Unlock()
	return proxy, ok
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

func (mProxy *mapProxy) setProxy(host string, port string, ssl bool) {
	proxy, err := newProxy(host, port, false)
	if err == nil {
		if !mProxy.existProxy(proxy.Hostname) {
			numIPs++
			mProxy.set(proxy)
		}
	}
}

func (mProxy *mapProxy) existProxy(hostname string) bool {
	mProxy.m.RLock()
	_, ok := mProxy.values[hostname]
	mProxy.m.RUnlock()
	return ok
}

func (mProxy *mapProxy) taskToProxy(task pool.Task) (Proxy, bool) {
	proxy, ok := mProxy.get(task.Proxy.String())
	if ok {
		proxy.Update = true
		proxy.UpdateAt = time.Now()
		proxy.Response = task.ResponceTime
		strBody := string(task.Body)
		if reRemoteIP.Match(task.Body) && !strings.Contains(strBody, myIP) {
			proxy.IsWork = true
			proxy.Checks = 0
			if strings.Count(strBody, "<p>") == 1 {
				proxy.IsAnon = true
			}
			return proxy, ok
		}
		proxy.Checks++
	}
	return proxy, ok
}

func proxyIsOld(proxy Proxy) bool {
	return proxy.UpdateAt == time.Time{} ||
		proxy.UpdateAt != time.Time{} &&
			time.Since(proxy.UpdateAt) > time.Duration(proxy.Checks)*time.Duration(60*24*7)*time.Minute
}
