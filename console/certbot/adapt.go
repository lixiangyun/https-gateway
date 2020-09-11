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

func NewCert(domain string) (*Cert, error) {
	key := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domain)
	cert := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domain)
	tlscfg, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logs.Error("load cert fail, %s", err.Error())
		return nil, err
	}
	logs.Info("%s -> %s", tlscfg.Leaf.NotBefore.String(), tlscfg.Leaf.NotAfter.String())
	return &Cert{CertFile: cert, CertKey: key, Expire: tlscfg.Leaf.NotAfter}, nil
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
	return nil, fmt.Errorf("return code %d, stdout:%s, stderr:%s",
		ret, cmd.Stdout(), cmd.Stderr())
}
