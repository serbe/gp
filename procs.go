package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/gocrawl"
)

func findProxy(db *sql.DB) {
	var mL *mapLink
	debugmsg("Start find proxy")
	p := gocrawl.New(numWorkers)
	p.SetTimeout(timeout)
	if testLink != "" {
		mL = newMapLink()
		link := mL.newLink(testLink)
		link.Iterate = true
		mL.set(link)
		log.Println(link)
	} else if addLink != "" {
		mL = newMapLink()
		link := mL.newLink(addLink)
		link.Insert = true
		link.Iterate = true
		mL.set(link)
		log.Println(link)
	} else {
		mL = getAllLinks(db)
	}
	mP := getAllProxy(db)

	if useFile != "" {
		fileBody, err := ioutil.ReadFile(useFile)
		if err != nil {
			errmsg("findProxy ReadFile", err)
		} else {
			var numProxy int64
			pList := getProxyList(fileBody)
			for _, p := range pList {
				if !mP.existProxy(p.Hostname) {
					mP.set(p)
					numProxy++
				}
			}
			log.Println("find", numProxy, "in", useFile)
		}
	}

	debugmsg("start add to pool")
	p.SetTimeout(timeout)
	p.SetQuitTimeout(1000)
	var addedLink int64
	for _, link := range mL.values {
		if link.Iterate && time.Since(link.UpdateAt) > time.Duration(1)*time.Hour {
			err := p.Add(link.Hostname, "")
			if err != nil {
				errmsg("findProxy p.Add", err)
			} else {
				addedLink++
			}
		}
	}
	debugmsg("end add to pool, added", addedLink, "links")
	if addedLink > 0 {
		debugmsg("get from chan")
		for result := range p.ResultChan {
			if result.Error == nil {
				mL.update(result.Hostname)
				links := grab(mP, mL, result)
				for _, l := range links {
					p.Add(l.Hostname, "")
					debugmsg("add to pool", l.Hostname)
				}
			}
		}
		if testLink == "" {
			debugmsg("save proxy")
			saveAllProxy(db, mP)
			saveAllLinks(db, mL)
		}
	}
	debugmsg("end findProxy")
}

func checkProxy(db *sql.DB) {
	debugmsg("start checkProxy")
	var (
		totalIP    int64
		totalProxy int64
		anonProxy  int64
		err        error
		mP         *mapProxy
	)
	if useCheckAll {
		mP = getAllProxy(db)
	} else {
		mP = getOldProxy(db)
	}
	p := gocrawl.New(numWorkers)
	p.SetTimeout(timeout)
	p.SetQuitTimeout(1000)
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
	myIP, err = getExternalIP()
	if err == nil {
		debugmsg("start add to pool")
		for _, proxy := range mP.values {
			if useCheckAll {
				totalIP++
				p.Add(targetURL, proxy.URL.Host)
			} else if proxyIsOld(proxy) {
				totalIP++
				p.Add(targetURL, proxy.URL.Host)
			}
		}
		debugmsg("end add to pool")
		log.Println("Start check", totalIP, "proxyes")
		if totalIP > 0 {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			var checked int
		checkProxyLoop:
			for {
				select {
				case task, ok := <-p.ResultChan:
					if ok {
						checked++
						proxy, isOk := mP.taskToProxy(task)
						if isOk {
							mP.set(proxy)
							if proxy.IsWork {
								log.Printf("%d/%d %-15v %-5v %-12v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
								totalProxy++
								if proxy.IsAnon {
									anonProxy++
								}
							}
						}
					} else {
						debugmsg("break loop by close chan ResultChan")
						break checkProxyLoop
					}
				case <-c:
					debugmsg("breal loop by pressing ctrl+c")
					break checkProxyLoop
				}
			}
			// updateAllProxy(db, mP)
			saveAllProxy(db, mP)
			log.Printf("checked %d ip\n", totalIP)
			log.Printf("%d is good\n", totalProxy)
			log.Printf("%d is anon\n", anonProxy)
			debugmsg("end checkProxy")
		}
	}
}

func checkOnMyIP(db *sql.DB) {
	debugmsg("start checkProxy")
	var (
		totalIP    int64
		totalProxy int64
		anonProxy  int64
		err        error
	)
	mP := getWorkingProxy(db)
	p := gocrawl.New(numWorkers)
	p.SetTimeout(timeout)
	targetURL := "http://myip.ru/"
	myIP, err = getExternalIP()
	if err == nil {
		debugmsg("start add to pool")
		for _, proxy := range mP.values {
			totalIP++
			p.Add(targetURL, proxy.URL.Host)
		}
		debugmsg("end add to pool")
		p.EndWaitingTasks()
		log.Println("Start check on myip", totalIP, "proxyes")
		if totalIP > 0 {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			var checked int
		checkProxyLoop:
			for {
				select {
				case task, ok := <-p.ResultChan:
					if ok {
						checked++
						proxy, isOk := mP.taskMYToProxy(task)
						if isOk {
							mP.set(proxy)
							if proxy.IsWork {
								log.Printf("%d/%d %-15v %-5v %-12v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
								totalProxy++
							}
						}
					} else {
						debugmsg("break loop by close chan ResultChan")
						break checkProxyLoop
					}
				case <-c:
					debugmsg("breal loop by pressing ctrl+c")
					break checkProxyLoop
				}
			}
			saveAllProxy(db, mP)
			log.Printf("checked %d ip\n", totalIP)
			log.Printf("%d is good\n", totalProxy)
			log.Printf("%d is anon\n", anonProxy)
			debugmsg("end checkProxy")
		}
	}
}
