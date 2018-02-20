package main

import (
	"flag"
	"log"
)

func checkFlags() {
	flag.BoolVar(&useCheckAll, "all", useCheckAll, "check all proxy")
	flag.BoolVar(&useCheck, "c", useCheck, "check proxy")
	flag.BoolVar(&useDebug, "d", useDebug, "show debug messages")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")
	flag.BoolVar(&useFind, "f", useFind, "find new proxy")
	flag.BoolVar(&useFUP, "fup", useFUP, "test all hosts with frequently used ports")
	flag.BoolVar(&useHttBinCheck, "h", useHttBinCheck, "check working proxy on httpbin.org")
	flag.BoolVar(&useMyIPCheck, "m", useMyIPCheck, "check working proxy on myip.ru")
	flag.BoolVar(&useNoAddLinks, "n", useNoAddLinks, "no add find links")
	flag.BoolVar(&useServer, "s", useServer, "start server")
	flag.BoolVar(&useTestScheme, "scheme", useTestScheme, "test all to scheme")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.Int64Var(&timeout, "t", timeout, "timeout")
	flag.Int64Var(&numWorkers, "w", numWorkers, "number of workers")
	flag.StringVar(&addLink, "add", addLink, "add primary link")
	flag.StringVar(&useFile, "file", useFile, "use file with proxy list")
	flag.StringVar(&testLink, "link", testLink, "link to test it")
	flag.StringVar(&targetURL, "target", targetURL, "target URL to check like 'http://127.0.0.1:12345'")
	flag.Parse()

	if addLink != "" {
		useAddLink = true
	}

	if testLink != "" {
		useTestLink = true
	}
}

func chkErr(str string, err error) {
	if err != nil {
		errmsg(str, err)
	}
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}

func debugmsg(str ...interface{}) {
	if useDebug {
		log.Println(str)
	}
}
