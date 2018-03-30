package main

import "fmt"

const (
	HTTP   = "http"
	HTTPS  = "https"
	SOCKS5 = "socks5"
)

func main() {
	ips := freeproxylist()
	for _, ip := range ips {
		fmt.Println(ip)
	}
	fmt.Println(len(ips))
}
