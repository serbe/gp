package main

import (
	"sync"
	"time"
)

type mapProxy struct {
	m      sync.RWMutex
	values map[string]Proxy
}

type mapLink struct {
	m      sync.RWMutex
	values map[string]Link
}

func newMapProxy() *mapProxy {
	return &mapProxy{values: make(map[string]Proxy)}
}

// func (mProxy *mapProxy) get(hostname string) Proxy {
// 	mProxy.m.RLock()
// 	proxy := mProxy.values[hostname]
// 	mProxy.m.RUnlock()
// 	return proxy
// }

func (mProxy *mapProxy) set(proxy Proxy) {
	mProxy.m.Lock()
	mProxy.values[proxy.URL.Hostname()] = proxy
	mProxy.m.Unlock()
}

func newMapLink() *mapLink {
	return &mapLink{values: make(map[string]Link)}
}

func (mLink *mapLink) get(s string) Link {
	mLink.m.RLock()
	link := mLink.values[s]
	mLink.m.RUnlock()
	return link
}

func (mLink *mapLink) set(s string) {
	mLink.m.Lock()
	var value Link
	value.Host = s
	value.CheckAt = time.Now()
	mLink.values[s] = value
	mLink.m.Unlock()
}
