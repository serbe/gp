package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func crawl(startURL string, depth int, finished chan bool) {
	if depth <= 0 {
		finished <- false
		return
	}

	// fmt.Println("parse: ", u)
	mutex.Lock()
	urlList[startURL] = true
	mutex.Unlock()

	_, urls, ips, err := fetch(startURL)
	if err != nil {
		fmt.Println(err)
		finished <- false
		return
	}

	saveIP(ips)

	urlCount := 0

	innerFinished := make(chan bool)

	for _, u := range urls {
		if !urlList[u] {
			urlCount++
			go crawl(u, depth-1, innerFinished)
		}
	}

	if urlCount > 0 {
		fmt.Printf("found: %s %d in depth: %d\n", startURL, urlCount, depth)
	}

	for i := 0; i < urlCount; i++ {
		<-innerFinished
	}

	finished <- true

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
