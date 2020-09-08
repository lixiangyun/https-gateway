package util

import (
	"io/ioutil"
	"net"
	"strings"
)

func EthGet() (map[string]string, error) {
	ips :=  make(map[string]string)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			ip, ipnet, err := net.ParseCIDR(v.String())
			if err != nil {
				continue
			}
			if ipnet.IP.To4() != nil {
				ips[byName.Name] = ip.String()
			}
		}
	}
	return ips, nil
}

func HostnameGet() string {
	hostname, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		return "unkown"
	}
	return strings.ReplaceAll(string(hostname),"\n","")
}

