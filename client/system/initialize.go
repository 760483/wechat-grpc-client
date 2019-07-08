package client_system

import (
	"flag"
	"time"
	"github.com/patrickmn/go-cache"
)

var logPath *string = flag.String("log", "./runtime/log/run.log", "Use -log <log output path>")

func init() {
	flag.Parse()
	cacheDriver = cache.New(30*time.Second, 30*24*time.Hour)
}
