package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func parseURL(iter int, u string) error {
	urlList[u] = true
	fmt.Println("parse: ", u, iter)
	body, err := getBody(u)
	if err != nil {
		return err
	}
	body = cleanBody(body)
	ips := getIP(body)
	fmt.Println("num of ips: ", len(ips))
	saveIP(ips)
	urls := getListURL(u, body)
	fmt.Println("num of urls: ", len(urls))
	if iter < 5 {
		if urls != nil {
			for _, u := range urls {
				if !urlList[u] {
					parseURL(iter+1, u)
				}
			}
		}
	}
	return nil
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
