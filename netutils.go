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
	defer func() {
		chkErr("r.Body.Close", resp.Body.Close())
	}()
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

func checkTarget() bool {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}
	_, err := client.Get(cfg.Target)
	return err == nil
}
