package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

var db *adb.ADB

func main() {
	checkFlags()
	if (useNoServer && (useMyIPCheck && useHttBinCheck)) ||
		(useMyIPCheck || useHttBinCheck) ||
		(useTargetURL != "" && (useMyIPCheck || useHttBinCheck)) {
		log.Println("use only one target")
		return
	}

	db = adb.InitDB("pr", "127.0.0.1:5432", "pr", "pr")

	startAppTime := time.Now()

	if !useNoServer {
		go startServer()
	}

	if useFind {
		findProxy()
	}

	if useCheck {
		checkProxy(getProxyListFromDB())
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
