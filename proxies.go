package main

import (
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

func proxyFromString(hostname string) (adb.Proxy, error) {
	var proxy adb.Proxy
	u, err := url.Parse(hostname)
	if err != nil {
		return proxy, err
	}
	proxy, err = newProxy(u.Hostname(), u.Port(), u.Scheme)
	return proxy, err
}

func newMapProxy() *mapProxy {
	return &mapProxy{values: make(map[string]adb.Proxy)}
}

// func (mp *mapProxy) fillMapProxy(proxyList []adb.Proxy) {
// 	for i := range proxyList {
// 		mp.set(proxyList[i])
// 	}
// }

func (mp *mapProxy) set(proxy adb.Proxy) {
	mp.Lock()
	mp.values[proxy.Hostname] = proxy
	mp.Unlock()
}

// func (mp *mapProxy) setFromString(hostname string) {
// 	proxy, err := proxyFromString(hostname)
// 	if err != nil {
// 		errmsg("setFromString proxyFromString", err)
// 		return
// 	}
// 	mp.Lock()
// 	mp.values[proxy.Hostname] = proxy
// 	mp.Unlock()
// }

func (mp *mapProxy) get(hostname string) (adb.Proxy, bool) {
	mp.Lock()
	proxy, ok := mp.values[hostname]
	mp.Unlock()
	return proxy, ok
}

func newProxy(host, port, scheme string) (adb.Proxy, error) {
	var proxy adb.Proxy
	if scheme == "" {
		scheme = HTTP
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

// func (mp *mapProxy) existProxy(hostname string) bool {
// 	_, ok := mp.get(hostname)
// 	return ok
// }

func (mp *mapProxy) taskToProxy(task *pool.Task, isNew bool, myIP string) (adb.Proxy, bool) {
	proxy, ok := mp.get(task.Proxy.String())
	if !ok {
		return proxy, ok
	}
	pattern := reRemoteIP
	if cfg.MyIPCheck {
		pattern = reMyIP
	}
	if !isNew {
		proxy.Update = true
	}
	proxy.UpdateAt = time.Now()
	proxy.Response = task.ResponseTime
	strBody := string(task.Body)
	if pattern.Match(task.Body) && !strings.Contains(strBody, myIP) {
		proxy.IsWork = true
		proxy.Checks = 0
		if !cfg.MyIPCheck && strings.Count(strBody, "<p>") == 1 {
			proxy.IsAnon = true
		}
		return proxy, ok
	}
	proxy.IsWork = false
	proxy.Checks++
	return proxy, ok
}

func proxyListFromSlice(ips []string) []adb.Proxy {
	list, err := db.CheckNotExists(ips)
	debugmsg("find new", len(list))
	chkErr("getNew CheckNotExists", err)
	var pList []adb.Proxy
	for i := range list {
		proxy, err := proxyFromString(list[i])
		if err != nil {
			errmsg("getNew proxyFromString", err)
		} else {
			pList = append(pList, proxy)
		}
	}
	return pList
}

// func (mp *mapProxy) loadProxyFromFile() {
// 	if testFile == "" {
// 		return
// 	}
// 	fileBody, err := ioutil.ReadFile(testFile)
// 	if err != nil {
// 		errmsg("findProxy ReadFile", err)
// 		return
// 	}
// 	var numProxy int64
// 	pl := getProxyList(fileBody)
// 	for _, p := range pl {
// 		if mp.existProxy(p.Hostname) {
// 			continue
// 		}
// 		mp.set(p)
// 		numProxy++
// 	}
// }

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
	list, err := db.ProxyGetAllScheme(HTTP)
	chkErr("getListWithScheme ProxyGetAllScheme", err)
	for i := range list {
		proxy, err := newProxy(list[i].Host, list[i].Port, HTTPS)
		if err == nil {
			newList = append(newList, proxy)
		}
		proxy, err = newProxy(list[i].Host, list[i].Port, SOCKS5)
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
	} else if useCheckAll {
		list, err = db.ProxyGetAll()
		chkErr("getProxyListFromDB ProxyGetAll", err)
	} else if useFUP {
		list = getFUPList()
	} else if useCheckScheme {
		list = getListWithScheme()
	} else {
		list, err = db.ProxyGetAllOld()
		chkErr("getProxyListFromDB ProxyGetAllOld", err)
	}
	return list
}

func saveProxy(p adb.Proxy) {
	if p.Update {
		chkErr("saveProxy ProxyUpdate", db.ProxyUpdate(&p))
	} else {
		chkErr("saveProxy ProxyInsert", db.ProxyInsert(&p))
	}
}

// func (mp *mapProxy) newProxyInTask(task *pool.Task) []adb.Proxy {
// 	var list []adb.Proxy
// 	body := cleanBody(task.Body)
// 	proxies := getProxyList(body)
// 	for _, p := range proxies {
// 		if mp.existProxy(p.Hostname) {
// 			continue
// 		}
// 		mp.set(p)
// 		// chkErr("numOfNewProxyInTask ProxyInsert", db.ProxyInsert(p))
// 		list = append(list, p)
// 	}
// 	return list
// }

// func getProxyList(body []byte) []adb.Proxy {
// 	var (
// 		pList []adb.Proxy
// 		err   error
// 	)
// 	for i := range baseDecode {
// 		re := regexp.MustCompile(baseDecode[i])
// 		if !re.Match(body) {
// 			continue
// 		}
// 		results := re.FindAllSubmatch(body, -1)
// 		for _, res := range results {
// 			var ip, port string
// 			ip, port, err = decodeIP(res[1])
// 			if err != nil {
// 				continue
// 			}
// 			var proxy adb.Proxy
// 			proxy, err = newProxy(ip, port, "")
// 			if err == nil {
// 				pList = append(pList, proxy)
// 			}
// 		}
// 	}
// 	for i := range base16 {
// 		re := regexp.MustCompile(base16[i])
// 		if !re.Match(body) {
// 			continue
// 		}
// 		results := re.FindAllSubmatch(body, -1)
// 		for _, res := range results {
// 			var proxy adb.Proxy
// 			port := convertPort(string(res[2]))
// 			proxy, err = newProxy(string(res[1]), port, "")
// 			if err == nil {
// 				pList = append(pList, proxy)
// 			}
// 		}
// 	}
// 	for i := range reCommaList {
// 		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
// 		if !re.Match(body) {
// 			continue
// 		}
// 		results := re.FindAllSubmatch(body, -1)
// 		for _, res := range results {
// 			var proxy adb.Proxy
// 			proxy, err = newProxy(string(res[1]), string(res[2]), "")
// 			if err == nil {
// 				pList = append(pList, proxy)
// 			}
// 		}
// 	}
// 	return pList
// }
