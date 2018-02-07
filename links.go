package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type mapLink struct {
	sync.RWMutex
	values map[string]Link
}

func newMapLink() *mapLink {
	return &mapLink{values: make(map[string]Link)}
}

func (mLink *mapLink) get(hostname string) Link {
	mLink.RLock()
	link := mLink.values[hostname]
	mLink.RUnlock()
	return link
}

func (mLink *mapLink) set(link Link) {
	mLink.Lock()
	mLink.values[link.Hostname] = link
	mLink.Unlock()
}

func (mLink *mapLink) update(hostname string) {
	mLink.Lock()
	link := mLink.values[hostname]
	link.Update = true
	link.UpdateAt = time.Now()
	mLink.values[hostname] = link
	mLink.Unlock()
}

func (mLink *mapLink) newLink(hostname string) Link {
	var link Link
	link.Hostname = hostname
	return link
}

func (mLink *mapLink) existLink(hostname string) bool {
	mLink.RLock()
	_, ok := mLink.values[hostname]
	mLink.RUnlock()
	return ok
}

func getMapLink(db *sql.DB) *mapLink {
	var mL *mapLink
	if testLink != "" {
		mL = newMapLink()
		link := mL.newLink(testLink)
		link.Iterate = true
		mL.set(link)
		log.Println(link)
	} else if addLink != "" {
		mL = newMapLink()
		link := mL.newLink(addLink)
		link.Insert = true
		link.Iterate = true
		mL.set(link)
		log.Println(link)
	} else {
		mL = getAllLinks(db)
	}
	return mL
}
