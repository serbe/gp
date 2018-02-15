package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

var db *adb.ADB

func main() {
	checkFlags()

	db = adb.InitDB("pr", "127.0.0.1:5432", "pr", "pr")

	if useDebug {
		db.ShowErrors(true)
	}

	startAppTime := time.Now()

	if useServer {
		go startServer()
	}

	if useFind {
		findProxy()
	}

	if useCheck {
		checkProxy()
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
