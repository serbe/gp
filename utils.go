package main

import (
	"encoding/base64"
	"flag"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/serbe/gocrawl"
)

func checkFlags() {
	flag.Int64Var(&numWorkers, "w", numWorkers, "number of workers")
	flag.Int64Var(&timeout, "t", timeout, "timeout")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.BoolVar(&useFind, "f", useFind, "find new proxy")
	flag.BoolVar(&useCheck, "c", useCheck, "check proxy")
	flag.BoolVar(&useCheckAll, "all", useCheckAll, "check all proxy")
	flag.BoolVar(&useMyIPCheck, "m", useMyIPCheck, "check working proxy on myip.ru")
	flag.BoolVar(&useServer, "s", useServer, "start server")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")
	flag.BoolVar(&useDebug, "d", useDebug, "show debug messages")
	flag.StringVar(&useFile, "pf", useFile, "use file with proxy list")
	flag.StringVar(&testLink, "test", testLink, "link to test it")
	flag.StringVar(&addLink, "a", addLink, "add primary link")
	flag.BoolVar(&useFUP, "fup", useFUP, "test all hosts with 4 frequently used ports")
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

func getLinkList(mL *mapLink, task *gocrawl.Task) []Link {
	var links []Link
	for i := range reURL {
		host, err := getHost(task.Hostname)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if !re.Match(task.Body) {
			continue
		}
		allResults := re.FindAllSubmatch(task.Body, -1)
		for _, result := range allResults {
			hostname := host + "/" + string(result[1])
			if mL.existLink(hostname) {
				continue
			}
			link := mL.newLink(hostname)
			link.Insert = true
			link.UpdateAt = time.Now()
			mL.set(link)
			links = append(links, link)
		}
	}
	return links
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

func getProxyList(body []byte) []Proxy {
	var (
		pList []Proxy
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
			var proxy Proxy
			proxy, err = newProxy(ip, port, false)
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
			var proxy Proxy
			port := convertPort(string(res[2]))
			proxy, err = newProxy(string(res[1]), port, false)
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
			var proxy Proxy
			proxy, err = newProxy(string(res[1]), string(res[2]), false)
			if err == nil {
				pList = append(pList, proxy)
			}
		}
	}
	return pList
}

func grab(mP *mapProxy, mL *mapLink, task *gocrawl.Task) []Link {
	var numProxy int64
	task.Body = cleanBody(task.Body)
	pList := getProxyList(task.Body)
	lList := getLinkList(mL, task)
	for _, p := range pList {
		if mP.existProxy(p.Hostname) {
			continue
		}
		mP.set(p)
		numProxy++
	}
	if numProxy > 0 {
		link := mL.get(task.Hostname)
		link.Num = link.Num + numProxy
		mL.set(link)
		debugmsg("find", numProxy, "in", task.Hostname)
	}
	return lList
}

func chkErr(str string, err error) {
	if err != nil {
		errmsg(str, err)
	}
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}

func debugmsg(str ...interface{}) {
	if useDebug {
		log.Println(str)
	}
}
