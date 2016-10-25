package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func parseURL(u string) int {
	mutex.Lock()
	urlList[u] = true
	mutex.Unlock()
	body, err := getBody(u)
	if err != nil {
		return 0
	}
	body = cleanBody(body)
	ips := getIP(body)
	saveIP(ips)
	urls := getListURL(u, body)
	fmt.Println("num of urls: ", len(urls))
	if urls != nil {
		for _, u := range urls {
			if !urlList[u] {
				jobs <- u
			}
		}
	}
	return len(ips)
}

func getBody(u string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	// if proxyADDR != "" {
	// 	proxyURL, err := url.Parse("http://" + proxyADDR)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	// }
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func getHost(u string) (string, error) {
	uParse, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	addon := ""
	if uParse.Host == "www.samair.ru" {
		addon = "/proxy"
	}
	return uParse.Scheme + "://" + uParse.Host + addon, nil
}
