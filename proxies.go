package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/serbe/adb"
)

func proxyFromString(hostname string) (adb.Proxy, error) {
	var proxy adb.Proxy
	u, err := url.Parse(hostname)
	if err != nil {
		return proxy, fmt.Errorf("failed parse url: %w", err)
	}
	proxy = newProxy(u)
	return proxy, nil
}

func newProxy(u *url.URL) adb.Proxy {
	var proxy adb.Proxy
	port := u.Port()
	proxy.Host = u.Hostname()
	proxy.Port, _ = strconv.Atoi(port)
	proxy.Scheme = u.Scheme
	proxy.CreateAt = time.Now()
	proxy.UpdateAt = time.Now()
	hostname := proxy.Scheme + "://" + proxy.Host + ":" + port
	proxy.Hostname = hostname
	return proxy
}

func taskToProxy(task Task, cfg *config) adb.Proxy {
	var pattern *regexp.Regexp
	proxy, _ := proxyFromString(task.Proxy)

	if cfg.MyIPCheck {
		pattern = regexp.MustCompile(`<td>(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})</td>`)
	} else {
		pattern = regexp.MustCompile(`<p>RemoteAddr: (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):\d{1,5}</p>`)
	}
	proxy.Response = task.Response
	strBody := string(task.Body)
	if pattern.Match(task.Body) && !strings.Contains(strBody, cfg.myIP) {
		proxy.IsWork = true
		if !cfg.MyIPCheck && strings.Count(strBody, "<p>") == 1 {
			proxy.IsAnon = true
		}
		return proxy
	}
	proxy.IsWork = false
	return proxy
}

// func proxyListFromSlice(ips []string, cfg *config) []adb.Proxy {
// 	list, err := cfg.db.CheckNotExists(ips)
// 	chkErr("getNew CheckNotExists", err)
// 	debugmsg(cfg.LogDebug, "find new", len(list))
// 	var pList []adb.Proxy
// 	for i := range list {
// 		proxy, err := proxyFromString(list[i])
// 		if err != nil {
// 			errmsg("getNew proxyFromString", err)
// 		} else {
// 			pList = append(pList, proxy)
// 		}
// 	}
// 	return pList
// }

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

// func (mp *mapProxy) newProxyInTask(task *Task) []adb.Proxy {
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
