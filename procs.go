package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/serbe/sites"
)

func findProxy() {
	debugmsg("Start find proxy")
	list := sites.ParseSites(cfg.LogDebug, cfg.LogErrors)
	list = removeDuplicates(list)
	checkProxy(list, false)
	debugmsg("End find proxy")
}

func checkProxy(list []string, isUpdate bool) {
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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	listLen := len(list)
	debugmsg("load proxies", listLen)

	p := NewPool(cfg.Workers)
	debugmsg("start add to pool")
	for i := range list {
		chkErr("add to pool", p.Add(cfg.Target, list[i]))
	}
	debugmsg("end add to pool")
	debugmsg(p.Added(), listLen)
	ch := p.outTasks
	if p.Added() > 0 {
		for p.Added() > p.Completed() {
			select {
			case task := <-ch:
				checked++
				proxy := taskToProxy(task, myIP)

				saveProxy(proxy, isUpdate)
				if proxy.IsWork {
					totalProxy++
					debugmsg(fmt.Sprintf("%d/%d/%d %v %v",
						totalProxy,
						checked,
						listLen,
						proxy.IsAnon,
						proxy.Hostname,
					))
					if proxy.IsAnon {
						anonProxy++
					}
				}
			case <-sig:
				debugmsg("press ctrl+c")
				break
			}
		}
	}
	if listLen > 0 {
		debugmsg(fmt.Sprintf("checked %d from %d", checked, listLen))
	}

	debugmsg(fmt.Sprintf("%d is good", totalProxy))
	debugmsg(fmt.Sprintf("%d is anon", anonProxy))
	debugmsg("end checkProxy")
}
