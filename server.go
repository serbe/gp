package main

import (
	"fmt"
	"log"
	"net/http"
)

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p>RemoteAddr: %s</p>", r.RemoteAddr)
		for _, header := range headers {
			str := r.Header.Get(header)
			if str != "" {
				fmt.Fprintf(w, "<p>%s: %s</p>", header, str)
			}
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}
