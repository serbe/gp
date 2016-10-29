package main

import (
	"flag"
	"fmt"
	"time"
)

// Grab - parse url
func Grab(host string) {
	body, err := fetch(host)
	if err != nil {
		finishTask <- true
		return
	}

	body = cleanBody(body)

	ips := getListIP(body)

	urls := getListURL(host, body)

	saveIP(ips)

	for _, item := range urls {
		if !urlList[item] {
			urlList[item] = true
			crawlChan <- item
		}
	}
	finishTask <- true
	return
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")

	flag.Parse()

	tm := InitTaskMaster(numWorkers, Grab)

	tm.Tasks = make(chan interface{}, 100000)
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
		tm.AddTask(s)
	}

Loop:
	for {
		select {
		case newWork := <-crawlChan:
			tm.AddTask(newWork)
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
