package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func parseURL(iter int, u string) error {
	body, err := getBody(u)
	if err != nil {
		return err
	}
	body = cleanBody(body)
	ips := getIP(body)
	saveIP(ips)
	urls := getListURL(u, body)
	if urls != nil {
		for _, u := range urls {
			parseURL(iter+1, u)
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
	uParse, err := url.Parse("http://bing.com/search?q=dotnet")
	if err != nil {
		return "", err
	}
	return uParse.Scheme + "://" + uParse.Host, nil
}
