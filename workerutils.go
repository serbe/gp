package main

import "fmt"

var (
	numWorkers = 5
)

func addTask(s string) {
	tasks <- s
	iter++
	numUrls++
}

func worker(id int, tasks chan string, quit <-chan bool) {
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			fmt.Printf("Worker %d Grab %s\n", id, task)
			grab(task)
		case <-quit:
			return
		}
	}
}

func grab(host string) {
	body, err := fetch(host)
	if err != nil {
		finishTask <- true
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
	finishTask <- true
	return
}
