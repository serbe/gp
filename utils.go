package main

import (
	"encoding/base64"
	"io"
	"log"
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

func setIP(ip string, port string, base int) {
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
	if ip != "0.0.0.0" && portInt < 65535 {
		ipWithPort := ip + ":" + portStr
		if ips.get(ipWithPort).Addr == "" {
			numIPs++
			ips.set(ipWithPort, newIP(ip, port, false))
		}
	}
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
					setIP(ip, port, 10)
				}
			}
		}
	}
	for i := range base16 {
		re := regexp.MustCompile(base16[i])
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				setIP(string(res[1]), string(res[2]), 16)
			}
		}
	}
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				setIP(string(res[1]), string(res[2]), 10)
			}
		}
	}
	return
}

func newIP(addr, port string, ssl bool) IP {
	var ip IP
	ip.Address = addr
	ip.Port = port
	ip.Ssl = ssl
	ip.CreateAt = time.Now()
	return ip
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

func check(task pool.Task) IP {
	var proxy IP
	proxy.Address = task.Proxy.Hostname()
	proxy.Port = task.Proxy.Port()
	proxy.UpdateAt = time.Now()
	proxy.IsWork = false
	proxy.IsAnon = false
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

func backupBase() error {
	origFile, err := os.Open("gp.zip")
	if err != nil {
		return err
	}
	defer func() {
		err = origFile.Close()
		if err != nil {
			errmsg("backupBase origFile.Close", err)
		}
	}()
	backupName := time.Now().Format("02-01-2006-15-04-05") + ".zip"
	newFile, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer func() {
		err = newFile.Close()
		if err != nil {
			errmsg("backupBase newFile.Close", err)
		}
	}()
	_, err = io.Copy(newFile, origFile)
	if err != nil {
		errmsg("backupBase io.Copy", err)
	}
	err = newFile.Sync()
	return err
}

func makeAddress(ip IP) string {
	var out string
	if ip.Ssl {
		out = "https://"
	} else {
		out = "http://"
	}
	out += ip.Address + ":" + ip.Port
	return out
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}
