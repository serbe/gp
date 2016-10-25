package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func worker(id int, jobs <-chan string, results chan<- int) {
	for u := range jobs {
		results <- parseURL(u)
	}
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

func getListURL(baseURL string, body []byte) []string {
	var urls []string
	for i := range reURL {
		host, err := getHost(baseURL)
		if err == nil {
			re := regexp.MustCompile(reURL[i])
			if re.Match(body) {
				allResults := re.FindAllSubmatch(body, -1)
				for _, result := range allResults {
					if result[1] != nil {
						fullURL := host + "/" + string(result[1])
						if !urlList[fullURL] {
							urls = append(urls, fullURL)
						}
					}
				}
			}
		}
	}
	return urls
}

func getIP(body []byte) []string {
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
					if !ipList[ip] {
						mutex.Lock()
						ipList[ip] = true
						mutex.Unlock()
						ips = append(ips, ip)
					}
				}
			}
		}
	}
	if ips == nil {
		re := regexp.MustCompile(reIP)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				if string(res[1]) != "0.0.0.0" {
					ip := string(res[1]) + ":80"
					if !ipList[ip] {
						mutex.Lock()
						ipList[ip] = true
						mutex.Unlock()
						ips = append(ips, ip)
					}
					ip = string(res[1]) + ":3128"
					if !ipList[ip] {
						mutex.Lock()
						ipList[ip] = true
						mutex.Unlock()
						ips = append(ips, ip)
					}
					ip = string(res[1]) + ":8080"
					if !ipList[ip] {
						mutex.Lock()
						ipList[ip] = true
						mutex.Unlock()
						ips = append(ips, ip)
					}
				}
			}
		}
	}
	return ips
}

func saveIP(ips []string) {
	writeSlice(ips, "ips.txt")
}

func getIPList() {
	ips := readSlice("ips.txt")
	for _, ip := range ips {
		ipList[ip] = true
	}
}

func writeSlice(slice []string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range slice {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func readSlice(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}

func existsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return createFile(file)
	}
	return true
}

func createFile(file string) bool {
	_, err := os.Create(file)
	if err != nil {
		return false
	}
	return true
}
