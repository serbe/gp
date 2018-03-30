package main

// func proxynovacom() []string {
// 	var ips []string
// 	body, err := crawl("https://www.proxynova.com/proxy-server-list/")
// 	if err != nil {
// 		errmsg("proxynovacom crawl", err)
// 		return ips
// 	}
// 	ips = append(ips, proxynovacomIPS(body)...)
// 	return ips
// }

// func proxynovacomIPS(body []byte) []string {
// 	var ips []string
// 	r := bytes.NewReader(body)
// 	dom, err := goquery.NewDocumentFromReader(r)
// 	if err != nil {
// 		errmsg("proxynovacomIPS NewDocumentFromReader", err)
// 		return ips
// 	}
// 	dom.Find("tr").Each(func(_ int, s *goquery.Selection) {
// 		// ip, exist := s.Find("td abbr").Attr("title")
// 		ip := s.Find("td abbr").Text()
// 		port := s.Find("td a").Text()
// 		// fmt.Println(ip, port)
// 		fmt.Println(ip)
// 		// if exist && port != "" {
// 		ips = append(ips, "http://"+ip+":"+port)
// 		// }
// 	})
// 	return ips
// }
