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
	Key        string
	Cert       string
	MakeInfo   string
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

func CertUsed(domain string, add int) error {
	cert, err := CertQuery(domain)
	if err != nil {
		return err
	}
	cert.Status += add
	return CertUpdate(cert)
}

func CertUpdate(cert * CertInfo) error {
	value, err := json.Marshal(cert)
	if err != nil {
		logs.Error("route json marshal fail", err.Error())
		return err
	}
	return CertTableGet().Update(cert.Domain[0], value)
}

func CertAdd(cert * CertInfo) error {
	if len(cert.Domain) < 1 {
		return fmt.Errorf("domain is null")
	}
	cert.Date = time.Now()
	value, err := json.Marshal(cert)
	if err != nil {
		logs.Error("route json marshal fail", err.Error())
		return err
	}
	return CertTableGet().Insert(cert.Domain[0], value)
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
