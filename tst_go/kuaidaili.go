package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func kuaidaili() []string {
	var ips []string
	links := kuaidailiLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			errmsg("kuaidaili crawl", err)
			continue
		}
		ips = append(ips, kuaidailiIPS(body)...)
	}
	return ips
}

func kuaidailiLinks() []string {
	var links []string
	for page := 1; page < 5; page++ {
		links = append(links, fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d", page))
	}
	return links
}

func kuaidailiIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		errmsg("kuaidailiIPS NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() == 7 {
			ips = append(ips, strings.ToLower(td.Eq(3).Text())+"://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
		}
	})
	return ips
}
