package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var (
	jobs  chan string
	mutex = &sync.Mutex{}
)

func main() {
	depth := flag.Int("d", 5, "Num of depth")

	flag.Parse()

	t0 := time.Now()

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	finished := make(chan bool)

	for _, u := range siteList {
		go crawl(u, *depth, finished)
	}

	<-finished

	t1 := time.Now()
	fmt.Printf("%v\n", t1.Sub(t0))
}
