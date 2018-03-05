package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

var db *adb.ADB

func main() {
	checkFlags()
	if (useMyIPCheck && useHttBinCheck) ||
		(useMyIPCheck || useHttBinCheck) ||
		(targetURL != "" && (useMyIPCheck || useHttBinCheck)) {
		log.Println("use only one target")
		return
	}

	myIP, err := getMyIP()
	if err != nil {
		errmsg("getMyIP", err)
		return
	}

	setTarget(myIP)
	if targetURL == "" {
		errmsg("targetURL is empty", nil)
		return
	}

	db = adb.InitDB("pr", "127.0.0.1:5432", "pr", "pr")

	startAppTime := time.Now()

	if useFind {
		findProxy()
	}

	if useCheck {
		checkProxy(getProxyListFromDB())
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
