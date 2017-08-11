package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/pool"
)

func findProxy() {
	debugmsg("Start find proxy")
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	mL = getAllLinks()
	mP = getAllProxy()

	links := mL.oldLinks()

	if len(links) > 0 {
		debugmsg("start add to pool")
		p.SetTaskTimeout(5)
		for i, link := range links {
			err := p.Add(link.Hostname, new(url.URL))
			if err != nil {
				errmsg("findProxy p.Add", err)
			}
			debugmsg("add to pool", i, link.Hostname)
		}
		debugmsg("end add to pool")
		debugmsg("get from chan")
		for result := range p.ResultChan {
			if result.Error == nil {
				mL.update(result.Hostname)
				urls := grab(result)
				for _, u := range urls {
					p.Add(u, new(url.URL))
					debugmsg("add to pool", u)
				}
			}
		}
		debugmsg("save proxy")
		saveAllProxy(mP)
		saveAllLinks(mL)
	}
	log.Printf("Add %d ip adress\n", numIPs)
}

func checkProxy() {
	var (
		totalIP    int64
		totalProxy int64
		anonProxy  int64
		err        error
	)
	mP = getAllProxy()
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
	myIP, err = getExternalIP()
	if err == nil {
		week := time.Duration(60*24*7) * time.Minute
		startTime := time.Now()
		for _, proxy := range mP.values {
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
				case task, ok := <-p.ResultChan:
					checked++
					if ok {
						proxy := taskToProxy(task)
						proxy.Response = task.ResponceTime
						mP.set(proxy)
						if proxy.IsWork {
							log.Printf("%d/%d %-15v %-5v %-10v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
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
			log.Printf("checked %d ip\n", totalIP)
			log.Printf("%d is good\n", totalProxy)
			log.Printf("%d is anon\n", anonProxy)
		}
	}
}
