package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getMyIP() (string, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
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

func getHost(u string) (string, error) {
	h, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return h.Scheme + "://" + h.Host, err
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

func cleanBody(body []byte) []byte {
	for i := range replace {
		re := regexp.MustCompile(replace[i][0])
		if re.Match(body) {
			body = re.ReplaceAll(body, []byte(replace[i][1]))
		}
	}
	if useTestLink && useDebug {
		ioutil.WriteFile("tmp.html", body, 0644)
	}
	return body
}

func decodeIP(src []byte) (string, string, error) {
	out, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return "", "", err
	}
	split := strings.Split(string(out), ":")
	if len(split) == 2 {
		return split[0], split[1], nil
	}
	return "", "", err
}

func getTarget(targetIP string) string {
	target := useTargetURL
	if useTargetURL == "" {
		if useMyIPCheck {
			target = "http://myip.ru/"
		} else if useHttBinCheck {
			target = "http://httpbin.org/get?show_env=1"
		} else {
			target = fmt.Sprintf("http://%s:%d/", targetIP, serverPort)
		}
	}
	return target
}
