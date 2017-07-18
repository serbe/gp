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
		server     = false
	)

	flag.IntVar(&numWorkers, "w", numWorkers, "number of workers")
	flag.IntVar(&timeout, "t", timeout, "timeout")
	flag.BoolVar(&findProxy, "f", findProxy, "find new proxy")
	flag.BoolVar(&checkProxy, "c", checkProxy, "check proxy")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.BoolVar(&backup, "b", backup, "backup database")
	flag.BoolVar(&server, "s", server, "start server")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")

	flag.Parse()

	if backup {
		err := backupBase()
		if err != nil {
			errmsg("backupBase", err)
		}
	}

	err := os.Remove("gp.db")
	if err != nil {
		errmsg("os.remove", err)
	}
	err = decompress("gp.zip")
	if err != nil {
		errmsg("decompress", err)
	}
	err = os.Remove("gp.zip")
	if err != nil {
		errmsg("os.Remove", err)
	}

	initDB()
	defer func() {
		err = db.Close()
		if err != nil {
			errmsg("db.Close", err)
		}
	}()

	startAppTime = time.Now()

	if server {
		go startServer()
	}

	if findProxy {
		p := pool.New(numWorkers)
		p.SetHTTPTimeout(timeout)
		links = getAllLinks()
		ips = getAllIP()
		for _, site := range siteList {
			links.set(site)
			p.Add(site, "")
		}
		p.SetTaskTimeout(timeout + 5)
		for result := range p.ResultChan {
			urls := grab(result)
			for _, u := range urls {
				p.Add(u, "")
			}
		}
		err = saveNewIP()
		if err != nil {
			errmsg("saveNewIP", err)
		}
		err = saveLinks()
		if err != nil {
			errmsg("saveLinks", err)
		}
		log.Printf("Add %d ip adress\n", numIPs)
	}

	if checkProxy {
		var (
			totalIP    int64
			totalProxy int64
			anonProxy  int64
		)
		ips = getAllIP()
		p := pool.New(numWorkers)
		p.SetHTTPTimeout(timeout)
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
				p.SetTaskTimeout(timeout + 5)
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
				err = saveAllIP()
				if err != nil {
					errmsg("saveAllIP", err)
				}
			}
		}
		log.Printf("checked %d ip\n", totalIP)
		log.Printf("%d is good\n", totalProxy)
		log.Printf("%d is anon\n", anonProxy)
	}

	err = db.Sync()
	if err != nil {
		errmsg("db.Sync", err)
	}
	err = db.Close()
	if err != nil {
		errmsg("db.Close", err)
	}

	err = compress("gp.db", "gp.zip")
	if err != nil {
		errmsg("compress", err)
	}

	err = os.Remove("gp.db")
	if err != nil {
		errmsg("os.Remove", err)
	}
	err = os.Remove("gp.db.lock")
	if err != nil {
		errmsg("os.Remove", err)
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
