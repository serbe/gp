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
	fmt.Printf("Start grab %s url\n", hostURL.(string))
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
			mutex.Lock()
			urlList[item] = true
			mutex.Unlock()
		}
	}
	resultChan <- urls
	return
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	tm := tasker.InitTasker(numWorkers, Grab)

	resultChan = make(chan []string)

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
				fmt.Printf("Get from chan %d urls\n", len(result))
				for _, r := range result {
					mutex.Lock()
					if !urlList[r] {
						tm.Work <- r
					}
					mutex.Unlock()
				}
			case <-tm.Quit:
				return
			}
		}
	}()

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Parse %d urls\n", len(urlList))
	fmt.Printf("%v second\n", t1.Sub(t0))
}
