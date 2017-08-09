package main

import (
	"log"
	"time"
)

func main() {
	checkFlags()
	initDB()

	startAppTime = time.Now()

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
