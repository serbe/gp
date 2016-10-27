package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	tasks = make(chan string, 100000)
	crawlChan = make(chan string)
	finishTask = make(chan bool)

	quit := make(chan bool)

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	for i := 0; i < numWorkers; i++ {
		go worker(i, tasks, quit)
	}

	t0 := time.Now()

	for _, s := range siteList {
		urlList[s] = true
		addTask(s)
	}

Loop:
	for {
		select {
		case newWork := <-crawlChan:
			addTask(newWork)
		case <-finishTask:
			iter--
			if iter == 0 {
				for i := 0; i < numWorkers; i++ {
					quit <- true
				}
				break Loop
			}
		}
	}

	close(tasks)

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Parse %d urls\n", numUrls)
	fmt.Printf("%v second\n", t1.Sub(t0))
}
