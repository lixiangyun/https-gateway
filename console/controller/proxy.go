package controller

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/util"
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
		Join: info.Date.Format("2006-01-02 15:04:05"),
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