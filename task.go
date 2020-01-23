package main

import (
	"time"
)

// Task - result from crawl
type Task struct {
	ID       int64
	Proxy    string
	Body     []byte
	Response time.Duration
	Error    error
}
