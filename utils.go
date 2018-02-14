package main

import (
	"flag"
	"log"
)

func checkFlags() {
	flag.Int64Var(&numWorkers, "w", numWorkers, "number of workers")
	flag.Int64Var(&timeout, "t", timeout, "timeout")
	flag.IntVar(&serverPort, "p", serverPort, "server port")
	flag.BoolVar(&useFind, "f", useFind, "find new proxy")
	flag.BoolVar(&useCheck, "c", useCheck, "check proxy")
	flag.BoolVar(&useCheckAll, "all", useCheckAll, "check all proxy")
	flag.BoolVar(&useMyIPCheck, "m", useMyIPCheck, "check working proxy on myip.ru")
	flag.BoolVar(&useServer, "s", useServer, "start server")
	flag.BoolVar(&logErrors, "e", logErrors, "logging all errors")
	flag.BoolVar(&useDebug, "d", useDebug, "show debug messages")
	flag.StringVar(&useFile, "pf", useFile, "use file with proxy list")
	flag.StringVar(&testLink, "test", testLink, "link to test it")
	flag.StringVar(&addLink, "a", addLink, "add primary link")
	flag.BoolVar(&useFUP, "fup", useFUP, "test all hosts with 4 frequently used ports")
	flag.Parse()
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
