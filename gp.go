package main

import (
	"log"
	"time"
)

func main() {
	cfg := initVars()

	startAppTime := time.Now()

	p := newPool(cfg)

	p.getHostList()

	p.run()

	log.Printf("Total time: %v\n", time.Since(startAppTime))
}
