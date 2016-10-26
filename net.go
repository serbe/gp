package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func crawl(u string) {
	fmt.Println("parse: ", u)
	mutex.Lock()
	urlList[u] = true
	mutex.Unlock()

	_, urls, ips, err := fetch(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	saveIP(ips)

	urlCount := 0

	for _, item := range urls {
		if !urlList[item] {
			urlCount++
			stringChan <- item
			fmt.Println("send ", item)
		}
	}

	if urlCount > 0 {
		fmt.Printf("found: %s %d\n", u, urlCount)
	}

	fmt.Println("finish ", u)
	return
}

func fetch(u string) ([]byte, []string, []string, error) {
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
		return nil, nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	body = cleanBody(body)
	ips := getIP(body)
	urls := getListURL(u, body)

	return body, urls, ips, err
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
