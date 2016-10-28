package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func fetch(u string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
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