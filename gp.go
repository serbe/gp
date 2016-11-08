package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/serbe/tasker"
)

var (
	resultChan   chan string
	numWorkers   = 5
	ips          *mapsIP
	links        *mapsLink
	startAppTime time.Time
)

// Grab - parse url
func Grab(hostURL interface{}) {
	fmt.Printf("Start grab %s\n", hostURL.(string))
	host := hostURL.(string)
	body, err := fetch(host)
	if err != nil {
		return
	}

	body = cleanBody(body)
	oldNumIP := numIPs
	getListIP(body)
	if numIPs-oldNumIP > 0 {
		fmt.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, hostURL.(string))
	}
	getListURL(host, body)

	return
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	initDB()
	defer db.Close()

	startAppTime = time.Now()

	tm := tasker.InitTasker(numWorkers, Grab)
	resultChan = make(chan string)
	links = newMapsLink()
	ips = getAllIP()

	for _, site := range siteList {
		links.set(site, true)
		tm.Queue(site)
	}

	func() {
		// var (
		// 	temptime  time.Time
		// 	startWait bool
		// )
		for {
			select {
			case host := <-resultChan:
				tm.Queue(host)
			case <-*tm.Finish:
				fmt.Println("finish")
				return
				// case <-time.After(time.Second):
				// 	if tm.QueueLen() == 0 {
				// 		if !startWait {
				// 			startWait = true
				// 			temptime = time.Now()
				// 			fmt.Println("Wait 30 secont to finish all tasks")
				// 		} else {
				// 			temptime2 := time.Now()
				// 			if temptime2.Sub(temptime) > time.Duration(30*time.Second) {
				// 				*tm.Finish <- true
				// 				break
				// 			}
				// 			fmt.Printf("Len queue %d, have task %d workers\n", tm.QueueLen(), tm.RunningWorkers())
				// 		}
				// 	}
			}
		}
	}()

	saveNewIP()

	endAppTime := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("%v second\n", endAppTime.Sub(startAppTime))
}
