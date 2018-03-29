package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func chkErr(str string, err error) {
	if err != nil {
		errmsg(str, err)
	}
}

func errmsg(str string, err error) {
	log.Println("Error in", str, err)
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

func ipsFromBytes(body []byte, scheme string) []string {
	var ips []string
	reIP := `((?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):\d{2,5})`
	reIPWScheme := `([http|https|socks5]://(?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):\d{2,5})`
	re := regexp.MustCompile(reIP)
	if scheme == "" {
		re = regexp.MustCompile(reIPWScheme)
	}
	if !re.Match(body) {
		return ips
	}
	results := re.FindAllSubmatch(body, -1)
	for _, res := range results {
		proxy := string(res[1])
		if scheme != "" {
			proxy = scheme + "://" + proxy
		}
		ips = append(ips, proxy)
	}
	return ips
}

func decodeBase64(src string) string {
	//	Proxy('MTE5LjI4LjIyMS4yODo4MDg4')
	src = strings.Replace(src, "Proxy('", "", -1)
	src = strings.Replace(src, "')", "", -1)
	out, _ := base64.StdEncoding.DecodeString(src)
	return string(out)
}
