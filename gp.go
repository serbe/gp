package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

var db adb.DB

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

	if !checkTarget() {
		log.Panic("Target", cfg.Target, "unavailable")
	}

	db = adb.InitDB(cfg.DatabaseURL)

	startAppTime := time.Now()

	if useFind {
		findProxy()
	} else if useCheck {
		checkProxy(getProxyListFromDB(), true)
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
