package main

import "regexp"

func cleanBody(body []byte) []byte {
	for i := range replace {
		re := regexp.MustCompile(replace[i][0])
		if re.Match(body) {
			body = re.ReplaceAll(body, []byte(replace[i][1]))
		}
	}
	return body
}

func getListURL(baseURL string, body []byte) []string {
	var urls []string
	for i := range regexpURL {
		host, err := getHost(baseURL)
		if err == nil {
			re := regexp.MustCompile(regexpURL[i])
			if re.Match(body) {
				allResults := re.FindAllSubmatch(body, -1)
				if allResults != nil {
					for _, result := range allResults {
						if result[1] != nil {
							urls = append(urls, host+"/"+string(result[1]))
						}
					}
				}
			}
		}
	}
	return urls
}

func getIP(body []byte) []string {
	var ips []string
	for i := range commaList {
		re := regexp.MustCompile(regexpIP + commaList[i] + regexpPort)
		if re.Match(body) {
			results := re.FindAllSubmatch(body, -1)
			for _, res := range results {
				if string(res[1]) != "0.0.0.0" {
					ips = append(ips, string(res[1])+":"+string(res[2]))
				}
			}
		}
	}
	return ips
}

func saveIP(ips []string) {

}
