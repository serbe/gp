package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
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

	os.Remove("gp.db")
	decompress("gp.zip")
	os.Remove("gp.zip")

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
			urls := grab(result)
			for _, u := range urls {
				p.Add(u, "")
			}
		}
		saveNewIP()
		saveLinks()
		fmt.Printf("Add %d ip adress\n", numIPs)
	}

	if checkProxy {
		var (
			totalIP    int64
			totalProxy int64
			anonProxy  int64
			err        error
		)
		targetURL := fmt.Sprintf("http://93.170.123.221:%d/", proxyPort)
		myIP, err = getExternalIP()
		if err == nil {
			month := time.Duration(30*60*24) * time.Minute
			startTime := time.Now()
			for _, v := range ips.values {
				if v.LastCheck.Sub(startTime) < time.Duration(v.ProxyChecks)*month || v.CreateAt.Sub(startTime) < time.Duration(v.ProxyChecks)*month {
					totalIP++
					p.Add(targetURL, makeAddress(v))
				}
			}
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
		checkProxyLoop:
			for {
				select {
				case result, ok := <-p.ResultChan:
					if ok {
						if result.Error == nil {
							proxy := check(result)
							ipString := proxy.Addr + ":" + proxy.Port
							ips.set(ipString, proxy)
							if proxy.isWork {
								totalProxy++
								if proxy.isAnon {
									anonProxy++
								}
								fmt.Println(ipString)
							}
						}
					} else {
						break checkProxyLoop
					}
				case <-c:
					break checkProxyLoop
				}
			}
			saveAllIP()
		}
		fmt.Printf("checked %d ip\n", totalIP)
		fmt.Printf("%d is good\n", totalProxy)
		fmt.Printf("%d is anon\n", anonProxy)
	}

	db.Sync()
	db.Close()

	compress("gp.db", "gp.zip")
	os.Remove("gp.db")
	os.Remove("gp.db.lock")

	endAppTime := time.Now()
	fmt.Printf("Total time: %v\n", endAppTime.Sub(startAppTime))
}
