package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func proxylistorg() []string {
	var ips []string
	links := proxylistorgLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			errmsg("proxylistorg crawl", err)
			continue
		}
		ips = append(ips, proxylistorgIPS(body)...)
	}
	return ips
}

func proxylistorgIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		errmsg("proxylistorg NewDocumentFromReader", err)
		return ips
	}
	dom.Find("ul").Each(func(_ int, s *goquery.Selection) {
		scheme := strings.ToLower(s.Find("li.https").Text())
		ip := decodeBase64(s.Find("li.proxy").Text())
		if ip != "" && (scheme == "http" || scheme == "https" || scheme == "socks5") {
			ips = append(ips, scheme+"://"+ip)
		}
	})
	return ips
}

func proxylistorgLinks() []string {
	var links []string
	for page := 1; page < 5; page++ {
		links = append(links, fmt.Sprintf("https://proxy-list.org/english/index.php?p=%d", page))
	}
	links = append(links, `https://proxy-list.org/english/search.php?search=anonymous-and-elite&country=any&type=anonymous-and-elite&port=any&ssl=any`)
	return links
}
