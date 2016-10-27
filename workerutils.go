package main

import (
	"fmt"
)

var (
	numWorkers = 2
)

func worker(tasks chan string, quit <-chan bool) {
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			grab(task)
		case <-quit:
			numWorkWorkers--
			return
		}
	}
}

func grab(host string) {
	fmt.Println("Grab :", host)

	body, err := fetch(host)
	if err != nil {
		return
	}

	body = cleanBody(body)

	ips := getIP(body)

	urls := getListURL(host, body)

	saveIP(ips)

	for _, item := range urls {
		if !urlList[item] {
			urlList[item] = true
			crawlChan <- item
		}
	}
	return
}
