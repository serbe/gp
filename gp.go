package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/gopool"
)

func main() {
	var (
		findProxy  = true
		checkProxy = false
		backup     = false
		err        error
	)

	flag.IntVar(&numWorkers, "w", numWorkers, "number of workers")
	flag.IntVar(&timeout, "t", timeout, "timeout")
	flag.BoolVar(&findProxy, "f", findProxy, "find new proxy")
	flag.BoolVar(&checkProxy, "c", checkProxy, "check proxy")
	flag.IntVar(&proxyPort, "p", proxyPort, "proxy port")
	flag.BoolVar(&backup, "b", backup, "backup database")

	flag.Parse()

	if backup {
		backupBase()
	}

	decompress("gp.gz", "gp.db")
	os.Remove("gp.gz")

	initDB()
	defer db.Close()

	startAppTime = time.Now()

	tm := gopool.New(numWorkers)
	tm.Run()

	links = getAllLinks()
	ips = getAllIP()

	if findProxy {
		for _, site := range siteList {
			links.set(site)
			tm.Add(grab, site)
		}
		r := tm.ResultChan(true)
	getResultFindLoop:
		for {
			select {
			case task := <-*r:
				if task.Result != nil {
					urls := task.Result.([]string)
					for _, u := range urls {
						tm.Add(grab, u)
					}
				}
			case <-time.After(time.Duration(100) * time.Millisecond):
				if tm.Done() {
					break getResultFindLoop
				}
			}
		}
		tm.ResultChan(false)
		saveNewIP()
		saveLinks()
	}

	if checkProxy {
		myIP, err = getExternalIP()
		if err == nil {
			month := time.Duration(30*60*24) * time.Minute
			timeNow := time.Now()
			for _, v := range ips.values {
				if v.LastCheck.Sub(timeNow) < time.Duration(v.ProxyChecks)*month || v.CreateAt.Sub(timeNow) < time.Duration(v.ProxyChecks)*month {
					tm.Add(check, v)
				}
			}
			r := tm.ResultChan(true)
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
		getResultCheckLoop:
			for {
				select {
				case task := <-*r:
					if task.Result != nil {
						ip := task.Result.(ipType)
						ipString := ip.Addr + ":" + ip.Port
						ips.set(ipString, ip)
						if ip.isWork {
							fmt.Println(ipString)
						}
					}
				case <-c:
					tm.Quit()
					break getResultCheckLoop
				case <-time.After(time.Duration(100) * time.Millisecond):
					if tm.Done() {
						break getResultCheckLoop
					}
				}
			}
			saveAllIP()
			tm.ResultChan(false)
		}
	}

	db.Sync()
	db.Close()

	compress("gp.db", "gp.gz")
	os.Remove("gp.db")
	os.Remove("gp.db.lock")

	endAppTime := time.Now()
	fmt.Printf("Add %d ip adress\n", numIPs)
	fmt.Printf("Total time: %v second\n", endAppTime.Sub(startAppTime))
}
