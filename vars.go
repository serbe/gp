package main

import (
	"regexp"
)

var (
	user       = "pr"
	pass       = "pr"
	dbname     = "pr"
	numWorkers = 5
	timeout    = 10
	serverPort = 19091
	logErrors  = false
	useFind    = false
	useCheck   = false
	useServer  = false
	useDebug   = false
	useFile    = ""
	testLink   = ""
	addLink    = ""

	myIP       string
	reRemoteIP = regexp.MustCompile(`<p>RemoteAddr: (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):\d{1,5}<\/p>`)

	// siteList = []string{
	// 	`https://hidester.com/proxydata/php/data.php?mykey=data&offset=0&limit=1000&orderBy=latest_check&sortOrder=DESC&country=&port=&type=undefined&anonymity=undefined&ping=undefined&gproxy=2`,
	// 	`http://gatherproxy.com/embed/`,
	// 	`http://txt.proxyspy.net/proxy.txt`,
	// 	`http://webanetlabs.net/publ/24`,
	// 	`http://awmproxy.com/freeproxy.php`,
	// 	`http://www.samair.ru/proxy/type-01.htm`,
	// 	`https://www.us-proxy.org/`,
	// 	`http://free-proxy-list.net/`,
	// 	`http://www.proxynova.com/proxy-server-list/`,
	// 	`http://proxyserverlist-24.blogspot.ru/`,
	// 	`http://gatherproxy.com/`,
	// 	`https://hidemy.name/ru/proxy-list/`,
	// 	`https://hidemy.name/en/proxy-list/?type=hs&anon=34#list`,
	// 	`https://free-proxy-list.com`,
	// 	`https://free-proxy-list.com/?search=1&page=&port=&type%5B%5D=http&type%5B%5D=https&level%5B%5D=anonymous&level%5B%5D=high-anonymous&speed%5B%5D=2&speed%5B%5D=3&connect_time%5B%5D=2&connect_time%5B%5D=3&up_time=40&search=Search`,
	// 	`http://www.idcloak.com/proxylist/free-proxy-servers-list.html`,
	// 	`https://premproxy.com/list/`,
	// 	`https://proxy-list.org/english/index.php`,

	// 	`https://www.sslproxies.org/`,

	// 	// `https://best-proxies.ru/proxylist/free/`,
	// 	// `http://spys.ru/proxies/`,

	// 	// `http://www.freeproxylists.com/elite.html`,
	// 	// `http://www.freeproxylists.com/anonymous.html`,
	// 	// `http://www.xroxy.com/proxylist.php?port=Standard&type=All_http&ssl=&country=&latency=1000&reliability=9000#table`,
	// 	// `http://www.prime-speed.ru/proxy/free-proxy-list/all-working-proxies.php`,
	// 	// `http://proxyservers.pro/`,
	// }

	reURL = []string{
		`href=(?:'|")/publ/(\d{1,3}-\d{1,3})(?:'|")\s`,
		`href=(?:'|")/(freeproxylist/proxylist.*?\.txt)(?:'|")`,
		`value=(?:'|")http://awmproxy.com/(freeproxy_\d{3,12}\.txt)(?:'|")`,
		`<a href=(?:'|")/(proxies/\d{1,3}/)(?:'|")>`,
		`<a class=(?:'|")page(?:'|") href=(?:'|")(type-\d{1,3}.htm)(?:'|")>`,
		`<a href=(?:'|")/(proxy/list/order/updated/order_dir/desc/page/\d{1,3})(?:'|")>`,
		`<a href='http://proxyserverlist-24.blogspot.ru/(\d{4}/\d{1,2}/\d{1,2}-\d{1,2}-\d{1,2}-free-proxy-server-list-\d{1,6}.html#more)'`,
		`<a href=(?:'|")((?:anon|elite)/\d{1,12}.html)(?:'|")>(?:anon|elite)`,
		`href=(?:'|")/(proxylist\.php\?.+?\#table)`,
		`href=(?:'|")/(\w\w/proxy-list/\?.*?#list)`,
		`<a class=(?:'|")page(?:'|") href=(?:'|")(\d{1,3}.htm)(?:'|")>`,
		`href=(?:'|")\./(index\.php\?p=\d{1,3})(?:'|")>`,
		`<a href=(?:'|")(\d{1,3}\.htm)(?:'|")>`,
	}
	reIP        = `((?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))`
	reCommaList = []string{
		`</td><td>`,
		`:`,
		`.*?(?:'|")Proxy Port \d{2,5}(?:'|")>`,
		`(?:'|")\,(?:'|")PORT(?:'|"):`,
		`(?:'|") data\-port=(?:'|")`,
	}
	rePort  = `(\d{2,5})`
	replace = [][]string{
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")1F90(?:'|")`, ":8080"},
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")C38(?:'|")`, ":3128"},
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")1FB6(?:'|")`, ":8118"},
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")22B8(?:'|")`, ":8888"},
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")270F(?:'|")`, ":9999"},
		{`(?:'|")\,(?:'|")PROXY_LAST_UPDATE(?:'|"):(?:'|")\d{1,3} \d{1,3}(?:'|"),(?:'|")PROXY_PORT(?:'|"):(?:'|")50(?:'|")`, ":80"},
		{`<script type=(?:'|")text/javascript(?:'|")>document\.write\((?:'|")<font class=spy2>`, ""},
		{`<span style=(?:'|")display:none(?:'|")>\d{1,3}</span>`, ""},
		{`</font>(?:'|")`, ""},
		{`<span>`, ""},
		{`</span>`, ""},
		{`\n`, ""},
	}
	baseDecode = []string{
		`text/javascript(?:'|")>Proxy\((?:'|")(.*?)(?:'|")\)</script>`,
	}
	base16 = []string{
		`PROXY_IP(?:'|"):(?:'|")((?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))(?:'|"),.*?,(?:'|")PROXY_PORT(?:'|"):(?:'|")([a-fA-F0-9]{1,4})(?:'|")`,
	}
	headers = []string{
		"HTTP_VIA",
		"HTTP_X_FORWARDED_FOR",
		"HTTP_FORWARDED_FOR",
		"HTTP_X_FORWARDED",
		"HTTP_FORWARDED",
		"HTTP_CLIENT_IP",
		"HTTP_FORWARDED_FOR_IP",
		"VIA",
		"X_FORWARDED_FOR",
		"FORWARDED_FOR",
		"X_FORWARDED",
		"FORWARDED",
		"CLIENT_IP",
		"FORWARDED_FOR_IP",
		"HTTP_PROXY_CONNECTION",
	}
)
