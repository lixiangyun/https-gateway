package nginx

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/proc"
	"github.com/lixiangyun/https-gateway/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
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
	PIDFile  string
	Redirect bool
	LogDir   string
	Proxy  []ProxyItem
}

var NGINX_HOME string
var NGINX_CONFIG_TEMPLATE string
var NGINX_CONFIG_PATH string
var NGINX_CONFIG_TMP_PATH string
var NGINX_PID string

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

	cmd := proc.NewCmd("nginx", "-t", "-c", cfg)
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx check fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}

func NginxStop() error {
	if !NginxRunning() {
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
	body, err := ioutil.ReadFile(NGINX_PID)
	if err != nil {
		logs.Warn("read nginx pid fail, %s", err.Error())
		return false
	}
	body2 := strings.ReplaceAll(string(body), "\r","")
	body2 = strings.ReplaceAll(body2, "\n","")
	pid, err := strconv.Atoi(body2)
	if err != nil {
		logs.Warn("[%s] atoi fail, %s", string(body), err.Error())
		return false
	}
	if pid > 0 && pid < 65535 {
		return true
	}
	return false
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

	var parms []string
	if NginxRunning() {
		parms = []string{"-s", "reload", "-c", cfg}
	} else {
		parms = []string{"-c", cfg}
	}

	logs.Info("nginx %v", parms)

	cmd := proc.NewCmd("nginx", parms...)
	retcode := cmd.Run()
	if retcode == 0 {
		return nil
	}

	return fmt.Errorf("nginx start fail, stdout:%s, stderr:%s", cmd.Stdout(), cmd.Stderr())
}


func SyncProxyToNginx() error {
	proxy, err := data.ProxyQueryAll()
	if err != nil {
		return err
	}

	var items []ProxyItem
	for _,v := range proxy {

		cert, err := data.CertQuery(v.Cert)
		if err != nil {
			return err
		}

		certFile := fmt.Sprintf("%s/%s/cert.pem", NGINX_HOME, v.Name )
		keyFile := fmt.Sprintf("%s/%s/key.pem", NGINX_HOME, v.Name )

		util.SaveToFile(certFile, []byte(cert.Cert))
		util.SaveToFile(keyFile, []byte(cert.Key))

		logDirs := fmt.Sprintf("%s/%s", NGINX_HOME, v.Name)
		os.MkdirAll(logDirs, 0644)

		items = append(items, ProxyItem{
			Https: v.HttpsPort,
			Name: v.Name,
			CertFile: certFile,
			CertKey: keyFile,
			Backend: v.Backend,
			LogDir: logDirs,
		})
	}

	return NginxConfig(&Config{
		PIDFile: NGINX_PID,
		Redirect: true,
		LogDir: NGINX_HOME,
		Proxy: items,
	})
}

func NginxSync()  {
	err := SyncProxyToNginx()
	if err != nil {
		logs.Warn("nginx sync fail", err.Error())
	}else {
		logs.Info("nginx sync success")
	}
}

func NginxInit(home string) error {
	err := os.MkdirAll(home, 0644)
	if err != nil {
		return err
	}

	NGINX_HOME = home
	NGINX_CONFIG_TEMPLATE = NGINX_HOME + "/nginx.conf.template"
	NGINX_CONFIG_PATH = NGINX_HOME + "/nginx.conf"
	NGINX_CONFIG_TMP_PATH = NGINX_HOME + "/nginx.conf.test"
	NGINX_PID = NGINX_HOME + "/nginx.pid"

	err = util.CopyToFile(NGINX_CONFIG_TEMPLATE, "./nginx.conf.template")
	if err != nil {
		return err
	}

	go func() {
		for  {
			if !NginxRunning() {
				NginxSync()
			}
			time.Sleep(time.Minute)
		}
	}()

	return nil
}
