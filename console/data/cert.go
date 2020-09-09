package data

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/etcdsdk"
	"github.com/lixiangyun/https-gateway/util"
	"time"
)

type CertInfo struct {
	Status     int
	Email      string
	Domain     []string
	Expire     time.Time
	Date       time.Time
	Auto       bool
	CertKey    string
	CertFile   string
}

func CertTableGet() etcdsdk.TableAPI {
	tab, err := DBAPI.NewTable("cert")
	if err != nil {
		logs.Error("attach cert table fail", err.Error())
		return nil
	}
	return tab
}

func CertDelete(domain []string) error {
	for _,v := range domain {
		CertTableGet().Delete(v)
	}
	return nil
}

func CertAdd(domain string, cert * CertInfo) error {
	cert.Date = time.Now()
	value, err := json.Marshal(cert)
	if err != nil {
		logs.Error("route json marshal fail", err.Error())
		return err
	}
	return CertTableGet().Insert(domain, value)
}

func CertQuery(domain string) (*CertInfo, error) {
	list, err := CertQueryAll()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		_, flag := util.StringListIndex(v.Domain, domain)
		if flag == true {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("domain cert not found")
}

func CertQueryAll() ([]CertInfo, error) {
	kvs, err := CertTableGet().Query()
	if err != nil {
		return nil, err
	}
	var output []CertInfo
	for _,v := range kvs {
		var proxy CertInfo
		err = json.Unmarshal(v.Value, &proxy)
		if err != nil {
			return nil, err
		}
		output = append(output, proxy)
	}
	return output, nil
}