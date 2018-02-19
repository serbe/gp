package main

import (
	"io/ioutil"
	"net/url"
	"regexp"
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

// func (mp *mapProxy) remove(hostname string) {
// 	mp.Lock()
// 	delete(mp.values, hostname)
// 	mp.Unlock()
// }

func newProxy(host, port, scheme string) (adb.Proxy, error) {
	var proxy adb.Proxy
	if scheme == "" {
		scheme = "http"
	}
	hostname := scheme + "://" + host + ":" + port
	_, err := url.Parse(hostname)
	if err != nil {
		return proxy, err
	}
	proxy.Hostname = hostname
	proxy.Insert = true
	proxy.Host = host
	proxy.Port = port
	proxy.Scheme = scheme
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

// func proxyIsOld(proxy adb.Proxy) bool {
// 	return time.Since(proxy.UpdateAt) > time.Duration(proxy.Checks)*time.Duration(60*24*7)*time.Minute
// }

func (mp *mapProxy) loadProxyFromFile() {
	if useFile == "" {
		return
	}
	fileBody, err := ioutil.ReadFile(useFile)
	if err != nil {
		errmsg("findProxy ReadFile", err)
		return
	}
	var numProxy int64
	pl := getProxyList(fileBody)
	for _, p := range pl {
		if mp.existProxy(p.Hostname) {
			continue
		}
		mp.set(p)
		numProxy++
	}
}

func getFUPList() []adb.Proxy {
	var list []adb.Proxy
	hosts, err := db.ProxyGetUniqueHosts()
	chkErr("getFUPList ProxyGetUniqueHosts", err)
	ports, err := db.ProxyGetFrequentlyUsedPorts()
	chkErr("getFUPList ProxyGetFrequentlyUsedPorts", err)
	for _, host := range hosts {
		for _, port := range ports {
			proxy, err := newProxy(host, port, "")
			if err == nil {
				list = append(list, proxy)
			}
		}
	}
	return list
}

func getListWithScheme() []adb.Proxy {
	var newList []adb.Proxy
	list, err := db.ProxyGetAllScheme("http")
	chkErr("getListWithScheme ProxyGetAllScheme", err)
	for _, item := range list {
		proxy, err := newProxy(item.Host, item.Port, "https")
		if err == nil {
			newList = append(newList, proxy)
		}
		proxy, err = newProxy(item.Host, item.Port, "socks5")
		if err == nil {
			newList = append(newList, proxy)
		}
	}
	return newList
}

func getProxyListFromDB() []adb.Proxy {
	var (
		list []adb.Proxy
		err  error
	)
	if useTestLink {
		return list
	} else if useCheckAll || useFind {
		list, err = db.ProxyGetAll()
		chkErr("getProxyListFromDB ProxyGetAll", err)
	} else if useFUP {
		list = getFUPList()
	} else if useTestScheme {
		list = getListWithScheme()
	} else {
		list, err = db.ProxyGetAllOld()
		chkErr("getProxyListFromDB ProxyGetAllOld", err)
	}
	return list
}

func saveProxy(p adb.Proxy) {
	debugmsg("save", p)
	if p.Update {
		chkErr("saveProxy ProxyUpdate", db.ProxyUpdate(p))
	} else {
		chkErr("saveProxy ProxyInsert", db.ProxyInsert(p))
	}
}

func (mp *mapProxy) numOfNewProxyInTask(task *pool.Task) int64 {
	var num int64
	body := cleanBody(task.Body)
	proxies := getProxyList(body)
	for _, p := range proxies {
		if mp.existProxy(p.Hostname) {
			continue
		}
		mp.set(p)
		chkErr("numOfNewProxyInTask ProxyInsert", db.ProxyInsert(p))
		num++
	}
	return num
}

func getProxyList(body []byte) []adb.Proxy {
	var (
		pList []adb.Proxy
		err   error
	)
	for i := range baseDecode {
		re := regexp.MustCompile(baseDecode[i])
		if !re.Match(body) {
			continue
		}
		results := re.FindAllSubmatch(body, -1)
		for _, res := range results {
			var ip, port string
			ip, port, err = decodeIP(res[1])
			if err != nil {
				continue
			}
			var proxy adb.Proxy
			proxy, err = newProxy(ip, port, "")
			if err == nil {
				pList = append(pList, proxy)
			}
		}
	}
	for i := range base16 {
		re := regexp.MustCompile(base16[i])
		if !re.Match(body) {
			continue
		}
		results := re.FindAllSubmatch(body, -1)
		for _, res := range results {
			var proxy adb.Proxy
			port := convertPort(string(res[2]))
			proxy, err = newProxy(string(res[1]), port, "")
			if err == nil {
				pList = append(pList, proxy)
			}
		}
	}
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if !re.Match(body) {
			continue
		}
		results := re.FindAllSubmatch(body, -1)
		for _, res := range results {
			var proxy adb.Proxy
			proxy, err = newProxy(string(res[1]), string(res[2]), "")
			if err == nil {
				pList = append(pList, proxy)
			}
		}
	}
	return pList
}
