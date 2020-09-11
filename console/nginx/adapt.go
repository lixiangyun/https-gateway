package nginx

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/proc"
	"os"
	"path/filepath"
	"text/template"
)

type ProxyItem struct {
	Https    int
	Name     string
	CertFile string
	CertKey  string
	Backend  string
	LogDir   string
}

type Config struct {
	Redirect bool
	LogDir   string
	Proxy  []ProxyItem
}

const NGINX_CONFIG_TEMPLATE = "/home/binary/nginx.conf.template"
const NGINX_CONFIG_PATH = "/home/binary/nginx.conf"
const NGINX_CONFIG_TMP_PATH = "/home/binary/nginx.conf.tmp"
const NGINX_PID = "/home/nginx.pid"

func NginxConfig(ctx *Config) error {
	_, err := os.Stat(NGINX_CONFIG_TMP_PATH)
	if err == nil {
		err = os.Remove(NGINX_CONFIG_TMP_PATH)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(NGINX_CONFIG_TMP_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	temp, err := template.ParseFiles(NGINX_CONFIG_TEMPLATE)
	if err != nil {
		return err
	}

	err = temp.Execute(file, ctx)
	if err != nil {
		logs.Error("executing template fail", err.Error())
		return err
	}

	err = nginxTest(NGINX_CONFIG_TMP_PATH)
	if err != nil {
		logs.Error("nginx test config fail", err.Error())
		return err
	}

	err = os.Rename(NGINX_CONFIG_TMP_PATH, NGINX_CONFIG_PATH)
	if err != nil {
		logs.Error("rename config fail", err.Error())
		return err
	}

	return NginxStart()
}

func nginxTest(name string) error {
	cfg, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	cmd :=proc.NewCmd("nginx", "-t", "-c", cfg)
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx check fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}

func NginxStop() error {
	_, err := os.Stat(NGINX_PID)
	if err != nil {
		logs.Info("nginx has been stop")
		return nil
	}

	cmd := proc.NewCmd("nginx", "-s", "stop")
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx stop fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}

func NginxRunning() bool {
	_, err := os.Stat(NGINX_PID)
	if err != nil {
		return false
	}
	return true
}

func NginxStart() error {
	cfg, err := filepath.Abs(NGINX_CONFIG_PATH)
	if err != nil {
		return err
	}

	_, err = os.Stat(cfg)
	if err != nil {
		logs.Error("nginx config not exist", cfg)
		return err
	}

	var flag []string
	_, err = os.Stat(NGINX_PID)
	if err == nil {
		flag = []string{"-s", "reload", "-c", cfg}
	} else {
		flag = []string{"-c", cfg}
	}

	cmd := proc.NewCmd("nginx", flag...)
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx start fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}
