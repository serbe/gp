package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	flag.Parse()

	stringChan = make(chan string, 100000)

	t0 := time.Now()

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	grab()

	for _, s := range siteList {
		stringChan <- s
	}

	for len(stringChan) > 1 {
		// stringChan <- s
	}

	t1 := time.Now()
	fmt.Printf("%v\n", t1.Sub(t0))
}
