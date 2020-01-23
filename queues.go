package main

import (
	"sync"
)

type reqQueue struct {
	sync.RWMutex
	nodes []req
	head  int
	tail  int
	cnt   int
}

func newReqQueue() reqQueue {
	return reqQueue{
		nodes: make([]req, 2),
	}
}

func (q *reqQueue) resize(n int) {
	nodes := make([]req, n)
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

func (q *reqQueue) put(req req) {
	q.Lock()
	if q.cnt == len(q.nodes) {
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = req
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
	q.Unlock()
}

func (q *reqQueue) get() (req, bool) {
	var req req
	q.Lock()
	if q.cnt == 0 {
		q.Unlock()
		return req, false
	}
	req = q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > 2 && q.cnt <= n {
		q.resize(n)
	}
	q.Unlock()
	return req, true
}

func (q *reqQueue) Cap() int {
	return cap(q.nodes)
}

func (q *reqQueue) Len() int {
	return q.cnt
}

type respQueue struct {
	sync.RWMutex
	nodes []resp
	head  int
	tail  int
	cnt   int
}

func newRespQueue() respQueue {
	return respQueue{
		nodes: make([]resp, 2),
	}
}

func (q *respQueue) resize(n int) {
	nodes := make([]resp, n)
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

func (q *respQueue) put(resp resp) {
	q.Lock()
	if q.cnt == len(q.nodes) {
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = resp
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
	q.Unlock()
}

func (q *respQueue) get() (resp, bool) {
	var resp resp
	q.Lock()
	if q.cnt == 0 {
		q.Unlock()
		return resp, false
	}
	resp = q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > 2 && q.cnt <= n {
		q.resize(n)
	}
	q.Unlock()
	return resp, true
}

func (q *respQueue) Cap() int {
	return cap(q.nodes)
}

func (q *respQueue) Len() int {
	return q.cnt
}
