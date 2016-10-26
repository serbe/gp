package main

import "fmt"

func crawl(u string) error {
	fmt.Println("parse: ", u)
	mutex.Lock()
	urlList[u] = true
	mutex.Unlock()

	body, err := fetch(u)
	if err != nil {
		fmt.Println(u, err)
		return err
	}

	body = cleanBody(body)
	ips := getIP(body)
	urls := getListURL(u, body)

	err = saveIP(ips)
	if err != nil {
		fmt.Println(err)
	}

	urlCount := 0

	for _, item := range urls {
		if !urlList[item] {
			urlCount++
			stringChan <- item
			fmt.Println("send ", item)
		}
	}

	if urlCount > 0 {
		fmt.Printf("found: %s %d\n", u, urlCount)
	}

	fmt.Println("finish ", u)
	return nil
}

func grab() {
	for i := 0; i < workers; i++ {
		go worker(i)
	}
}

func worker(i int) {
	for {
		select {
		case u := <-stringChan:
			fmt.Println(i, len(stringChan), " get ", u)
			err := crawl(u)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
