package main

func rawlist() []string {
	var ips []string
	links := rawlistLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			chkErr("freeProxyListCom crawl", err)
			continue
		}
		ips = append(ips, ipsFromBytes(body, "http")...)
	}
	return ips
}

func rawlistLinks() []string {
	links := []string{
		"https://www.rmccurdy.com/scripts/proxy/good.txt",
		"http://www.proxylists.net/http_highanon.txt",
		"http://ab57.ru/downloads/proxylist.txt",
	}
	return links
}
