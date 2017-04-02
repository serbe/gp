package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/serbe/pool"
)

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

func getListIP(body []byte) {
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				ip := string(res[1])
				port := string(res[2])
				portInt, _ := strconv.Atoi(port)
				if ip != "0.0.0.0" && portInt < 65535 {
					ipWithPort := ip + ":" + port
					if ips.get(ipWithPort).Addr == "" {
						numIPs++
						ips.set(ipWithPort, newIP(ip, port, false))
					}
				}
			}
		}
	}
	return
}

func newIP(addr, port string, ssl bool) ipType {
	var ip ipType
	ip.Addr = addr
	ip.Port = port
	ip.Ssl = ssl
	ip.CreateAt = time.Now()
	return ip
}

func (ip ipType) encode() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(ip)
	return b.Bytes(), err
}

func (ip *ipType) decode(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(&ip)
}

func (link linkType) encode() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(link)
	return b.Bytes(), err
}

func (link *linkType) decode(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	return dec.Decode(&link)
}

func isOld(link linkType) bool {
	currentTime := time.Now()
	return currentTime.Sub(link.CheckAt) > time.Duration(720)*time.Minute
}

func grab(task pool.Task) []string {
	task.Body = cleanBody(task.Body)
	oldNumIP := numIPs
	getListIP(task.Body)
	if numIPs-oldNumIP > 0 {
		fmt.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, task.Target.String())
	}
	urls := getListURL(task)
	return urls
}

func check(task pool.Task) ipType {
	var proxy ipType
	startTimeCheck := time.Now()
	endTimeCheck := time.Now()
	proxy.Addr = task.Proxy.Hostname()
	proxy.Port = task.Proxy.Port()
	duration := endTimeCheck.Sub(startTimeCheck)
	if task.Error != nil {
		proxy.ProxyChecks++
		proxy.LastCheck = endTimeCheck
		proxy.isWork = false
		proxy.Response = duration
		proxy.LastCheck = endTimeCheck
		return proxy
	}
	strBody := string(task.Body)
	if reRemoteIP.Match(task.Body) && !strings.Contains(strBody, myIP) {
		if strings.Count(strBody, "<p>") == 1 {
			proxy.ProxyChecks = 0
			proxy.isWork = true
			proxy.isAnon = true
			proxy.Response = duration
			proxy.LastCheck = endTimeCheck
			return proxy
		}
		proxy.ProxyChecks = 0
		proxy.isWork = true
		proxy.isAnon = false
		proxy.Response = duration
		proxy.LastCheck = endTimeCheck
		return proxy
	}
	proxy.ProxyChecks++
	proxy.isWork = false
	proxy.Response = duration
	proxy.LastCheck = endTimeCheck
	return proxy
}

func backupBase() error {
	origFile, err := os.Open("gp.zip")
	if err != nil {
		return err
	}
	defer origFile.Close()
	backupName := time.Now().Format("02-01-2006-15-04-05") + ".zip"
	newFile, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(origFile, newFile)
	return err
}

func makeAddress(ip ipType) string {
	var out string
	if ip.Ssl {
		out = "https://"
	} else {
		out = "http://"
	}
	out += ip.Addr + ":" + ip.Port
	return out
}
