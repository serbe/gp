package main

import "fmt"

func main() {
	ips := freeproxylist()
	for _, ip := range ips {
		fmt.Println(ip)
	}
	fmt.Println(len(ips))
}
