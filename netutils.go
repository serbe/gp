package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	debugmsg("Get External IP")
	body, err := fetchBody("http://myexternalip.com/raw", Proxy{})
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func startServer() {
	debugmsg("Start server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "<p>RemoteAddr: %s</p>", r.RemoteAddr)
		if err != nil {
			errmsg("startServer fmt.Fprintf", err)
		}
		for _, header := range headers {
			str := r.Header.Get(header)
			if str != "" {
				_, err = fmt.Fprintf(w, "<p>%s: %s</p>", header, str)
				if err != nil {
					errmsg("startServer fmt.Fprintf", err)
				}
			}
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func convPort(port string, base int) string {
	portInt, _ := strconv.ParseInt(port, base, 32)
	return strconv.Itoa(int(portInt))
}
