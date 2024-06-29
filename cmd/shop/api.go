package main

import (
	"math/rand"
	"mxshop/app/mxshop/api"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	api.NewApp("api-server").Run()
}
