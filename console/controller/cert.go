package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/console/certbot"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/console/nginx"
	"github.com/lixiangyun/https-gateway/util"
	"github.com/lixiangyun/https-gateway/weberr"
	"time"
)

type CertInfo struct {
	Date    string `json:"date"`
	Expire  string `json:"expire"`
	Domains   string `json:"domain"`
	Email     string `json:"email"`
	Next      string `json:"next"`
	Status  string `json:"status"`
	Make    string `json:"make"`
	Detail  string `json:"detail"`
}

type CertInfoRsponse struct {
	Code    int    `json:"code"`
	Count   int    `json:"count"`
	Message string     `json:"msg"`
	Data    []CertInfo `json:"data"`
}

func CertInfo2Console(info data.CertInfo) CertInfo {
	var detail string
	if info.Cert != "" && info.Key != "" {
		detail = fmt.Sprintf("cert: %s", info.Cert)
	}

	return CertInfo{
		Date: info.Date.Format("2006-01-02"),
		Expire: info.Expire.Format("2006-01-02"),
		Next: info.Expire.AddDate(0,-1,0).Format("2006-01-02"),
		Domains: util.StringList(info.Domain),
		Email: info.Email,
		Status: util.Status(info.Status),
		Detail: detail,
		Make: info.MakeInfo,
	}
}

func CertInfo2ConsoleList(infos []data.CertInfo) []CertInfo {
	var output []CertInfo
	for _, v := range infos {
		output = append(output, CertInfo2Console(v))
	}
	return output
}

func CertInfoControllerGet(ctx *context.Context)  {
	instances, _ := data.CertQueryAll()

	var rsp CertInfoRsponse
	rsp.Code = 0
	rsp.Message = ""
	rsp.Count = len(instances)

	if len(instances) > 0 {
		begin, end := TablePage(ctx, len(instances))
		rsp.Data = CertInfo2ConsoleList(instances[begin: end])
	}

	result, _ := json.Marshal(&rsp)
	ctx.WriteString(string(result))
}

type CertAddRequest struct {
	Email      string `json:"email"`
	Domains  []string `json:"domains"`
	Auto       string `json:"update"`
}

func CertInfoControllerAdd(ctx *context.Context)  {
	var req CertAddRequest

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

	if len(req.Email) == 0 {
		logs.Error("email %s is invalid", req.Email )
		werr = weberr.WebErrMake(weberr.WEB_ERR_PARAM)
		return
	}

	nginx.NginxStop()
	defer nginx.NginxStart()

	for _, v := range req.Domains {
		cert, _ := data.CertQuery(v)
		if cert != nil {
			logs.Error("domain %s cert has been exist", v)
			werr = weberr.WebErrMake(weberr.WEB_ERR_HAS_CERT)
			return
		}

		err = util.ListenAndConnect(v, 80)
		if err != nil {
			logs.Error("domain %s connect fail", v )
			werr = weberr.WebErrMake(weberr.WEB_ERR_PARAM)
			return
		}
	}

	var auto bool
	if req.Auto == "auto" {
		auto = true
	}

	err = data.CertAdd(&data.CertInfo{
		Email: req.Email,
		Domain: req.Domains,
		Expire: time.Now(),
		Auto: auto,
		Date: time.Now(),
		MakeInfo: "making",
	})

	if err != nil {
		logs.Error("add cert fail, %s", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_ADD_CERT, err.Error())
		return
	}

	logs.Info("add cert success!")

	go makeCert(req.Domains[0])
}

func makeCert(domain string) {
	certinfo, err := data.CertQuery(domain)
	if err != nil {
		logs.Warn("cert query fail", err.Error())
		return
	}

	cert, err := certbot.CertMake(certinfo.Domain, certinfo.Email)
	if err != nil {
		certinfo.MakeInfo = err.Error()
	} else {
		certinfo.Key      = string(cert.Key)
		certinfo.Cert     = string(cert.Cert)
		certinfo.Expire   = cert.Expire
		certinfo.MakeInfo = "success"
	}

	data.CertUpdate(certinfo)
}

func updateCert() {
	certs, err := data.CertQueryAll()
	if err != nil {
		logs.Warn("cert query all fail", err.Error())
		return
	}

	err = certbot.CertUpdate()
	if err != nil {
		logs.Error("cert update fail", err.Error())
		return
	}

	for _, v := range certs {
		info, err := certbot.NewCert(v.Domain[0])
		if err != nil {
			logs.Error("load new cert fail", err.Error())
			continue
		}

		v.Key      = string(info.Key)
		v.Cert     = string(info.Cert)
		v.Expire   = info.Expire
		v.MakeInfo = "success"

		data.CertUpdate(&v)
	}

	logs.Info("cert async update success")

	nginx.NginxSync()
}

type CertDelRequest struct {
	Domain string `json:"domain"`
}

func CertInfoControllerDelete(ctx *context.Context)  {
	var req CertDelRequest

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

	cert, err := data.CertQuery(req.Domain)
	if err != nil {
		logs.Error("cert %s not exist, %s", req.Domain, err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_NOT_CERT)
		return
	}

	if cert.Status > 0 {
		logs.Error("cert %s using", req.Domain)
		werr = weberr.WebErrMake(weberr.WEB_ERR_USED_CERT)
		return
	}

	err = data.CertDelete([]string{req.Domain})
	if err != nil {
		logs.Error("cert %s not exist, %s", req.Domain, err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_NOT_CERT)
		return
	}

	logs.Info("delete cert %s success!", req.Domain)
}

func CertInfoControllerUpdate(ctx *context.Context)  {
	werr := weberr.WEB_ERR_OK
	defer func() {
		ctx.WriteString(weberr.WebErr(werr))
	}()

	updateCert()

	logs.Info("update cert success!")
}

type DomainInfoRsponse struct {
	Code    int      `json:"code"`
	Count   int      `json:"count"`
	Message string   `json:"msg"`
	Data    []string `json:"data"`
}

func DomainInfoControllerGet(ctx *context.Context)  {
	instances, _ := data.CertQueryAll()

	var rsp DomainInfoRsponse
	rsp.Code = 0
	rsp.Message = ""
	rsp.Data = make([]string, 0)

	for _,v := range instances {
		if v.Domain != nil {
			rsp.Data = append(rsp.Data, v.Domain...)
		}
	}
	rsp.Count = len(rsp.Data)

	result, _ := json.Marshal(&rsp)
	ctx.WriteString(string(result))
}