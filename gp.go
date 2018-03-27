package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

var db *adb.ADB

func main() {
	getConfig()

	checkFlags()
	if (cfg.MyIPCheck && cfg.HTTPBinCheck) ||
		(cfg.MyIPCheck || cfg.HTTPBinCheck) ||
		(cfg.Target != "" && (cfg.MyIPCheck || cfg.HTTPBinCheck)) {
		log.Panic("use only one target")
	}

	// myIP, err := getMyIP()
	// if err != nil {
	// 	log.Panic("getMyIP", err)
	// }

	setTarget()
	if cfg.Target == "" {
		log.Panic("Target is empty", nil)
	}

	db = adb.InitDB("pr", "127.0.0.1:5432", "pr", "pr")

	startAppTime := time.Now()

	if useFind {
		findProxy()
	} else if useCheck {
		checkProxy(getProxyListFromDB())
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
