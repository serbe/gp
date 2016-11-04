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
	ips := getListIP(body)
	urls := getListURL(host, body)

	saveIP(ips)

	resultChan <- urls
	return
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	tm := tasker.InitTasker(numWorkers, Grab)

	resultChan = make(chan []string)

	existsFile("ips.txt")
	urlList = newMaps()
	ipList = newMaps()
	getIPList()

	t0 := time.Now()

	for _, site := range siteList {
		urlList.set(site, true)
		tm.Queue(site)
	}

	func() {
		for {
			select {
			case result := <-resultChan:
				if len(result) > 0 {
					fmt.Printf("Get from chan %d urls\n", len(result))
				}
				for _, r := range result {
					tm.Queue(r)
				}
			case <-tm.Finish:
				return
			case <-time.After(time.Second):
				fmt.Printf("Queue len: %v num of running workers: %v\n", tm.QueueLen(), tm.RunningWorkers())
			}
		}
	}()

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Parse %d urls\n", urlList.len())
	fmt.Printf("%v second\n", t1.Sub(t0))
}
