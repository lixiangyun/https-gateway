package certbot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/proc"
	"time"
)

type Cert struct {
	CertFile string
	CertKey  string
	Expire   time.Time
}

func NewCert(domain string, output string) (*Cert, error) {
	key := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain)
	cert := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain)

	tlscfg, err := tls.LoadX509KeyPair(cert, key)
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

	return &Cert{CertFile: cert, CertKey: key, Expire: certinfo.NotAfter}, nil
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
		return NewCert(domain[0], cmd.Stdout())
	}
	return nil, fmt.Errorf("certbot make code %d, stdout:%s, stderr:%s",
		ret, cmd.Stdout(), cmd.Stderr())
}
