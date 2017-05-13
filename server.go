package main

import (
	"fmt"
	"log"
	"net/http"
)

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "<p>RemoteAddr: %s</p>", r.RemoteAddr)
		if err != nil {
			errmsg("startServer fmt.Fprintf", err)
		}
		for _, header := range headers {
			str := r.Header.Get(header)
			if str != "" {
				_, err = fmt.Fprintf(w, "<p>%s: %s</p>", header, str)
				if err != nil {
					errmsg("startServer fmt.Fprintf", err)
				}
			}
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}
