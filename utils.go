package main

import (
	"regexp"
	"strconv"
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
				if !urlList.get(fullURL) {
					urlList.set(fullURL, true)
					urls = append(urls, fullURL)
				}
			}
		}
	}
	return urls
}

func getListIP(body []byte) []string {
	var ips []string
	for i := range reCommaList {
		re := regexp.MustCompile(reIP + reCommaList[i] + rePort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				ip := string(res[1])
				port := string(res[2])
				portInt, _ := strconv.Atoi(port)
				if ip != "0.0.0.0" && portInt < 65535 {
					ip = ip + ":" + port
					mutex.Lock()
					if !ipList.get(ip) {
						numIPs++
						ipList.set(ip, true)
						ips = append(ips, ip)
					}
					mutex.Unlock()
				}
			}
		}
	}
	return ips
}

func saveIP(ips []string) error {
	return writeSlice(ips, "ips.txt")
}

func getIPList() {
	ips := readSlice("ips.txt")
	for _, ip := range ips {
		ipList.set(ip, true)
	}
	return
}
