package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/serbe/gopool"
)

func main() {
	var (
		findProxy  = false
		checkProxy = true
		backup     = true
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

	loopFind:
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
				break loopFind
			}
		}
		saveNewIP()
		saveLinks()
	}

	if checkProxy {
		go func() {
			month := time.Duration(30*60*24) * time.Minute
			timeNow := time.Now()
			for _, v := range ips.values {
				if v.LastCheck.Sub(timeNow) < time.Duration(v.ProxyChecks)*month || v.CreateAt.Sub(timeNow) < time.Duration(v.ProxyChecks)*month {
					tm.Add(check, v)
				}
			}
		}()
	loopCheck:
		for {
			task := tm.GetTask()
			if task.Result != nil {
				ip := task.Result.(ipType)
				ipString := ip.Addr + ":" + ip.Port
				ips.set(ipString, ip)
				if ip.isWork {
					fmt.Println(ipString)
				}
			}
			added, running, completed := tm.Status()
			if running == 0 && added > 0 && added == completed {
				break loopCheck
			}
		}
		saveNewIP()
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
