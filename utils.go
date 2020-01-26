package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func getConfig() config {
	var cfg config
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Panic("getConfig ReadFile", err)
	}
	if err = json.Unmarshal(file, &cfg); err != nil {
		log.Panic("getConfig Unmarshal", err)
	}
	logErrors = cfg.LogErrors
	logDebug = cfg.LogDebug
	return cfg
}

func checkFlags(cfg *config) {
	flag.BoolVar(&cfg.LogDebug, "d", cfg.LogDebug, "logging debug messages")
	flag.BoolVar(&cfg.LogErrors, "e", cfg.LogErrors, "logging error messages")
	flag.BoolVar(&cfg.HTTPBinCheck, "h", cfg.HTTPBinCheck, "check working proxy with httpbin.org")
	flag.BoolVar(&cfg.MyIPCheck, "m", cfg.MyIPCheck, "check working proxy with myip.ru")
	flag.Int64Var(&cfg.Timeout, "t", cfg.Timeout, "timeout in millisecond")
	flag.Int64Var(&cfg.Workers, "w", cfg.Workers, "number of workers")
	flag.StringVar(&cfg.Target, "target", cfg.Target, "target URL to check like 'http://127.0.0.1:12345/target'")

	flag.BoolVar(&cfg.UseCheckAll, "all", cfg.UseCheckAll, "check all proxy")
	flag.BoolVar(&cfg.UseCheck, "c", cfg.UseCheck, "check proxy")
	flag.BoolVar(&cfg.UseFind, "f", cfg.UseFind, "find new proxy")
	flag.BoolVar(&cfg.UseFUP, "fup", cfg.UseFUP, "test all hosts with frequently used ports")
	// flag.BoolVar(&cfg.UseNoAddLinks, "test", cfg.UseNoAddLinks, "no add find links")
	flag.BoolVar(&cfg.UseCheckScheme, "scheme", cfg.UseCheckScheme, "check all http to https and socks5 scheme ")

	// flag.StringVar(&cfg.PrimaryLink, "p", cfg.PrimaryLink, "add primary link")
	// flag.StringVar(&testFile, "file", testFile, "use file with proxy list")
	flag.StringVar(&cfg.TestLink, "link", cfg.TestLink, "link to test it")

	flag.Parse()

	// if primaryLink != "" {
	// 	useAddLink = true
	// }

	if cfg.TestLink != "" {
		cfg.useTestLink = true
	}
}

func removeDuplicates(newList, oldList []string) []string {
	var (
		mapList map[string]bool
		list    []string
	)

	mapList = make(map[string]bool)

	for i := range oldList {
		if !mapList[oldList[i]] {
			mapList[oldList[i]] = true
		}
	}
	for i := range newList {
		if !mapList[newList[i]] {
			mapList[newList[i]] = true
			list = append(list, newList[i])
		}
	}
	return list
}

func chkErr(str string, err error) {
	if err != nil {
		errmsg(str, err)
	}
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}

func debugmsg(str ...interface{}) {
	if logDebug {
		log.Println(str...)
	}
}

func setTarget(cfg *config) {
	if cfg.Target == "" {
		if cfg.MyIPCheck {
			cfg.Target = "http://myip.ru/"
		} else if cfg.HTTPBinCheck {
			cfg.Target = "http://httpbin.org/get?show_env=1"
		}
	}
}

func saveToFile(filename string, list []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	for i := range list {
		_, err := f.WriteString(fmt.Sprintln(list[i]))
		if err != nil {
			log.Panic(err)
		}
	}
}
