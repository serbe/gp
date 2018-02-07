package main

import (
	"log"
	"time"
)

func main() {
	checkFlags()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	startAppTime := time.Now()

	if useFUP {
		// fup(db)
	} else {
		if useServer {
			go startServer()
		}

		if useFind {
			findProxy(db)
		}

		if useCheck {
			checkProxy(db)
		}

		if useMyIPCheck {
			checkOnMyIP(db)
		}
	}

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
