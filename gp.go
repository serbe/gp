package main

import (
	"log"
	"time"

	"github.com/serbe/adb"
)

func main() {
	checkFlags()
	db, err := adb.InitDB("pr", "127.0.0.1:5432", "pr", "pr")
	if err != nil {
		errmsg("initDB", err)
		return
	}

	if useDebug {
		db.UseDebug(true)
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
