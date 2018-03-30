package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func gatherproxycom() []string {
	var ips []string
	links := gatherproxycomLinks()
	for _, link := range links {
		body, err := crawl(link)
		if err != nil {
			errmsg("gatherproxycom crawl", err)
			continue
		}
		ips = append(ips, gatherproxycomIPS(body)...)
	}
	return ips
}

func gatherproxycomLinks() []string {
	links := []string{
		"http://www.gatherproxy.com/embed/",
		"http://www.gatherproxy.com/embed/?t=Elite&p=&c=",
		"http://www.gatherproxy.com/embed/?t=Anonymous&p=&c=",
	}
	return links
}

func gatherproxycomIPS(body []byte) []string {
	var ips []string
	reAddr := `"PROXY_IP":"(.+?)".+?"PROXY_PORT":"(.+?)"`
	re := regexp.MustCompile(reAddr)
	if !re.Match(body) {
		return ips
	}
	results := re.FindAllSubmatch(body, -1)
	for _, res := range results {
		ip := string(res[1])
		hexPort := string(res[2])
		port, err := strconv.ParseInt(hexPort, 16, 64)
		if err != nil {
			continue
		}
		ips = append(ips, fmt.Sprintf("http://%v:%d", ip, port))
	}
	return ips
}
