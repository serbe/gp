package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// Config all vars
type config struct {
	HTTPBinCheck bool   `json:"http_bin_check"`
	LogErrors    bool   `json:"log_errors"`
	LogDebug     bool   `json:"log_debug"`
	MyIPCheck    bool   `json:"my_ip_check"`
	Workers      int64  `json:"workers"`
	Timeout      int64  `json:"timeout"`
	Database     string `json:"database"`
	DBAddress    string `json:"db_address"`
	Password     string `json:"password"`
	Target       string `json:"target"`
	Username     string `json:"username"`
}

func getConfig() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Panic("getConfig ReadFile", err)
	}
	if err = json.Unmarshal(file, &cfg); err != nil {
		log.Panic("getConfig Unmarshal", err)
	}
}

func checkFlags() {
	flag.BoolVar(&cfg.LogDebug, "d", cfg.LogDebug, "logging debug messages")
	flag.BoolVar(&cfg.LogErrors, "e", cfg.LogErrors, "logging error messages")
	flag.BoolVar(&cfg.HTTPBinCheck, "h", cfg.HTTPBinCheck, "check working proxy with httpbin.org")
	flag.BoolVar(&cfg.MyIPCheck, "m", cfg.MyIPCheck, "check working proxy with myip.ru")
	flag.Int64Var(&cfg.Timeout, "t", cfg.Timeout, "timeout in millisecond")
	flag.Int64Var(&cfg.Workers, "w", cfg.Workers, "number of workers")
	flag.StringVar(&cfg.Target, "target", cfg.Target, "target URL to check like 'http://127.0.0.1:12345/target'")

	flag.BoolVar(&useCheckAll, "all", useCheckAll, "check all proxy")
	flag.BoolVar(&useCheck, "c", useCheck, "check proxy")
	flag.BoolVar(&useFind, "f", useFind, "find new proxy")
	flag.BoolVar(&useFUP, "fup", useFUP, "test all hosts with frequently used ports")
	flag.BoolVar(&useNoAddLinks, "test", useNoAddLinks, "no add find links")
	flag.BoolVar(&useCheckScheme, "scheme", useCheckScheme, "check all http to https and socks5 scheme ")

	flag.StringVar(&primaryLink, "p", primaryLink, "add primary link")
	flag.StringVar(&testFile, "file", testFile, "use file with proxy list")
	flag.StringVar(&testLink, "link", testLink, "link to test it")

	flag.Parse()

	if primaryLink != "" {
		useAddLink = true
	}

	if testLink != "" {
		useTestLink = true
	}
}

func chkErr(str string, err error) {
	if err != nil {
		errmsg(str, err)
	}
}

func errmsg(str string, err error) {
	if cfg.LogErrors {
		log.Println("Error in", str, err)
	}
}

func debugmsg(str ...interface{}) {
	if cfg.LogDebug {
		log.Println(str)
	}
}

func decodeBase64(src string) string {
	src = strings.Replace(src, "Proxy('", "", -1)
	src = strings.Replace(src, "')", "", -1)
	out, _ := base64.StdEncoding.DecodeString(src)
	return string(out)
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
