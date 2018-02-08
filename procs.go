package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/gocrawl"
)

func findProxy(db *sql.DB) {
	debugmsg("Start find proxy")
	p := gocrawl.New(numWorkers)
	p.SetTimeout(timeout)
	mL := getMapLink(db)
	mP := getAllProxy(db)

	loadProxyFromFile(mP)

	debugmsg("start add to pool")
	p.SetTimeout(timeout)
	p.SetQuitTimeout(2000)
	for _, link := range mL.values {
		if link.Iterate && time.Since(link.UpdateAt) > time.Duration(1)*time.Hour {
			err := p.Add(link.Hostname, nil)
			chkErr("findProxy p.Add", err)
		}
	}
	if p.GetAddedTasks() == 0 {
		debugmsg("not added tasks to pool")
		return
	}
	debugmsg("end add to pool, added", p.GetAddedTasks(), "links")
	debugmsg("start get from chan")
	for result := range p.ResultChan {
		if result.Error != nil {
			continue
		}
		mL.update(result.Hostname)
		links := grab(mP, mL, result)
		for _, l := range links {
			p.Add(l.Hostname, nil)
			debugmsg("add to pool", l.Hostname)
		}
	}
	if testLink == "" {
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
	mP := getMapProxy(db)
	p := gocrawl.New(numWorkers)
	p.SetTimeout(timeout)
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
	myIP, err = getExternalIP()
	if err != nil {
		errmsg("checkProxy getExternalIP", err)
		return
	}
	debugmsg("start add to pool")
	for _, proxy := range mP.values {
		if useCheckAll || proxyIsOld(proxy) {
			totalIP++
			chkErr("add to pool", p.Add(targetURL, proxy.URL))
		}
	}
	debugmsg("end add to pool")
	p.EndWaitingTasks()
	log.Println("Start check", totalIP, "proxies")
	if p.GetAddedTasks() == 0 {
		debugmsg("no task added to pool")
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var checked int
checkProxyLoop:
	for {
		select {
		case task, ok := <-p.ResultChan:
			if !ok {
				debugmsg("break loop by close chan ResultChan")
				break checkProxyLoop
			}
			checked++
			proxy, isOk := mP.taskToProxy(task)
			if !isOk {
				continue
			}
			mP.set(proxy)
			if proxy.IsWork {
				log.Printf("%d/%d %-15v %-5v %-12v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
				totalProxy++
				if proxy.IsAnon {
					anonProxy++
				}
			}
		case <-c:
			debugmsg("break loop by pressing ctrl+c")
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
	if err != nil {
		errmsg("getExternalIP", err)
		return
	}
	debugmsg("start add to pool")
	for _, proxy := range mP.values {
		totalIP++
		p.Add(targetURL, proxy.URL)
	}
	debugmsg("end add to pool")
	p.EndWaitingTasks()
	log.Println("Start check on myIP", totalIP, "proxies")
	if totalIP > 0 {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		var checked int
	checkProxyLoop:
		for {
			select {
			case task, ok := <-p.ResultChan:
				if !ok {
					debugmsg("break loop by close chan ResultChan")
					break checkProxyLoop
				}
				checked++
				proxy, isOk := mP.taskMYToProxy(task)
				if !isOk {
					continue
				}
				mP.set(proxy)
				if proxy.IsWork {
					log.Printf("%d/%d %-15v %-5v %-12v anon=%v\n", checked, totalIP, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
					totalProxy++
				}
			case <-c:
				debugmsg("break loop by pressing ctrl+c")
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
