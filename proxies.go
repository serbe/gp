package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/serbe/adb"
	"github.com/serbe/pool"
)

func proxyFromString(hostname string) (adb.Proxy, error) {
	var proxy adb.Proxy
	u, err := url.Parse(hostname)
	if err != nil {
		return proxy, err
	}
	proxy = newProxy(u)
	return proxy, err
}

func newProxy(u *url.URL) adb.Proxy {
	var proxy adb.Proxy
	port := u.Port()
	proxy.Host = u.Hostname()
	proxy.Port, _ = strconv.Atoi(port)
	proxy.Scheme = u.Scheme
	proxy.CreateAt = time.Now()
	hostname := proxy.Scheme + "://" + proxy.Host + ":" + port
	proxy.Hostname = hostname
	return proxy
}

func taskToProxy(task *pool.TaskResult, myIP string) adb.Proxy {
	var proxy adb.Proxy
	proxy.Hostname = task.Proxy
	pattern := reRemoteIP
	if cfg.MyIPCheck {
		pattern = reMyIP
	}
	proxy.CreateAt = time.Now()
	proxy.UpdateAt = time.Now()
	proxy.Response = task.ResponseTime
	strBody := string(task.Body)
	if pattern.Match(task.Body) && !strings.Contains(strBody, myIP) {
		proxy.IsWork = true
		if !cfg.MyIPCheck && strings.Count(strBody, "<p>") == 1 {
			proxy.IsAnon = true
		}
		return proxy
	}
	proxy.IsWork = false
	return proxy
}

func proxyListFromSlice(ips []string) []adb.Proxy {
	list, err := db.CheckNotExists(ips)
	chkErr("getNew CheckNotExists", err)
	debugmsg("find new", len(list))
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

func getFUPList() []string {
	var list []string
	hosts, err := db.GetUniqueHosts()
	chkErr("getFUPList ProxyGetUniqueHosts", err)
	ports, err := db.GetFrequentlyUsedPorts()
	chkErr("getFUPList ProxyGetFrequentlyUsedPorts", err)
	for _, host := range hosts {
		for _, port := range ports {
			u := "http://" + host + ":" + strconv.Itoa(port)
			_, err := url.Parse(u)
			if err == nil {
				list = append(list, u)
			}
		}
	}
	return list
}

func getListWithScheme() []string {
	var newList []string
	list, err := db.GetAllScheme(HTTP)
	chkErr("getListWithScheme ProxyGetAllScheme", err)
	for i := range list {
		u, err := url.Parse(list[i])
		if err == nil {
			newList = append(newList, "https://"+u.Host+":"+u.Port())
			newList = append(newList, "socks5://"+u.Host+":"+u.Port())
		}
	}
	return newList
}

func getProxyListFromDB() []string {
	var (
		list []string
		err  error
	)
	if useTestLink {
		return list
	} else if useCheckAll {
		list, err = db.GetAll()
		chkErr("getProxyListFromDB ProxyGetAll", err)
	} else if useFUP {
		list = getFUPList()
	} else if useCheckScheme {
		list = getListWithScheme()
	} else {
		list, err = db.GetAllOld()
		chkErr("getProxyListFromDB ProxyGetAllOld", err)
	}
	return list
}

func saveProxy(p adb.Proxy, isUpdate bool) {
	if isUpdate {
		chkErr("saveProxy Update "+p.Hostname, db.Update(&p))
	} else {
		chkErr("saveProxy Insert "+p.Hostname, db.Insert(&p))
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
