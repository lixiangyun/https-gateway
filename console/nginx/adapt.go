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

func NginxConfig(name string, ctx *Config) (error) {
	_, err := os.Stat(name)
	if err == nil {
		err = os.Remove(name)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	temp, err := template.ParseFiles("./nginx.conf.template")
	if err != nil {
		return err
	}

	err = temp.Execute(file, ctx)
	if err != nil {
		logs.Error("executing template fail", err.Error())
		return err
	}

	return nil
}

func NginxCheck(name string) error {
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

func NginxReload(name string) error {
	cfg, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	cmd :=proc.NewCmd("nginx", "-s", "reload", "-c", cfg)
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx reload fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}