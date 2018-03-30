package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/serbe/adb"
	"github.com/serbe/pool"
)

func findProxy() {
	debugmsg("Start find proxy")
	ips := parseSites()

	list := proxyListFromSlice(ips)

	checkProxy(list)
	debugmsg("End find proxy")
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	listLen := len(list)
	debugmsg("load proxies", listLen)

breakCheckProxyLoop:
	for j := 0; j < listLen; {
		mp := newMapProxy()
		var r = 2000
		if j+2000 > listLen {
			r = listLen % 2000
		}
		for i := 0; i < r; i++ {
			mp.set(list[j])
			j++
		}
		p := pool.New(cfg.Workers)
		// defer p.Quit()
		p.SetTimeout(cfg.Timeout)
		debugmsg("start add to pool")
		for _, proxy := range mp.values {
			proxyURL, err := url.Parse(proxy.Hostname)
			chkErr("parse url", err)
			chkErr("add to pool", p.Add(cfg.Target, proxyURL))
		}
		debugmsg("end add to pool")
		debugmsg(j, p.GetAddedTasks(), listLen)
		p.EndWaitingTasks()
		p.SetQuitTimeout(cfg.Timeout + 1000)
		if p.GetAddedTasks() > 0 {
		checkProxyLoop:
			for {
				select {
				case task, ok := <-p.ResultChan:
					if !ok {
						debugmsg(j, p.GetAddedTasks(), p.GetCompletedTasks(), listLen)
						debugmsg("break loop by close chan ResultChan")
						break checkProxyLoop
					}
					checked++
					isNew := false
					if useFUP || useCheckScheme || useFind {
						isNew = true
					}
					proxy, isOk := mp.taskToProxy(task, isNew, myIP)
					if !isOk {
						continue
					}
					mp.set(proxy)
					if !(useFUP || useCheckScheme) {
						saveProxy(proxy)
					}
					if proxy.IsWork {
						if useFUP || useCheckScheme {
							saveProxy(proxy)
						}
						totalProxy++
						debugmsg(fmt.Sprintf("%d/%d/%d %-15v %-5v %-6v %v",
							totalProxy,
							checked,
							listLen,
							task.Proxy.Hostname(),
							task.Proxy.Port(),
							task.Proxy.Scheme,
							proxy.IsAnon,
						))
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
		if listLen > 0 {
			debugmsg(fmt.Sprintf("checked %d from %d", checked, listLen))
		}
	}
	debugmsg(fmt.Sprintf("%d is good", totalProxy))
	debugmsg(fmt.Sprintf("%d is anon", anonProxy))
	debugmsg("end checkProxy")
}

func parseSites() []string {
	type parser struct {
		name string
		ips  []string
	}
	var ips []string
	debugmsg("start parse sites")
	ch := make(chan parser)
	go func() {
		data := parser{name: "freeproxylist", ips: freeproxylist()}
		ch <- data
	}()
	go func() {
		data := parser{name: "freeproxylistcom", ips: freeproxylistcom()}
		ch <- data
	}()
	go func() {
		data := parser{name: "gatherproxycom", ips: gatherproxycom()}
		ch <- data
	}()
	go func() {
		data := parser{name: "kuaidaili", ips: kuaidaili()}
		ch <- data
	}()
	go func() {
		data := parser{name: "proxylistorg", ips: proxylistorg()}
		ch <- data
	}()
	go func() {
		data := parser{name: "proxyserverlist24top", ips: proxyserverlist24top()}
		ch <- data
	}()
	go func() {
		data := parser{name: "rawlist", ips: rawlist()}
		ch <- data
	}()
	go func() {
		data := parser{name: "webanetlabs", ips: webanetlabs()}
		ch <- data
	}()
	for i := 0; i < 8; i++ {
		data := <-ch
		ips = append(ips, data.ips...)
	}
	debugmsg("end parse sites, found", len(ips), "proxy")
	return ips
}
