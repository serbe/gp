package main

import (
	"sync"
)

// Queue - no race queue
type Queue struct {
	sync.RWMutex
	nodes []string
	head  int
	tail  int
	cnt   int
}

func newQueue() Queue {
	return Queue{
		nodes: make([]string, 2),
	}
}

func (q *Queue) resize(n int) {
	nodes := make([]string, n)
	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.tail])
	}

	q.tail = q.cnt % n
	q.head = 0
	q.nodes = nodes
}

func (q *Queue) put(value string) {
	q.Lock()
	if q.cnt == len(q.nodes) {
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = value
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
	q.Unlock()
}

func (q *Queue) get() (string, bool) {
	var value string
	q.Lock()
	if q.cnt == 0 {
		q.Unlock()
		return value, false
	}
	value = q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > 2 && q.cnt <= n {
		q.resize(n)
	}
	q.Unlock()
	return value, true
}

func (q *Queue) cap() int {
	return cap(q.nodes)
}

func (q *Queue) len() int {
	return q.cnt
}
