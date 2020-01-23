package main

import (
	"time"
)

type resp struct {
	ID       int64
	Hostname string
	Proxy    string
	Body     []byte
	Response time.Duration
	Error    error
}

type req struct {
	ID       int64
	Hostname string
	Proxy    string
}
