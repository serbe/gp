package main

import (
	"log"
	"time"
)

func main() {
	var err error
	checkFlags()
	db, err = initDB()
	if err != nil {
		errmsg("initDB", err)
		return
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
