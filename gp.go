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

	setTarget()
	if cfg.Target == "" {
		log.Panic("Target is empty", nil)
	}

	db = adb.InitDB(cfg.Database, cfg.DBAddress, cfg.Username, cfg.Password)

	startAppTime := time.Now()

	if useFind {
		findProxy()
	} else if useCheck {
		checkProxy(getProxyListFromDB())
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
