package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

// Config all vars
type config struct {
	Target       string `json:"target"`
	FindWorkers  int64  `json:"find_workers"`
	CheckWorkers int64  `json:"check_workers"`
	Timeout      int64  `json:"timeout"`
	LogErrors    bool   `json:"log_errors"`
	LogDebug     bool   `json:"log_debug"`
	MyIPCheck    bool   `json:"my_ip_check"`
	HTTPBinCheck bool   `json:"http_bin_check"`
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
	flag.Int64Var(&cfg.CheckWorkers, "cw", cfg.CheckWorkers, "number of workers to check")
	flag.Int64Var(&cfg.FindWorkers, "fw", cfg.FindWorkers, "number of workers to find")
	flag.BoolVar(&cfg.LogDebug, "d", cfg.LogDebug, "logging debug messages")
	flag.BoolVar(&cfg.LogErrors, "e", cfg.LogErrors, "logging error messages")
	flag.BoolVar(&cfg.HTTPBinCheck, "h", cfg.HTTPBinCheck, "check working proxy with httpbin.org")
	flag.BoolVar(&cfg.MyIPCheck, "m", cfg.MyIPCheck, "check working proxy with myip.ru")
	flag.Int64Var(&cfg.Timeout, "t", cfg.Timeout, "timeout")
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
