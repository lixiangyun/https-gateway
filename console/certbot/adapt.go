package certbot

import (
	"crypto/tls"
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

func NewCert(domain string) *Cert {
	key := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain)
	cert := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain)
	tlscfg, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logs.Error("load cert fail, %s", err.Error())
		return nil
	}
	logs.Info("%s -> %s", tlscfg.Leaf.NotBefore.String(), tlscfg.Leaf.NotAfter.String())
	return &Cert{CertFile: cert, CertKey: key, Expire: tlscfg.Leaf.NotAfter}
}

func CertMake(domain []string, email string) (*Cert, error) {
	var input string
	input += fmt.Sprintf("%s\r", email)
	input += fmt.Sprintf("A\r")
	input += fmt.Sprintf("Y\r")

	parms := []string{"certonly", "--standalone"}
	for _, v := range domain {
		parms = append(parms, "-d", v)
	}
	cmd := proc.NewCmd("certbot", parms...)
	ret := cmd.RunWithStdin(input)

	if ret == 0 {
		return NewCert(domain[0]), nil
	}

	return nil, fmt.Errorf("return code %d, stdout:%s, stderr:%s",
		ret, cmd.Stdout(), cmd.Stderr())
}
