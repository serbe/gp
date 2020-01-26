package main

import (
	"errors"
	"log"
)

// Config all vars
type config struct {
	HTTPBinCheck   bool   `json:"http_bin_check"`
	LogErrors      bool   `json:"log_errors"`
	LogDebug       bool   `json:"log_debug"`
	MyIPCheck      bool   `json:"my_ip_check"`
	UseFUP         bool   `json:"use_fup"`
	UseFind        bool   `json:"use_find"`
	UseCheck       bool   `json:"use_check"`
	UseCheckAll    bool   `json:"use_check_all"`
	UseCheckScheme bool   `json:"use_check_scheme"`
	Workers        int64  `json:"workers"`
	Timeout        int64  `json:"timeout"`
	DatabaseURL    string `json:"database"`
	Target         string `json:"target"`
	TestLink       string `json:"test_link"`
	isUpdate       bool
	useTestLink    bool
	myIP           string
	// db             adb.DB
}

// protocols
const (
	HTTP = "http"
	// HTTPS  = "https"
	// SOCKS5 = "socks5"
)

var (
	logErrors bool
	logDebug  bool

	errEmptyHostname = errors.New("error: empty hostname")
	// errNotRun        = errors.New("error: pool is not running")
	// errNotWait       = errors.New("error: pool is not waiting tasks")
	// useFUP      = false
	// useFind     = false
	// useCheck    = false
	// useCheckAll = false
	// // useAddLink     = false
	// useNoAddLinks  = false
	// useTestLink    = false
	// useCheckScheme = false
	// testFile       = ""
	// testLink       = ""
	// primaryLink    = ""
)

func initVars() config {
	cfg := getConfig()

	checkFlags(&cfg)
	if (cfg.MyIPCheck && cfg.HTTPBinCheck) ||
		(cfg.MyIPCheck || cfg.HTTPBinCheck) ||
		(cfg.Target != "" && (cfg.MyIPCheck || cfg.HTTPBinCheck)) {
		log.Panic("use only one target")
	}

	setTarget(&cfg)
	if cfg.Target == "" {
		log.Panic("Target is empty", nil)
	}

	if !checkTarget(&cfg) {
		log.Panic("Target ", cfg.Target, " unavailable")
	}

	err := getMyIP(&cfg)
	if err != nil {
		log.Panic("Not get my IP", err)
	}

	return cfg
}
