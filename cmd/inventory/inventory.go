package main

import (
	"math/rand"
	"mxshop/app/inventory/srv"
	"os"
	"runtime"
	"time"
)

func main() {
	//rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	//启动grpc服务项目
	srv.NewApp("inventory-server").Run()

}
