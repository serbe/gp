package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/serbe/tasker"
)

var (
	resultChan chan string
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
	oldNumIP := numIPs
	getListIP(body)
	fmt.Printf("Find %d new ip address\n", numIPs-oldNumIP)
	getListURL(host, body)

	return
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	initDB()
	defer db.Close()

	tm := tasker.InitTasker(numWorkers, Grab)

	resultChan = make(chan string)

	t0 := time.Now()

	for _, site := range siteList {
		saveLink(site)
		tm.Queue(site)
	}

	func() {
		for {
			select {
			case host := <-resultChan:
				tm.Queue(host)
			case <-tm.Finish:
				fmt.Println("finish")
				return
			case <-time.After(time.Second):
				fmt.Printf("Queue len: %v num of running workers: %v\n", tm.QueueLen(), tm.RunningWorkers())
			}
		}
	}()

	t1 := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("%v second\n", t1.Sub(t0))
}
