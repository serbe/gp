package main

import (
	"bytes"
	"encoding/gob"
	"regexp"
	"strconv"
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

func getListURL(baseURL string, body []byte) error {
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

					resultChan <- fullURL
				}
			}
		}
	}
	return nil
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

// func newLink(addr, port string) linkType {
// 	var link linkType
// 	ip.Addr = addr
// 	ip.Port = port
// 	link.UpdateAt = time.Now()
// 	return link
// }

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
	return currentTime.Sub(link.CheckAt) > time.Duration(15*time.Second)
}
