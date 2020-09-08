package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/lixiangyun/https-gateway/util"
)

var (
	Help   bool
	Debug  bool

	LogDir      string
	HealthCheck string

	Address string
	Port    int
)

func init()  {
	flag.StringVar(&LogDir, "log", "/log/https-gateway/", "log dir")
	flag.StringVar(&HealthCheck, "healthcheck", "127.0.0.1:18001", "healthcheck for docker")

	flag.IntVar(&Port, "port", 18000, "port for listen")
	flag.StringVar(&Address, "address", "0.0.0.0", "address for listen")

	flag.BoolVar(&Debug, "debug",false,"enable debug")
	flag.BoolVar(&Help,"help",false,"usage help")
}

func main()  {
	flag.Parse()

	if Help {
		flag.Usage()
		return
	}

	if !Debug {
		util.LogInit(LogDir,"console.log")
	}

	util.HealthCheckInit(HealthCheck)

	beego.SetStaticPath("/", "static")

	beego.Post("/proxy", controller)
	beego.Get("/proxy", controller)
	beego.Delete("/proxy", controller)

	beego.Post("/cert", controller)
	beego.Get("/cert", controller)
	beego.Delete("/cert", controller)

	beego.Run()
}
