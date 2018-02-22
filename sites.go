package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func freeProxyListlLinks() []string {
	var links []string
	for page := 1; page < 5; page++ {
		links = append(links, fmt.Sprintf("https://free-proxy-list.com/?page=%d", page))
	}
	return links
}

func kuaidailiLinks() []string {
	var links []string
	for page := 2; page < 5; page++ {
		links = append(links, fmt.Sprintf("https://www.kuaidaili.com/free/inha/%d", page))
	}
	return links
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

func freeProxyList(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("freeProxyList NewDocumentFromReader", err)
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

func kuaidaili(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("kuaidaili NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		if td.Length() == 7 {
			ips = append(ips, td.Eq(3).Text()+"://"+td.Eq(0).Text()+":"+td.Eq(1).Text())
		}
	})
	return ips
}

func usProxy(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("usProxy NewDocumentFromReader", err)
		return ips
	}
	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
		td := s.Find("td")
		log.Println(td.Length())
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

func webanetlabs(body []byte) []string {
	var ips []string
	r := bytes.NewReader(body)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		chkErr("usProxy NewDocumentFromReader", err)
		return ips
	}
	dom.Find("p").Each(func(_ int, s *goquery.Selection) {
		sp := strings.Split(s.Text(), "\n")
		if len(sp) > 9 {
			for _, ip := range sp {
				if len(ip) > 10 {
					ips = append(ips, ip)
				}
			}
		}
	})
	return ips
}
