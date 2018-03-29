package main

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func freeProxyListCom() []string {
	var ips []string
	links := freeProxyListlComLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			chkErr("freeProxyListCom crawl", err)
			continue
		}
		ips = append(ips, freeProxyListComIPS(body)...)
	}
	return ips
}

func freeProxyListlComLinks() []string {
	var links []string
	for page := 1; page < 5; page++ {
		links = append(links, fmt.Sprintf("https://free-proxy-list.com/?page=%d", page))
	}
	links = append(links, `https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=3&up_time=60&search=Search`)
	return links
}

func freeProxyListComIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("freeProxyListComIPS NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() == 11 {
			ips = append(ips, td.Eq(8).Text()+"://"+td.Eq(0).Text()+":"+td.Eq(2).Text())
		}
	})
	return ips
}
