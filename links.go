package main

import (
	"sync"
	"time"
)

type mapLink struct {
	m      sync.RWMutex
	values map[string]Link
}

func newMapLink() *mapLink {
	return &mapLink{values: make(map[string]Link)}
}

func (mLink *mapLink) get(hostname string) Link {
	mLink.m.RLock()
	link := mLink.values[hostname]
	mLink.m.RUnlock()
	return link
}

func (mLink *mapLink) set(link Link) {
	mLink.m.Lock()
	mLink.values[link.Hostname] = link
	mLink.m.Unlock()
}

func (mLink *mapLink) update(hostname string) {
	mLink.m.Lock()
	link := mLink.values[hostname]
	link.Update = true
	link.UpdateAt = time.Now()
	mLink.values[hostname] = link
	mLink.m.Unlock()
}

func (mLink *mapLink) newLink(hostname string) Link {
	var link Link
	link.Hostname = hostname
	return link
}

func (mLink *mapLink) existLink(hostname string) bool {
	mLink.m.RLock()
	_, ok := mLink.values[hostname]
	mLink.m.RUnlock()
	return ok
}
