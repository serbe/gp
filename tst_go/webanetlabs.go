package main

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func webanetlabs() []string {
	var ips []string
	body, err := crawl("http://webanetlabs.net/publ/24")
	if err != nil {
		chkErr("webanetlabs crawl", err)
		return ips
	}
	ips = append(ips, webanetlabsIPS(body)...)
	links := webanetlabsLinks(body)
	for _, link := range links {
		body, err := crawl("http://webanetlabs.net" + link)
		if err != nil {
			chkErr("webanetlabs crawl", err)
			continue
		}
		ips = append(ips, ipsFromBytes(body, "http")...)
	}
	return ips
}

func webanetlabsLinks(body []byte) []string {
	var links []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("webanetlabsLinks NewDocumentFromReader", err)
		return links
	}
	dom.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, ".txt") {
			links = append(links, href)
		}
	})
	return links
}

func webanetlabsIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("webanetlabs NewDocumentFromReader", err)
		return ips
	}
	dom.Find("p").Each(func(_ int, s *goquery.Selection) {
		sp := strings.Split(s.Text(), "\n")
		if len(sp) > 9 {
			for _, ip := range sp {
				if len(ip) > 10 {
					ips = append(ips, "http://"+ip)
				}
			}
		}
	})
	return ips
}
