package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/serbe/gopool"
)

var (
	numWorkers   = 5
	ips          *mapsIP
	links        *mapsLink
	startAppTime time.Time
)

func grab(args ...interface{}) interface{} {
	var urls []string
	host := args[0].(string)
	fmt.Printf("Start grab %s\n", host)
	body, err := fetch(host)
	if err != nil {
		return urls
	}

	body = cleanBody(body)
	oldNumIP := numIPs
	getListIP(body)
	if numIPs-oldNumIP > 0 {
		fmt.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, host)
	}
	urls = getListURL(host, body)

	return urls
}

func main() {
	flag.IntVar(&numWorkers, "w", numWorkers, "количество рабочих")
	flag.Parse()

	decompress("ips.gz", "ips.db")
	os.Remove("ips.gz")

	initDB()
	defer db.Close()

	startAppTime = time.Now()

	tm := gopool.New(numWorkers)
	tm.Run()

	links = getAllLinks()
	ips = getAllIP()

	for _, site := range siteList {
		links.set(site)
		tm.Add(grab, site)
	}

loop:
	for {
		task := tm.GetTask()

		if task.Result != nil {
			urls := task.Result.([]string)
			for _, u := range urls {
				tm.Add(grab, u)
			}
		}
		added, running, completed := tm.Status()
		if running == 0 && added > 0 && added == completed {
			break loop
		}
	}

	saveNewIP()
	saveLinks()

	db.Sync()
	db.Close()

	compress("ips.db", "ips.gz")
	os.Remove("ips.db")

	endAppTime := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Total time: %v second\n", endAppTime.Sub(startAppTime))
}
