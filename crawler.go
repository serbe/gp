package main

import "fmt"

func crawl(u string) error {
	urlList[u] = true

	body, err := fetch(u)
	if err != nil {
		fmt.Println(u, err)
		return err
	}

	body = cleanBody(body)
	ips := getIP(body)

	urls := getListURL(u, body)

	fmt.Println("len urls ", len(urls))

	go saveIP(ips)

	urlCount := 0

	for _, item := range urls {
		if !urlList[item] {
			stringInChan <- item
			numUrls++
			urlCount++
		}
	}

	if urlCount > 0 {
		fmt.Printf("in %s found %d urls\n", u, urlCount)
	}

	return nil
}

func grab() {
	for {
		select {
		case u := <-stringInChan:
			crawl(u)
			numUrls--
		}
	}
}
