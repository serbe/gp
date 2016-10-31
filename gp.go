package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/serbe/tasker"
)

var (
	resultChan chan []string
	numWorkers = 5
)

// Grab - parse url
func Grab(hostURL interface{}) {
	resultChan = make(chan []string)
	host := hostURL.(string)
	body, err := fetch(host)
	if err != nil {
		return
	}

	body = cleanBody(body)
	getListIP(body)
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

	tm := tasker.InitTasker(numWorkers, Grab)

	crawlChan = make(chan string)
	finishTask = make(chan bool)

	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()

	t0 := time.Now()

	for _, site := range siteList {
		urlList[site] = true
		tm.Work <- site
	}

	func() {
		for {
			select {
			case result := <-resultChan:
				for i := range result {
					tm.Work <- result[i]
				}
			case <-tm.Quit:
				return
			}
		}
	}()

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Parse %d urls\n", numUrls)
	fmt.Printf("%v second\n", t1.Sub(t0))
}
