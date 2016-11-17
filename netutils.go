package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func fetchBody(targetURL string, proxy ipType) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if proxy.Addr != "" {
		var (
			proxyURL *url.URL
			err      error
		)
		if proxy.Ssl {
			proxyURL, err = url.Parse("https://" + proxy.Addr + ":" + proxy.Port)
		} else {
			proxyURL, err = url.Parse("http://" + proxy.Addr + ":" + proxy.Port)
		}
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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

func getExternalIP() string {
	body, err := fetchBody("http://myexternalip.com/raw", ipType{})
	if err != nil {
		panic(err)
	}
	return string(body)
}
