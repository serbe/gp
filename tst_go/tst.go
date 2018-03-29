package main

import "fmt"

func main() {
	ips := rawlist()
	for _, ip := range ips {
		fmt.Println(ip)
	}
	fmt.Println(len(ips))
}
