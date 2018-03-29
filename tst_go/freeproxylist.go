package main

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func freeProxyList() []string {
	var ips []string
	links := freeProxyListLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			chkErr("freeProxyList crawl", err)
			continue
		}
		ips = append(ips, freeProxyListIPS(body)...)
	}
	return ips
}

func freeProxyListLinks() []string {
	var links = []string{
		"https://www.us-proxy.org/",
		"http://free-proxy-list.net/",
		"https://www.sslproxies.org/",
	}
	return links
}

func freeProxyListIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("usProxy NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() == 8 {
			if td.Eq(6).Text() == "yes" {
				ips = append(ips, "https://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
			} else {
				ips = append(ips, "http://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
			}
		}
	})
	return ips
}
