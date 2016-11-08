package main

import (
	"sync"
)

type mapsIP struct {
	m      sync.Mutex
	values map[string]ipType
}

type mapsLink struct {
	m      sync.Mutex
	values map[string]bool
}

func newMapsIP() *mapsIP {
	return &mapsIP{values: make(map[string]ipType)}
}

func (mIP *mapsIP) get(s string) ipType {
	mIP.m.Lock()
	defer mIP.m.Unlock()
	return mIP.values[s]
}

func (mIP *mapsIP) set(s string, value ipType) {
	mIP.m.Lock()
	defer mIP.m.Unlock()
	mIP.values[s] = value
	return
}

func (mIP *mapsIP) len() int {
	mIP.m.Lock()
	defer mIP.m.Unlock()
	return len(mIP.values)
}

func newMapsLink() *mapsLink {
	return &mapsLink{values: make(map[string]bool)}
}

func (mLink *mapsLink) get(s string) bool {
	mLink.m.Lock()
	defer mLink.m.Unlock()
	return mLink.values[s]
}

func (mLink *mapsLink) set(s string, value bool) {
	mLink.m.Lock()
	defer mLink.m.Unlock()
	mLink.values[s] = value
	return
}

func (mLink *mapsLink) len() int {
	mLink.m.Lock()
	defer mLink.m.Unlock()
	return len(mLink.values)
}
