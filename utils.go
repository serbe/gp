package main

import (
	"bytes"
	"compress/flate"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func getListURL(baseURL string, body []byte) []string {
	var urls []string
	for i := range reURL {
		host, err := getHost(baseURL)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if re.Match(body) {
			allResults := re.FindAllSubmatch(body, -1)
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
						ips.set(ipWithPort, newIP(ip, port))
					}
				}
			}
		}
	}
	return
}

func newIP(addr, port string) ipType {
	var ip ipType
	ip.Addr = addr
	ip.Port = port
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

func compress(inputFile, outputFile string) error {
	i, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer i.Close()
	o, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer o.Close()
	f, err := flate.NewWriter(o, flate.BestCompression)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, i)
	return err
}

func decompress(inputFile, outputFile string) error {
	i, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer i.Close()
	f := flate.NewReader(i)
	if err != nil {
		return err
	}
	defer f.Close()
	o, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer o.Close()
	_, err = io.Copy(o, f)
	return err
}

// host string
func grab(args ...interface{}) interface{} {
	var urls []string
	host := args[0].(string)
	fmt.Printf("Start grab %s\n", host)
	body, err := fetchBody(host, ipType{})
	if err != nil {
		return urls
	}
	body = cleanBody(body)
	oldNumIP := numIPs
	getListIP(body)
	if numIPs-oldNumIP > 0 {
		fmt.Printf("Find %d new ip address in %s\n", numIPs-oldNumIP, host)
	}
	urls = getListURL(host, body)
	return urls
}

// proxy ipType
func check(args ...interface{}) interface{} {
	targetURL := fmt.Sprintf("http://93.170.123.221:%d/", proxyPort)
	proxy := args[0].(ipType)
	startTimeCheck := time.Now()
	body, err := fetchBody(targetURL, proxy)
	endTimeCheck := time.Now()
	duration := endTimeCheck.Sub(startTimeCheck)
	if err != nil {
		proxy.ProxyChecks++
		proxy.LastCheck = endTimeCheck
		proxy.isWork = false
		proxy.Response = duration
		proxy.LastCheck = endTimeCheck
		return proxy
	}
	strBody := string(body)
	if reRemoteIP.Match(body) && !strings.Contains(strBody, myIP) {
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
	origFile, err := os.Open("gp.gz")
	if err != nil {
		return err
	}
	defer origFile.Close()
	backupName := time.Now().Format("02-01-2006-15-04-05") + ".gz"
	newFile, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(origFile, newFile)
	return err
}
