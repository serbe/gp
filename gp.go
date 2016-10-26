package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	flag.Parse()

	stringInChan = make(chan string)

	t0 := time.Now()

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	go grab()

	for _, s := range siteList {
		stringInChan <- s
		numUrls++
	}

	for numUrls > 0 {

	}

	t1 := time.Now()
	fmt.Printf("%v\n", t1.Sub(t0))
}
