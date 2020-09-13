package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/console/controller"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/console/nginx"
	"github.com/lixiangyun/https-gateway/util"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Help   bool
	Debug  bool

	LogDir      string
	HealthCheck string
	Etcds       string
	HOME        string
	HTTPS       bool

	Address string
	Port    int
)

func init()  {
	flag.StringVar(&LogDir, "log", "/var/log/", "log dir")
	flag.StringVar(&HOME, "home", "/var/run/", "home dir")

	flag.StringVar(&HealthCheck, "healthcheck", "127.0.0.1:18001", "healthcheck for docker")

	flag.IntVar(&Port, "port", 18000, "port for listen")
	flag.StringVar(&Address, "address", "0.0.0.0", "address for listen")
	flag.StringVar(&Etcds, "etcds", "http://127.0.0.1:2379", "address for etcd server")

	flag.BoolVar(&HTTPS, "tls", false, "using https visit website")
	flag.BoolVar(&Debug, "debug",false,"enable debug")
	flag.BoolVar(&Help,"help",false,"usage help")
}

func BeegoConfig()  {
	if HTTPS {
		cert, key := util.NewSelfCert(nil)
		err := util.SaveToFile("./cert.pem", cert)
		if err != nil {
			logs.Error("save self signature cert fail", err.Error())
			os.Exit(1)
		}
		util.SaveToFile("./key.pem", key)
		if err != nil {
			logs.Error("save self signature key fail", err.Error())
			os.Exit(1)
		}

		beego.BConfig.Listen.EnableHTTP = false
		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.HTTPSAddr = Address
		beego.BConfig.Listen.HTTPSPort = Port
		beego.BConfig.Listen.HTTPSCertFile = "./cert.pem"
		beego.BConfig.Listen.HTTPSKeyFile = "./key.pem"
	} else {
		beego.BConfig.Listen.HTTPAddr = Address
		beego.BConfig.Listen.HTTPPort = Port
	}

	beego.BConfig.AppName = "https-gateway"
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.StaticExtensionsToGzip = []string{".js",".css",".png",".woff"}
	if Debug {
		beego.BConfig.RunMode = "dev"
	}
	beego.BConfig.EnableGzip = true
	beego.BConfig.ServerName = "https gateway console"
}

var nologin_page_type = []string {"login.html", "login"}

func ignorePath(url string, path []string) bool {
	for _, one := range path {
		if len(one) > 0 {
			if true == strings.HasSuffix(url, one) {
				return true
			}
		}
	}
	return false
}

func StaticFilterLogin(ctx * context.Context) {
	path := ctx.Request.URL.Path
	if path != "/" && !strings.HasSuffix(path, ".html") {
		return
	}
	if ignorePath(path, nologin_page_type) {
		time.Sleep(time.Second)
		return
	}
	user := controller.LoginSessionGet(ctx)
	if user == "" {
		time.Sleep(time.Second)
		ctx.Redirect(http.StatusFound,"/login.html")
	}
}

func RouterFilterLogin(ctx * context.Context) {
	path := ctx.Request.URL.Path
	if ignorePath(path, nologin_page_type) {
		time.Sleep(time.Second)
		return
	}
	user := controller.LoginSessionGet(ctx)
	if user == "" {
		time.Sleep(time.Second * 5)
		controller.NoLogin(ctx)
	}
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

	err := util.CheckPortFree(80)
	if err != nil {
		logs.Error("80 port been used")
		return
	}

	util.HealthCheckInit(HealthCheck)
	data.DataInit([]string{Etcds})

	err = nginx.NginxInit(HOME)
	if err != nil {
		logs.Error("nginx init fail, %s", err.Error())
		return
	}

	BeegoConfig()

	beego.SetStaticPath("/", "static")

	beego.Get("/domain", controller.DomainInfoControllerGet)
	beego.Get("/console", controller.ConsoleInfoControllerGet)

	beego.Put("/proxy", controller.ProxyInfoControllerUpdate)
	beego.Post("/proxy", controller.ProxyInfoControllerAdd)
	beego.Get("/proxy", controller.ProxyInfoControllerGet)
	beego.Delete("/proxy", controller.ProxyInfoControllerDelete)

	beego.Put("/cert", controller.CertInfoControllerUpdate)
	beego.Post("/cert", controller.CertInfoControllerAdd)
	beego.Get("/cert", controller.CertInfoControllerGet)
	beego.Delete("/cert", controller.CertInfoControllerDelete)

	beego.Post("/login", controller.LoginControllerPost)
	beego.Put("/password", controller.ChangePwdControllerPut)
	beego.Any("/logout", controller.LogoutControllerGet)

	beego.InsertFilter("/*", beego.BeforeStatic, StaticFilterLogin)
	beego.InsertFilter("/*", beego.BeforeRouter, RouterFilterLogin)

	beego.Run()
}
