package main

import (
	"flag"
	"fmt"
	"sync"
)

var (
	jobs  chan string
	mutex = &sync.Mutex{}
)

func main() {
	workers := flag.Int("w", 5, "Num of workers")

	flag.Parse()

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	jobs = make(chan string)
	results := make(chan int)

	for w := 1; w <= *workers; w++ {
		go worker(w, jobs, results)
	}

	for _, u := range siteList {
		jobs <- u
	}

	for r := range results {
		if r > 0 {
			fmt.Println("Add ", r, " ip")
		}
	}
}
