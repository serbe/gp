package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	tasks := make(chan string, 1000000)
	crawlChan = make(chan string)

	quit := make(chan bool)

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	for i := 0; i < numWorkers; i++ {
		numWorkWorkers++
		go worker(tasks, quit)
	}

	t0 := time.Now()

	for _, s := range siteList {
		urlList[s] = true
		tasks <- s
	}

	for numWorkWorkers > 0 {
		select {
		case newWork := <-crawlChan:
			tasks <- newWork
		}
	}

	close(tasks)

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("%v second\n", t1.Sub(t0))
}
