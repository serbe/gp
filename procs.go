package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/serbe/pool"
	"github.com/serbe/sites"
)

func findProxy() {
	debugmsg("Start find proxy")
	list := sites.ParseSites(cfg.LogDebug, cfg.LogErrors)
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listLen := len(list)
	debugmsg("load proxies", listLen)

	p := pool.New(cfg.Workers)
	// defer p.Quit()
	p.SetTimeout(cfg.Timeout)
	debugmsg("start add to pool")
	for i := range list {
		chkErr("add to pool", p.Add(cfg.Target, list[i]))
	}
	debugmsg("end add to pool")
	debugmsg(p.GetAddedTasks(), listLen)
	p.EndWaitingTasks()
	p.SetQuitTimeout(cfg.Timeout + 1000)

	if p.GetAddedTasks() > 0 {
	checkProxyLoop:
		for {
			select {
			case task, ok := <-p.ResultChan:
				if !ok {
					debugmsg(p.GetAddedTasks(), p.GetCompletedTasks(), listLen)
					debugmsg("break loop by close chan ResultChan")
					break checkProxyLoop
				}
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
			case <-c:
				debugmsg("break loop by pressing ctrl+c")
				break checkProxyLoop
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
