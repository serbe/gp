package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/pool"
)

func main() {
	var (
		findProxy  = false
		checkProxy = false
		server     = false
	)

	flag.IntVar(&numWorkers, "w", numWorkers, "number of workers")
	flag.IntVar(&timeout, "t", timeout, "timeout")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.BoolVar(&findProxy, "f", findProxy, "find new proxy")
	flag.BoolVar(&checkProxy, "c", checkProxy, "check proxy")
	flag.BoolVar(&server, "s", server, "start server")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")
	flag.Parse()

	initDB()
	defer func() {
		err := db.Close()
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
		ips = getAllProxy()
		for _, site := range siteList {
			links.set(site)
			p.Add(site, new(url.URL))
		}
		p.SetTaskTimeout(3)
		for result := range p.ResultChan {
			urls := grab(result)
			for _, u := range urls {
				p.Add(u, new(url.URL))
			}
		}
		saveAllProxy(ips)
		log.Printf("Add %d ip adress\n", numIPs)
	}

	if checkProxy {
		var (
			totalIP    int64
			totalProxy int64
			anonProxy  int64
			err        error
		)
		ips = getAllProxy()
		p := pool.New(numWorkers)
		p.SetHTTPTimeout(timeout)
		targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
		myIP, err = getExternalIP()
		if err == nil {
			week := time.Duration(60*24*7) * time.Minute
			startTime := time.Now()
			for _, proxy := range ips.values {
				if (proxy.UpdateAt == time.Time{} || proxy.UpdateAt != time.Time{} && startTime.Sub(proxy.UpdateAt) > time.Duration(proxy.Checks)*week) {
					totalIP++
					p.Add(targetURL, proxy.URL)
				}
			}
			log.Println("Start check", totalIP, "proxyes")
			if totalIP > 0 {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)
				p.SetTaskTimeout(2)
				var checked int
			checkProxyLoop:
				for {
					select {
					case result, ok := <-p.ResultChan:
						checked++
						if ok {
							proxy := check(result)
							proxy.Response = result.ResponceTime
							ips.set(proxy)
							if proxy.IsWork {
								log.Printf("%d/%d %-15v %-5v %-10v anon=%v\n", checked, totalIP, result.Proxy.Hostname(), result.Proxy.Port(), result.ResponceTime, proxy.IsAnon)
								totalProxy++
								if proxy.IsAnon {
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
			}
		}
		log.Printf("checked %d ip\n", totalIP)
		log.Printf("%d is good\n", totalProxy)
		log.Printf("%d is anon\n", anonProxy)
	}
	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
