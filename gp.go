package main

import (
	"flag"
	"fmt"
	"log"
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
	)

	flag.IntVar(&numWorkers, "w", numWorkers, "number of workers")
	flag.IntVar(&timeout, "t", timeout, "timeout")
	flag.BoolVar(&findProxy, "f", findProxy, "find new proxy")
	flag.BoolVar(&checkProxy, "c", checkProxy, "check proxy")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
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

	if findProxy {
		p := pool.New(numWorkers)
		p.SetTimeout(timeout)
		links = getAllLinks()
		ips = getAllIP()
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
		log.Printf("Add %d ip adress\n", numIPs)
	}

	if checkProxy {
		var (
			totalIP    int64
			totalProxy int64
			anonProxy  int64
			err        error
		)
		ips = getAllIP()
		p := pool.New(numWorkers)
		p.SetTimeout(timeout)
		targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
		myIP, err = getExternalIP()
		if err == nil {
			week := time.Duration(60*24*7) * time.Minute
			startTime := time.Now()
			for _, v := range ips.values {
				if (v.LastCheck == time.Time{} || v.LastCheck != time.Time{} && startTime.Sub(v.LastCheck) > time.Duration(v.ProxyChecks)*week) {
					totalIP++
					p.Add(targetURL, makeAddress(v))
				}
			}
			log.Println("Start check", totalIP, "proxyes")
			if totalIP > 0 {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				var checked int
			checkProxyLoop:
				for {
					select {
					case result, ok := <-p.ResultChan:
						checked++
						if ok {
							proxy := check(result)
							proxy.Response = result.ResponceTime
							ipString := proxy.Addr + ":" + proxy.Port
							ips.set(ipString, proxy)
							if proxy.isWork {
								log.Printf("%d/%d %-15v %-5v %-10v anon=%v\n", checked, totalIP, result.Proxy.Hostname(), result.Proxy.Port(), result.ResponceTime, proxy.isAnon)
								totalProxy++
								if proxy.isAnon {
									anonProxy++
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
		}
		log.Printf("checked %d ip\n", totalIP)
		log.Printf("%d is good\n", totalProxy)
		log.Printf("%d is anon\n", anonProxy)
	}
	db.Sync()
	db.Close()

	compress("gp.db", "gp.zip")
	os.Remove("gp.db")
	os.Remove("gp.db.lock")

	endAppTime := time.Now()
	log.Printf("Total time: %v\n", endAppTime.Sub(startAppTime))
}
