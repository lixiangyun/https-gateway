package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/lixiangyun/https-gateway/console/controller"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/util"
)

var (
	Help   bool
	Debug  bool

	LogDir      string
	HealthCheck string
	Etcds       string

	Address string
	Port    int
)

func init()  {
	flag.StringVar(&LogDir, "log", "/var/log/https-gateway/", "log dir")
	flag.StringVar(&HealthCheck, "healthcheck", "127.0.0.1:18001", "healthcheck for docker")

	flag.IntVar(&Port, "port", 18000, "port for listen")
	flag.StringVar(&Address, "address", "0.0.0.0", "address for listen")
	flag.StringVar(&Etcds, "etcds", "http://127.0.0.1:2379", "address for etcd server")

	flag.BoolVar(&Debug, "debug",false,"enable debug")
	flag.BoolVar(&Help,"help",false,"usage help")
}


func BeegoConfig()  {
	//beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPAddr = Address
	beego.BConfig.Listen.HTTPPort = Port

	beego.BConfig.AppName = "https-gateway"
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.StaticExtensionsToGzip = []string{".js",".css",".png",".woff"}
	if Debug {
		beego.BConfig.RunMode = "dev"
	}
	beego.BConfig.EnableGzip = true
	beego.BConfig.ServerName = "https gateway console"
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
	data.DataInit([]string{Etcds})

	BeegoConfig()

	beego.SetStaticPath("/", "static")

	beego.Post("/proxy", controller.ProxyInfoControllerAdd)
	beego.Get("/proxy", controller.ProxyInfoControllerGet)
	beego.Delete("/proxy", controller.ProxyInfoControllerDelete)

	beego.Post("/cert", controller)
	beego.Get("/cert", controller)
	beego.Delete("/cert", controller)

	beego.Run()
}