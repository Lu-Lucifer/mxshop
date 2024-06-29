package main

import (
	"math/rand"
	"mxshop/app/mxshop/admin"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	//启动grpc服务项目
	admin.NewApp("admin-server").Run()

}
