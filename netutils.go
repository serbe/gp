package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getMyIP() (string, error) {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.Replace(string(body), "\n", "", 1)
	return ip, err
}

func setTarget() {
	if cfg.Target == "" {
		if cfg.MyIPCheck {
			cfg.Target = "http://myip.ru/"
		} else if cfg.HTTPBinCheck {
			cfg.Target = "http://httpbin.org/get?show_env=1"
		}
	}
}

func crawl(target string) ([]byte, error) {
	timeout := time.Duration(15000 * time.Millisecond)
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", "https://www.google.com/")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	ioutil.WriteFile("tmp.html", body, 0644)
	return body, err
}
