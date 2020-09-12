package certbot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/proc"
	"io/ioutil"
	"time"
)

type Cert struct {
	Cert    []byte
	Key     []byte
	Expire  time.Time
}

func NewCert(domain string) (*Cert, error) {
	key, err := ioutil.ReadFile(fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	cert, err := ioutil.ReadFile(fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain))
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	tlscfg, err := tls.X509KeyPair(cert, key)
	if err != nil {
		logs.Error("load cert fail, %s", err.Error())
		return nil, err
	}

	certinfo, err := x509.ParseCertificate(tlscfg.Certificate[0])
	if err != nil {
		logs.Error("parse cert fail, %s", err.Error())
		return nil, err
	}

	value, err := json.Marshal(certinfo)
	if err != nil {
		logs.Error("json marshal fail", err.Error())
		return nil, err
	}
	logs.Info("x509 info: %s", string(value))

	return &Cert{Cert: cert, Key: key, Expire: certinfo.NotAfter}, nil
}

func CertUpdate() error {
	cmd := proc.NewCmd("certbot", "renew")
	ret := cmd.Run()
	if ret == 0 {
		return nil
	}
	return fmt.Errorf("certbot renew code %d, stdout:%s, stderr:%s",
		ret, cmd.Stdout(), cmd.Stderr())
}

func CertMake(domain []string, email string) (*Cert, error) {
	var input string
	input += fmt.Sprintf("Y\r")

	parms := []string{"certonly", "--standalone"}
	for _, v := range domain {
		parms = append(parms, "-d", v)
	}

	parms = append(parms, "-m", email, "--agree-tos")
	cmd := proc.NewCmd("certbot", parms...)
	ret := cmd.RunWithStdin(input)

	if ret == 0 {
		return NewCert(domain[0])
	}

	return nil, fmt.Errorf("certbot make code %d, stdout:%s, stderr:%s",
		ret, cmd.Stdout(), cmd.Stderr())
}
