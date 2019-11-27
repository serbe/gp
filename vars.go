package main

import (
	"regexp"
)

// protocols
const (
	HTTP = "http"
	// HTTPS  = "https"
	// SOCKS5 = "socks5"
)

var (
	cfg config

	useFUP      = false
	useFind     = false
	useCheck    = false
	useCheckAll = false
	// useAddLink     = false
	useNoAddLinks  = false
	useTestLink    = false
	useCheckScheme = false
	testFile       = ""
	testLink       = ""
	primaryLink    = ""

	reRemoteIP = regexp.MustCompile(`<p>RemoteAddr: (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):\d{1,5}</p>`)
	reMyIP     = regexp.MustCompile(`<td>(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})</td>`)
)
