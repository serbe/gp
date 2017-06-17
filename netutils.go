package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
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
	body, err := fetchBody("http://myexternalip.com/raw", ipType{})
	if err != nil {
		return "", err
	}
	return string(body), nil
}
