package main

import (
	"time"
)

// Task - result from crawl
type Task struct {
	Proxy    string
	Body     []byte
	Response time.Duration
	Error    error
}
