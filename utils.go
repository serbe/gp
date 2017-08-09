package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/serbe/pool"
)

func checkFlags() {
	flag.IntVar(&numWorkers, "w", numWorkers, "number of workers")
	flag.IntVar(&timeout, "t", timeout, "timeout")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.BoolVar(&useFind, "f", useFind, "find new proxy")
	flag.BoolVar(&useCheck, "c", useCheck, "check proxy")
	flag.BoolVar(&useServer, "s", useServer, "start server")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")
	flag.BoolVar(&createTables, "m", createTables, "create tables in new database")
	flag.Parse()
}

func cleanBody(body []byte) []byte {
	for i := range replace {
		re := regexp.MustCompile(replace[i][0])
		if re.Match(body) {
			body = re.ReplaceAll(body, []byte(replace[i][1]))
		}
	}
	return body
}

func getListURL(task pool.Task) []string {
	var urls []string
	for i := range reURL {
		host, err := getHost(task.Target.String())
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if re.Match(task.Body) {
			allResults := re.FindAllSubmatch(task.Body, -1)
			for _, result := range allResults {
				fullURL := host + "/" + string(result[1])
				if isOld(links.get(fullURL)) {
					links.set(fullURL)
					urls = append(urls, fullURL)
				}
			}
		}
	}
	return urls
}

func setProxy(host string, port string, base int, ssl bool) {
	portInt, err := strconv.ParseInt(port, base, 32)
	if err != nil {
		return
	}
	var portStr string
	if base == 10 {
		portStr = port
	} else {
		portStr = strconv.Itoa(int(portInt))
	}
	proxy, err := newProxy(host, portStr, false)
	if err == nil {
		numIPs++
		ips.set(proxy)
	}
}

func newProxy(host, port string, ssl bool) (Proxy, error) {
	var (
		proxy  Proxy
		schema string
	)
	if ssl {
		schema = "https://"
	} else {
		schema = "http://"
	}
	URL, err := url.Parse(schema + host + ":" + port)
	proxy.URL = URL
	proxy.CreateAt = time.Now()
	return proxy, err
}

func decodeIP(src []byte) (string, string, error) {
	out, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return "", "", err
	}
	split := strings.Split(string(out), ":")
	if len(split) == 2 {
		return split[0], split[1], nil
	}
	return "", "", err
}

func getListIP(body []byte) {
	for i := range baseDecode {
		re := regexp.MustCompile(baseDecode[i])
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				ip, port, err := decodeIP(res[1])
				if err == nil {
					setProxy(ip, port, 10, false)
				}
			}
		}
	}
	for i := range base16 {
		re := regexp.MustCompile(base16[i])
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				setProxy(string(res[1]), string(res[2]), 16, false)
			}
		}
	}
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				setProxy(string(res[1]), string(res[2]), 10, false)
			}
		}
	}
}

func ipFromProxy(proxy Proxy) (IP, error) {
	var (
		ip  IP
		err error
	)
	ip.Hostname = proxy.URL.Hostname()
	ip.Checks = proxy.Checks
	ip.IsAnon = proxy.IsAnon
	ip.IsWork = proxy.IsWork
	ip.Response = proxy.Response
	ip.CreateAt = proxy.CreateAt
	ip.UpdateAt = proxy.UpdateAt
	return ip, err
}

func proxyFromIP(ip IP) (Proxy, error) {
	var (
		proxy Proxy
		err   error
	)
	proxy.URL, err = url.Parse(ip.Hostname)
	proxy.Checks = ip.Checks
	proxy.IsAnon = ip.IsAnon
	proxy.IsWork = ip.IsWork
	proxy.Response = ip.Response
	proxy.CreateAt = ip.CreateAt
	proxy.UpdateAt = ip.UpdateAt
	return proxy, err
}

func isOld(link Link) bool {
	currentTime := time.Now()
	return currentTime.Sub(link.CheckAt) > time.Duration(720)*time.Minute
}

func grab(task pool.Task) []string {
	task.Body = cleanBody(task.Body)
	oldNumIP := numIPs
	getListIP(task.Body)
	if numIPs-oldNumIP > 0 {
		log.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, task.Target.String())
	}
	urls := getListURL(task)
	return urls
}

func check(task pool.Task) Proxy {
	proxy := Proxy{
		URL:      task.Proxy,
		UpdateAt: time.Now(),
	}
	if task.Error == nil {
		strBody := string(task.Body)
		if reRemoteIP.Match(task.Body) && !strings.Contains(strBody, myIP) {
			proxy.IsWork = true
			proxy.Checks = 0
			if strings.Count(strBody, "<p>") == 1 {
				proxy.IsAnon = true
			}
			return proxy
		}
	}
	proxy.Checks++
	return proxy
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}

func findProxy() {
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	links = getAllLinks()
	ips = getAllProxy()
	for _, link := range links.values {
		if time.Since(link.CheckAt) > time.Duration(5)*time.Minute {
			p.Add(link.Host, new(url.URL))
		}
	}
	p.SetTaskTimeout(3)
	for result := range p.ResultChan {
		urls := grab(result)
		for _, u := range urls {
			p.Add(u, new(url.URL))
		}
	}
	saveAllProxy(ips)
	saveAllLinks(links)
	log.Printf("Add %d ip adress\n", numIPs)
}

func checkProxy() {
	var (
		totalIP    int64
		totalProxy int64
		anonProxy  int64
		err        error
	)
	ips = getAllProxy()
	p := pool.New(numWorkers)
	p.SetHTTPTimeout(timeout)
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", serverPort)
	myIP, err = getExternalIP()
	if err == nil {
		week := time.Duration(60*24*7) * time.Minute
		startTime := time.Now()
		for _, proxy := range ips.values {
			if (proxy.UpdateAt == time.Time{} || proxy.UpdateAt != time.Time{} && startTime.Sub(proxy.UpdateAt) > time.Duration(proxy.Checks)*week) {
				totalIP++
				p.Add(targetURL, proxy.URL)
			}
		}
		log.Println("Start check", totalIP, "proxyes")
		if totalIP > 0 {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			p.SetTaskTimeout(2)
			var checked int
		checkProxyLoop:
			for {
				select {
				case result, ok := <-p.ResultChan:
					checked++
					if ok {
						proxy := check(result)
						proxy.Response = result.ResponceTime
						ips.set(proxy)
						if proxy.IsWork {
							log.Printf("%d/%d %-15v %-5v %-10v anon=%v\n", checked, totalIP, result.Proxy.Hostname(), result.Proxy.Port(), result.ResponceTime, proxy.IsAnon)
							totalProxy++
							if proxy.IsAnon {
								anonProxy++
							}
						}
					} else {
						break checkProxyLoop
					}
				case <-c:
					break checkProxyLoop
				}
			}
			log.Printf("checked %d ip\n", totalIP)
			log.Printf("%d is good\n", totalProxy)
			log.Printf("%d is anon\n", anonProxy)
		}
	}
}
