package main

import (
	"sync"
	"time"
)

// Link - link unit
type Link struct {
	Insert   bool      `sql:"-"         json:"-"`
	Update   bool      `sql:"-"         json:"-"`
	Hostname string    `sql:"hostname"  json:"hostname"`
	UpdateAt time.Time `sql:"update_at" json:"-"`
}

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
	link.UpdateAt = time.Now()
	mLink.values[hostname] = link
	mLink.m.Unlock()
}

func (mLink *mapLink) newLink(hostname string) Link {
	var link Link
	link.Hostname = hostname
	return link
}

func (mLink *mapLink) isOldLink(hostname string) bool {
	link := mLink.get(hostname)
	return time.Since(link.UpdateAt) > time.Duration(10)*time.Minute
}

func (mLink *mapLink) oldLinks() []Link {
	debugmsg("create list of old Link")
	var links []Link
	for _, link := range mLink.values {
		if time.Since(link.UpdateAt) > time.Duration(10)*time.Minute {
			links = append(links, link)
		}
	}
	debugmsg("old links len is", len(links))
	return links
}

func (mLink *mapLink) existLink(hostname string) bool {
	mLink.m.RLock()
	_, ok := mLink.values[hostname]
	mLink.m.RUnlock()
	return ok
}
