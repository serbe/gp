package main

import (
	"encoding/base64"
	"flag"
	"log"
	"regexp"
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
	flag.BoolVar(&useDebug, "d", useDebug, "show debug messages")
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
		host, err := getHost(task.Hostname)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if re.Match(task.Body) {
			allResults := re.FindAllSubmatch(task.Body, -1)
			for _, result := range allResults {
				hostname := host + "/" + string(result[1])
				if mL.existLink(hostname) {
					if mL.isOldLink(hostname) {
						mL.update(hostname)
						urls = append(urls, hostname)
					}
				} else {
					link := mL.newLink(hostname)
					link.Insert = true
					link.UpdateAt = time.Now()
					mL.set(link)
					urls = append(urls, hostname)
				}
			}
		}
	}
	return urls
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
					setProxy(ip, port, false)
				}
			}
		}
	}
	for i := range base16 {
		re := regexp.MustCompile(base16[i])
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				port := convPort(string(res[2]), 16)
				setProxy(string(res[1]), port, false)
			}
		}
	}
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				setProxy(string(res[1]), string(res[2]), false)
			}
		}
	}
}

func grab(task pool.Task) []string {
	task.Body = cleanBody(task.Body)
	oldNumIP := numIPs
	getListIP(task.Body)
	if numIPs-oldNumIP > 0 {
		log.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, task.Hostname)
	}
	urls := getListURL(task)
	return urls
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
