package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/serbe/pool"
)

func main() {
	var (
		findProxy  = true
		checkProxy = false
		backup     = false
		// err        error
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

	p := pool.New(numWorkers)
	p.SetTimeout(timeout)

	links = getAllLinks()
	ips = getAllIP()

	if findProxy {
		for _, site := range siteList {
			links.set(site)
			p.Add(site, "")
		}
		for result := range p.ResultChan {
			urls := grab(result.Address, result.Body)
			for _, u := range urls {
				p.Add(u, "")
			}
		}
		saveNewIP()
		saveLinks()
		fmt.Printf("Add %d ip adress\n", numIPs)
	}

	// if checkProxy {
	// 	var (
	// 		totalIP    int64
	// 		totalProxy int64
	// 		anonProxy  int64
	// 	)
	// 	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", proxyPort)
	// 	myIP, err := getExternalIP()
	// 	if err == nil {
	// 		month := time.Duration(30*60*24) * time.Minute
	// 		timeNow := time.Now()
	// 		for _, v := range ips.values {
	// 			if v.LastCheck.Sub(timeNow) < time.Duration(v.ProxyChecks)*month || v.CreateAt.Sub(timeNow) < time.Duration(v.ProxyChecks)*month {
	// 				totalIP++
	// 				p.Add(targetURL, makeAddress(v))
	// 			}
	// 		}
	// 		c := make(chan os.Signal, 1)
	// 		signal.Notify(c, os.Interrupt)
	// 	checkProxyLoop:
	// 		for {
	// 			select {
	// 			case result, ok := <-p.ResultChan:
	// 				if ok {
	// 					if task.Result != nil {
	// 						ip := task.Result.(ipType)
	// 						ipString := ip.Addr + ":" + ip.Port
	// 						ips.set(ipString, ip)
	// 						if ip.isWork {
	// 							totalProxy++
	// 							if ip.isAnon {
	// 								anonProxy++
	// 							}
	// 							fmt.Println(ipString)
	// 						}
	// 					}
	// 				} else {
	// 					break checkProxyLoop
	// 				}
	// 			case <-c:
	// 				p.Quit()
	// 				break checkProxyLoop
	// 			case <-time.After(time.Duration(100) * time.Millisecond):
	// 				if p.Done() {
	// 					break checkProxyLoop
	// 				}
	// 			}
	// 		}
	// 		saveAllIP()
	// 	}
	// 	fmt.Printf("checked %d ip\n", totalIP)
	// 	fmt.Printf("%d is good\n", totalProxy)
	// 	fmt.Printf("%d is anon\n", anonProxy)
	// }

	db.Sync()
	db.Close()

	compress("gp.db", "gp.gz")
	os.Remove("gp.db")
	os.Remove("gp.db.lock")

	endAppTime := time.Now()
	fmt.Printf("Total time: %v\n", endAppTime.Sub(startAppTime))
}
