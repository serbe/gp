package main

// import (
// 	"bytes"
// 	"fmt"

// 	"github.com/PuerkitoBio/goquery"
// )

// func proxy5050() []string {
// 	var ips []string
// 	_, err := crawl("http://proxy50-50.blogspot.com.tr/")
// 	if err != nil {
// 		errmsg("proxy5050 crawl", err)
// 	}
// 	return ips
// }

// func proxy5050IPS(body []byte) []string {
// 	var ips []string
// 	r := bytes.NewReader(body)
// 	dom, err := goquery.NewDocumentFromReader(r)
// 	if err != nil {
// 		errmsg("proxy5050IPS NewDocumentFromReader", err)
// 		return ips
// 	}
// 	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
// 		fmt.Println(s.Text())
// 		td := s.Find("td")
// 		ip := td.Find("a").Text()
// 		port := td.Eq(2).Text()
// 		ips = append(ips, "http://"+ip+":"+port)
// 	})
// 	return ips
// }
