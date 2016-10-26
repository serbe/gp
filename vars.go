package main

import (
	"flag"
	"sync"
)

var (
	stringChan chan string
	mutex      = &sync.Mutex{}

	workers = 2

	siteList = []string{
		`http://webanetlabs.net/publ/24`,
		`http://awmproxy.com/freeproxy.php`,
		`http://www.prime-speed.ru/proxy/free-proxy-list/all-working-proxies.php`,
		`http://spys.ru/proxies/`,
		`http://www.samair.ru/proxy/type-01.htm`,
		`https://www.us-proxy.org/`,
		`http://free-proxy-list.net/`,
		`http://www.proxynova.com/proxy-server-list/`,
		`http://proxyserverlist-24.blogspot.ru/`,
		`http://www.xroxy.com/proxylist.php?port=Standard&type=All_http&ssl=&country=&latency=1000&reliability=9000#table`,
		`http://www.freeproxylists.com/anonymous.html`,
		`http://www.freeproxylists.com/elite.html`,
	}

	// `http://txt.proxyspy.net/proxy.txt`,
	// `https://best-proxies.ru/proxylist/free/`,
	// `https://hidester.com/proxylist/`,
	// `http://proxyservers.pro/`,
	// `http://gatherproxy.com/`,
	// `http://proxylistdailyupdated.blogspot.ru/p/blog-page.html`,

	reURL = []string{
		`href=(?:'|")/(publ/\d{1,3}-\d{1,3})(?:'|")\s`,
		`href=(?:'|")/(freeproxylist/proxylist.*?\.txt)(?:'|")`,
		`value=(?:'|")http://awmproxy.com/(freeproxy_\d{3,12}\.txt)(?:'|")`,
		`<a href=(?:'|")/(proxies/\d{1,3}/)(?:'|")>`,
		`<a class=(?:'|")page(?:'|") href=(?:'|")(type-\d{1,3}.htm)(?:'|")>`,
		`<a href=(?:'|")/(proxy/list/order/updated/order_dir/desc/page/\d{1,3})(?:'|")>`,
		`<a href='http://proxyserverlist-24.blogspot.ru/(\d{4}/\d{1,2}/\d{1,2}-\d{1,2}-\d{1,2}-free-proxy-server-list-\d{1,6}.html#more)'`,
		`<a href=(?:'|")((?:anon|elite)/\d{1,12}.html)(?:'|")>(?:anon|elite)`,
		`href=(?:'|")/(proxylist\.php\?.+?\#table)`,
	}
	reIP        = `((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))`
	reCommaList = []string{
		`</td><td>`,
		`:`,
		`.*?(?:'|")Proxy Port \d{2,5}(?:'|")>`,
	}
	rePort  = `(\d{2,5})`
	replace = [][]string{
		// {`+(b2u1x4^t0s9)`, "1"},
		// {`+(w3q7a1^y5m3)`, "2"},
		// {`+(h8h8h8^v2b2)`, "3"},
		// {`+(y5w3r8^g7e5)`, "4"},
		// {`+(h8q7u1^v2c3)`, "5"},
		// {`+(q7x4k1^e5g7)`, "7"},
		// {`+(t0d4y5^l2p6)`, "8"},
		// {`+(g7h8m3^b2z6)`, "9"},
		// {`+(p6j0j0^o5u1)`, "0"},
		// {`1F90`, "8080"},
		// {`C38`, "3128"},
		// {`1FB6`, "8118"},
		// {`22B8`, "8888"},
		// {`270F`, "9999"},
		{`<script type=(?:'|")text/javascript(?:'|")>document\.write\((?:'|")<font class=spy2>`, ""},
		{`</font>(?:'|")`, ""},
		{`<span>`, ""},
		{`</span>`, ""},
		{`\n`, ""},
	}
	urlList map[string]bool
	ipList  map[string]bool
)

func init() {
	flag.IntVar(&workers, "w", workers, "Num of workers")
}
