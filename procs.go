package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/pool"
)

func findProxy(db *sql.DB) {
	debugmsg("Start find proxy")
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	mL := getAllLinks(db)
	mP := getAllProxy(db)

	// var (
	// 	fileBody []byte
	// 	err      error
	// )

	// if useFile != "" {
	// 	fileBody, err = ioutil.ReadFile(useFile)
	// }

	debugmsg("start add to pool")
	p.SetTaskTimeout(5)
	var addedLink int64
	for _, link := range mL.values {
		if link.Iterate && time.Since(link.UpdateAt) > time.Duration(1)*time.Hour {
			err := p.Add(link.Hostname, new(url.URL))
			if err != nil {
				errmsg("findProxy p.Add", err)
			} else {
				addedLink++
			}
		}
	}
	debugmsg("end add to pool")
	if addedLink > 0 {
		debugmsg("get from chan")
		for result := range p.ResultChan {
			if result.Error == nil {
				mL.update(result.Hostname)
				links := grab(mP, mL, result)
				for _, l := range links {
					p.Add(l.Hostname, new(url.URL))
					debugmsg("add to pool", l.Hostname)
				}
			}
		}
		debugmsg("save proxy")
		saveAllProxy(db, mP)
		saveAllLinks(db, mL)
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
	)
	mP := getOldProxy(db)
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	p.SetTaskTimeout(2)
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
	myIP, err = getExternalIP()
	if err == nil {
		debugmsg("start add to pool")
		for _, proxy := range mP.values {
			if proxyIsOld(proxy) {
				totalIP++
				p.Add(targetURL, proxy.URL)
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
								log.Printf("%d/%d %-15v %-5v %-10v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
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
			updateAllProxy(db, mP)
			log.Printf("checked %d ip\n", totalIP)
			log.Printf("%d is good\n", totalProxy)
			log.Printf("%d is anon\n", anonProxy)
			debugmsg("end checkProxy")
		}
	}
}
