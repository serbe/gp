package main

import (
	"flag"
	"fmt"
	"time"
)

// Grab - parse url
func Grab(hostURL interface{}) {
	host := hostURL.(string)
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

	// tm.StartWorkers()

	crawlChan = make(chan string)
	finishTask = make(chan bool)

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	t0 := time.Now()

	for _, s := range siteList {
		urlList[s] = true
		tm.AddTask(s)
	}

	close(tasks)

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Parse %d urls\n", numUrls)
	fmt.Printf("%v second\n", t1.Sub(t0))
}
