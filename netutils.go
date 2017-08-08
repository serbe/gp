package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func fetchBody(targetURL string, proxy Proxy) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if proxy.URL.Host != "" {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(proxy.URL),
			DisableKeepAlives: true,
		}
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			errmsg("fetchBody resp.Body.Close", err)
		}
	}()
	return body, err
}

func getHost(u string) (string, error) {
	_, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return u[:strings.LastIndex(u, "/")], nil
}

func getExternalIP() (string, error) {
	body, err := fetchBody("http://myexternalip.com/raw", Proxy{})
	if err != nil {
		return "", err
	}
	return string(body), nil
}
