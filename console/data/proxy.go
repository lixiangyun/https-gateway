package data

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/etcdsdk"
	"time"
)

type ProxyInfo struct {
	Name       string
	HttpsPort  int
	Backend    string
	Cert       string
	Redirct    bool
	Date       time.Time
	Status     string
}

func proxyTableGet() etcdsdk.TableAPI {
	tab, err := DBAPI.NewTable("proxy")
	if err != nil {
		logs.Error("attach proxy table fail", err.Error())
		return nil
	}
	return tab
}

func ProxyDelete(https string) error {
	return proxyTableGet().Delete(https)
}

func ProxyAdd(proxy * ProxyInfo) error {
	proxy.Date = time.Now()
	value, err := json.Marshal(proxy)
	if err != nil {
		logs.Error("route json marshal fail", err.Error())
		return err
	}
	return proxyTableGet().Insert(fmt.Sprintf("%d", proxy.HttpsPort), value)
}

func ProxyQuery(https int) (*ProxyInfo, error) {
	kv, err := proxyTableGet().QueryKey(fmt.Sprintf("%d", https))
	if err != nil {
		return nil, err
	}
	var proxy ProxyInfo
	err = json.Unmarshal(kv.Value, &proxy)
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}

func ProxyQueryAll() ([]ProxyInfo, error) {
	kvs, err := proxyTableGet().Query()
	if err != nil {
		return nil, err
	}
	var output []ProxyInfo
	for _,v := range kvs {
		var proxy ProxyInfo
		err = json.Unmarshal(v.Value, &proxy)
		if err != nil {
			return nil, err
		}
		output = append(output, proxy)
	}
	return output, nil
}
