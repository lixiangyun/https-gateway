package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/console/nginx"
	"github.com/lixiangyun/https-gateway/util"
	"github.com/lixiangyun/https-gateway/weberr"
	"net/url"
)

type ProxyInfo struct {
	Join    string `json:"date"`
	Name    string `json:"name"`

	Https     string `json:"https"`
	Backend   string `json:"http"`
	Redirect  string `json:"redirect"`

	Status  string `json:"status"`
	Detail  string `json:"detail"`
}

type ProxyInfoRsponse struct {
	Code    int    `json:"code"`
	Count   int    `json:"count"`
	Message string      `json:"msg"`
	Data    []ProxyInfo `json:"data"`
}

func ProxyInfo2Console(info data.ProxyInfo) ProxyInfo {
	return ProxyInfo{
		Name: info.Name,
		Status: info.Status,
		Detail: fmt.Sprintf("Domain: %s", info.Cert),
		Join: info.Date.Format("2006-01-02"),
		Https: fmt.Sprintf("%d", info.HttpsPort),
		Backend: info.Backend,
		Redirect: fmt.Sprintf("%v", info.Redirct),
	}
}

func ProxyInfo2ConsoleList(infos []data.ProxyInfo) []ProxyInfo {
	var output []ProxyInfo
	for _, v := range infos {
		output = append(output, ProxyInfo2Console(v))
	}
	return output
}

func TablePage(ctx *context.Context, total int) (int, int) {
	return util.TablePage(
		util.StringToInt(ctx.Request.FormValue("page")),
		util.StringToInt(ctx.Request.FormValue("limit")),
		total)
}

func ProxyInfoControllerGet(ctx *context.Context)  {
	instances, _ := data.ProxyQueryAll()

	var rsp ProxyInfoRsponse
	rsp.Code = 0
	rsp.Message = ""
	rsp.Count = len(instances)

	if len(instances) > 0 {
		begin, end := TablePage(ctx, len(instances))
		rsp.Data = ProxyInfo2ConsoleList(instances[begin: end])
	}

	result, _ := json.Marshal(&rsp)
	ctx.WriteString(string(result))
}

type ProxyAddRequest struct {
	Name  string `json:"name"`
	Https string `json:"https"`
	Backend  string `json:"backend"`
	Redirect string `json:"redirect"`
	Domain   string `json:"domain"`
}

func ProxyInfoControllerAdd(ctx *context.Context)  {
	var req ProxyAddRequest

	werr := weberr.WEB_ERR_OK
	defer func() {
		ctx.WriteString(weberr.WebErr(werr))
	}()

	body := ctx.Input.RequestBody
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Error("json unmarshal fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_JSON_UCODER)
		return
	}

	httpsPort := util.StringToInt(req.Https)
	if httpsPort < 1 || httpsPort > 65535 {
		logs.Error("https port %d invalid", httpsPort)
		werr = weberr.WebErrMake(weberr.WEB_ERR_PARAM)
		return
	}

	err = util.CheckPortFree(httpsPort)
	if err != nil {
		logs.Error("https port %d has been used", httpsPort)
		werr = weberr.WebErrMake(weberr.WEB_ERR_HTTPS_PORT)
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 32 {
		logs.Error("proxy name %d invalid", len(req.Name) )
		werr = weberr.WebErrMake(weberr.WEB_ERR_PARAM)
		return
	}

	_, err = url.Parse(req.Backend)
	if err != nil {
		logs.Error("backend http address is invalid", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_BACKED_FAIL )
		return
	}

	instances, _ := data.ProxyQueryAll()
	for _, v := range instances {
		if v.HttpsPort == httpsPort {
			logs.Error("https port %d beed used", httpsPort)
			werr = weberr.WebErrMake(weberr.WEB_ERR_HTTPS_PORT)
			return
		}
	}

	_, err = data.CertQuery(req.Domain)
	if err != nil {
		logs.Error("domain %s cert not exist", req.Domain)
		werr = weberr.WebErrMake(weberr.WEB_ERR_NOT_CERT)
		return
	}

	var redirct bool
	if req.Redirect == "on" {
		redirct = true
	}

	err = data.ProxyAdd(&data.ProxyInfo{
		Name: req.Name,
		HttpsPort: httpsPort,
		Backend: req.Backend,
		Cert: req.Domain,
		Redirct: redirct,
		Status: "running",
	})

	if err != nil {
		logs.Error("add proxy fail, %s", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_ADD_PROXY)
		return
	}

	err = SyncProxyToNginx()
	if err != nil {
		logs.Error("sync nginx fail, %s", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_NG_PROXY)
		return
	}

	logs.Info("add proxy success!")
}

type ProxyDelRequest struct {
	Https string `json:"https"`
}

func ProxyInfoControllerDelete(ctx *context.Context)  {
	var req ProxyDelRequest

	werr := weberr.WEB_ERR_OK
	defer func() {
		ctx.WriteString(weberr.WebErr(werr))
	}()

	body := ctx.Input.RequestBody
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Error("json unmarshal fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_JSON_UCODER)
		return
	}

	err = data.ProxyDelete(req.Https)
	if err != nil {
		logs.Error("user register fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_DEL_PROXY)
		return
	}

	logs.Info("delete proxy success!")
}

func SyncProxyToNginx() error {
	proxy, err := data.ProxyQueryAll()
	if err != nil {
		return err
	}

	var items []nginx.ProxyItem
	for _,v := range proxy {
		cert, err := data.CertQuery(v.Cert)
		if err != nil {
			return err
		}

		items = append(items, nginx.ProxyItem{
			Https: v.HttpsPort,
			Name: v.Name,
			CertFile: cert.CertFile,
			CertKey: cert.CertKey,
			Backend: v.Backend,
			LogDir: fmt.Sprintf("/home/log/nginx/%s/", v.Name),
		})
	}

	return nginx.NginxConfig(&nginx.Config{
		Redirect: true,
		LogDir: fmt.Sprintf("/home/log/nginx/"),
		Proxy: items,
	})
}

func NginxInit()  {
	err := SyncProxyToNginx()
	if err != nil {
		logs.Warn("nginx init fail", err.Error())
	}
}