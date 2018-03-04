package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/serbe/adb"

	"github.com/serbe/pool"
)

func findProxy() {
	var (
		addedProxy int64
		newList    []adb.Proxy
	)
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
		newProxy := mp.newProxyInTask(result)
		num := int64(len(newProxy))
		if num > 0 {
			debugmsg("find", num, "proxy in", result.Hostname)
			if link, ok := ml.get(result.Hostname); ok {
				link.Num = link.Num + num
				ml.set(link)
			}
			for _, np := range newProxy {
				newList = append(newList, np)
			}
		}
		if !useNoAddLinks {
			for _, l := range links {
				chkErr("findProxy add to pool", p.Add(l.Hostname, nil))
				debugmsg("add to pool", l.Hostname)
			}
		} else {
			debugmsg("find", len(links), "links in", result.Hostname)
		}
		addedProxy = addedProxy + num
	}
	if !useTestLink {
		debugmsg("save proxy")
		ml.saveAll()
	}
	p.Quit()
	debugmsg(addedProxy, "new proxy found")
	debugmsg("end findProxy")
	checkProxy(newList)
}

func checkProxy(list []adb.Proxy) {
	debugmsg("start checkProxy")
	var (
		checked    int64
		totalProxy int64
		anonProxy  int64
	)

	myIP, err := getMyIP()
	if err != nil {
		errmsg("checkProxy getMyIP", err)
		return
	}
	targetURL := getTarget(myIP)

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
		defer p.Quit()
		p.SetTimeout(timeout)
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
				isNew := false
				if useFUP || useTestScheme {
					isNew = true
				}
				proxy, isOk := mp.taskToProxy(task, isNew, myIP)
				if !isOk {
					continue
				}
				mp.set(proxy)
				if !(useFUP || useTestScheme) {
					saveProxy(proxy)
				}
				if proxy.IsWork {
					if useFUP || useTestScheme {
						saveProxy(proxy)
					}
					totalProxy++
					log.Printf("%d/%d/%d %-15v %-5v %-6v anon=%v\n",
						totalProxy,
						checked,
						listLen,
						task.Proxy.Hostname(),
						task.Proxy.Port(),
						task.Proxy.Scheme,
						proxy.IsAnon,
					)
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
