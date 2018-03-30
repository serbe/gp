package main

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func freeproxylist() []string {
	var ips []string
	links := freeproxylistLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			errmsg("freeproxylist crawl", err)
			continue
		}
		ips = append(ips, freeproxylistIPS(body)...)
	}
	return ips
}

func freeproxylistLinks() []string {
	var links = []string{
		"https://www.us-proxy.org/",
		"http://free-proxy-list.net/",
		"https://www.sslproxies.org/",
		"https://www.socks-proxy.net/",
	}
	return links
}

func freeproxylistIPS(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		errmsg("freeproxylistIPS NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() == 8 {
			if strings.ToLower(td.Eq(4).Text()) == "socks5" {
				ips = append(ips, "socks5://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
			} else if strings.ToLower(td.Eq(4).Text()) != "socks4" {
				if td.Eq(6).Text() == "yes" {
					ips = append(ips, "https://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
				} else {
					ips = append(ips, "http://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
				}
			}
		}
	})
	return ips
}
