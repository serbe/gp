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

func fetchMyIPBody(proxy *url.URL) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	if proxy != nil {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(proxy),
			DisableKeepAlives: true,
		}
	}
	resp, err := client.Get("http://myexternalip.com/raw")
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
	debugmsg("start getExternalIP")
	body, err := fetchMyIPBody(nil)
	if err != nil {
		return "", err
	}
	debugmsg("end getExternalIP")
	return string(body), nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "<p>RemoteAddr: %s</p>", r.RemoteAddr)
	chkErr("startServer fmt.Fprintf", err)
	for _, header := range headers {
		str := r.Header.Get(header)
		if str == "" {
			continue
		}
		_, err = fmt.Fprintf(w, "<p>%s: %s</p>", header, str)
		chkErr("startServer fmt.Fprintf", err)
	}
}

func startServer() {
	debugmsg("Start server")
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func convertPort(port string) string {
	portInt, _ := strconv.ParseInt(port, 16, 32)
	return strconv.Itoa(int(portInt))
}
