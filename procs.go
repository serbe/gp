package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/pool"
)

func findProxy() {
	var addedProxy int64
	debugmsg("Start find proxy")
	p := pool.New(numWorkers)
	p.SetTimeout(timeout)
	ml := getMapLink()
	mp := newMapProxy()
	list := getProxyListFromDB()
	mp.fillMapProxy(list)
	mp.loadProxyFromFile()

	debugmsg("load links", len(ml.values))
	debugmsg("load proxies", len(mp.values))
	debugmsg("start add to pool")
	p.SetTimeout(timeout)
	p.SetQuitTimeout(2000)
	for _, link := range ml.values {
		if useAddLink || useTestLink || link.Iterate && time.Since(link.UpdateAt) > time.Duration(1)*time.Hour {
			err := p.Add(link.Hostname, nil)
			if err == nil {
				addedProxy++
			}
			chkErr("findProxy p.Add", err)
		}
	}
	if addedProxy == 0 {
		debugmsg("not added tasks to pool")
		return
	}
	debugmsg("end add to pool, added", p.GetAddedTasks(), "links")
	debugmsg("start get from chan")
	for result := range p.ResultChan {
		if result.Error != nil {
			errmsg("result", result.Error)
			continue
		}
		ml.update(result.Hostname)
		links := ml.getNewLinksFromTask(result)
		num := mp.numOfNewProxyInTask(result)
		if num > 0 {
			debugmsg("find", num, "proxy in", result.Hostname)
			if link, ok := ml.get(result.Hostname); ok {
				link.Num = link.Num + num
				ml.set(link)
				debugmsg("find", num, "in", result.Hostname)
			}
		}
		for _, l := range links {
			chkErr("findProxy add to pool", p.Add(l.Hostname, nil))
			debugmsg("add to pool", l.Hostname)
		}
		addedProxy = addedProxy + num
	}
	if !useTestLink {
		debugmsg("save proxy")
		ml.saveAll()
	}
	debugmsg(addedProxy, "new proxy found")
	debugmsg("end findProxy")
}

func checkProxy() {
	debugmsg("start checkProxy")
	var (
		checked    int64
		totalProxy int64
		anonProxy  int64
		err        error
	)
	list := getProxyListFromDB()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listLen := len(list)
	debugmsg("load proxies", listLen)

breakCheckProxyLoop:
	for j := 0; j < listLen; {
		mp := newMapProxy()
		var r = 10000
		if j+10000 > listLen {
			r = listLen % 10000
		}
		for i := 0; i < r; i++ {
			mp.set(list[j])
			j++
		}
		p := pool.New(numWorkers)
		p.SetTimeout(timeout)
		targetURL := getTarget()
		myIP, err = getExternalIP()
		if err != nil {
			errmsg("checkProxy getExternalIP", err)
			return
		}
		debugmsg("start add to pool")
		for _, proxy := range mp.values {
			proxyURL, err := url.Parse(proxy.Hostname)
			chkErr("parse url", err)
			chkErr("add to pool", p.Add(targetURL, proxyURL))
		}
		debugmsg("end add to pool")
		p.EndWaitingTasks()
		if p.GetAddedTasks() == 0 {
			debugmsg("no task added to pool")
			return
		}
	checkProxyLoop:
		for {
			select {
			case task, ok := <-p.ResultChan:
				if !ok {
					debugmsg("break loop by close chan ResultChan")
					break checkProxyLoop
				}
				checked++
				proxy, isOk := mp.taskToProxy(task)
				if !isOk {
					continue
				}
				mp.set(proxy)
				if !useFUP {
					saveProxy(proxy)
				}
				if proxy.IsWork {
					if useFUP {
						saveProxy(proxy)
					}
					log.Printf("%d/%d %-15v %-5v %-12v anon=%v\n", checked, listLen, task.Proxy.Hostname(), task.Proxy.Port(), task.ResponceTime, proxy.IsAnon)
					totalProxy++
					if proxy.IsAnon {
						anonProxy++
					}
				}
			case <-c:
				debugmsg("break loop by pressing ctrl+c")
				break breakCheckProxyLoop
			}
		}
	}
	log.Printf("%d is good\n", totalProxy)
	log.Printf("%d is anon\n", anonProxy)
	debugmsg("end checkProxy")
}
