package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/serbe/adb"
	"github.com/serbe/pool"
)

type mapProxy struct {
	sync.RWMutex
	values map[string]adb.Proxy
}

func (mp *mapProxy) fillMapProxy(proxyList []adb.Proxy) {
	for _, proxy := range proxyList {
		mp.set(proxy)
	}
}

func newMapProxy() *mapProxy {
	return &mapProxy{values: make(map[string]adb.Proxy)}
}

func (mp *mapProxy) set(proxy adb.Proxy) {
	mp.Lock()
	mp.values[proxy.Hostname] = proxy
	mp.Unlock()
}

func (mp *mapProxy) get(hostname string) (adb.Proxy, bool) {
	mp.Lock()
	proxy, ok := mp.values[hostname]
	mp.Unlock()
	return proxy, ok
}

func (mp *mapProxy) remove(hostname string) {
	mp.Lock()
	delete(mp.values, hostname)
	mp.Unlock()
}

func newProxy(host, port string, ssl bool) (adb.Proxy, error) {
	var (
		proxy  adb.Proxy
		schema string
	)
	if ssl {
		schema = "https://"
	} else {
		schema = "http://"
	}
	hostname := schema + host + ":" + port
	_, err := url.Parse(hostname)
	if err != nil {
		return proxy, err
	}
	proxy.Hostname = hostname
	proxy.Insert = true
	proxy.Host = host
	proxy.Port = port
	proxy.CreateAt = time.Now()
	return proxy, err
}

func (mp *mapProxy) existProxy(hostname string) bool {
	_, ok := mp.get(hostname)
	return ok
}

func (mp *mapProxy) taskToProxy(task *pool.Task) (adb.Proxy, bool) {
	proxy, ok := mp.get(task.Proxy.String())
	if !ok {
		return proxy, ok
	}
	pattern := reRemoteIP
	if useMyIPCheck {
		pattern = reMyIP
	}
	proxy.Update = true
	proxy.UpdateAt = time.Now()
	proxy.Response = task.ResponceTime
	strBody := string(task.Body)
	if pattern.Match(task.Body) && !strings.Contains(strBody, myIP) {
		proxy.IsWork = true
		proxy.Checks = 0
		if !useMyIPCheck && strings.Count(strBody, "<p>") == 1 {
			proxy.IsAnon = true
		}
		return proxy, ok
	}
	proxy.IsWork = false
	proxy.Checks++
	return proxy, ok
}

func proxyIsOld(proxy adb.Proxy) bool {
	return time.Since(proxy.UpdateAt) > time.Duration(proxy.Checks)*time.Duration(60*24*7)*time.Minute
}

func loadProxyFromFile(mp *mapProxy) {
	if useFile == "" {
		return
	}
	fileBody, err := ioutil.ReadFile(useFile)
	if err != nil {
		errmsg("findProxy ReadFile", err)
		return
	}
	var numProxy int64
	pList := getProxyList(fileBody)
	for _, p := range pList {
		if mp.existProxy(p.Hostname) {
			continue
		}
		mp.set(p)
		numProxy++
	}
	log.Println("find", numProxy, "in", useFile)
}

// func getFUPList() *mapProxy {
// 	mp := getAllProxy()
// 	hosts := uniqueHosts()
// 	ports := frequentlyUsedPorts()
// 	for _, host := range hosts {
// 		for _, port := range ports {
// 			proxy, err := newProxy(host, port, false)
// 			if err == nil {
// 				if !mp.existProxy(proxy.Hostname) {
// 					mp.set(proxy)
// 				}
// 			}
// 		}
// 	}
// 	return mp
// }

func getMapProxy() *mapProxy {
	mp := newMapProxy()
	if useCheckAll {
		mp.fillMapProxy(db.ProxyGetAll())
		// } else if useFUP {
		// 	mp = getFUPList()
	} else {
		mp.fillMapProxy(db.ProxyGetAllOld())
	}
	return mp
}

func saveProxy(p adb.Proxy) error {
	if p.Update {
		return db.ProxyUpdate(p)
	}
	return db.ProxyCreate(p)
}
