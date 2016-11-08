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
	links = getAllLinks()
	ips = getAllIP()

	for _, site := range siteList {
		links.set(site)
		tm.Queue(site)
	}

	func() {
		for {
			select {
			case host := <-resultChan:
				tm.Queue(host)
			case <-*tm.Finish:
				fmt.Println("finish")
				return
				// case <-time.After(time.Duration(5) * time.Second):
				// 	fmt.Printf("queue len %d worked %d\n", tm.QueueLen(), tm.RunningWorkers())
				// 	infoWS := tm.GetWorkersInfo()
				// 	for _, info := range infoWS {
				// 		fmt.Printf("w.id %d start %v task %v\n", info.ID, info.StartTime, info.Task)
				// 	}
			}
		}
	}()

	saveNewIP()
	saveLinks()

	endAppTime := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("%v second\n", endAppTime.Sub(startAppTime))
}
