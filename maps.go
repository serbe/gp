package main

import (
	"sync"
)

type maps struct {
	m      sync.Mutex
	values map[string]bool
}

func newMaps() *maps {
	return &maps{values: make(map[string]bool)}
}

func (m *maps) get(s string) bool {
	m.m.Lock()
	defer m.m.Unlock()
	return m.values[s]
}

func (m *maps) set(s string, b bool) {
	m.m.Lock()
	defer m.m.Unlock()
	m.values[s] = b
	return
}

func (m *maps) len() int {
	m.m.Lock()
	defer m.m.Unlock()
	return len(m.values)
}
