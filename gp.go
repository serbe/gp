package main

func main() {
	existsFile("ips.txt")
	urlList = make(map[string]bool)
	ipList = make(map[string]bool)
	getIPList()
	for _, u := range siteList {
		parseURL(1, u)
	}
}
