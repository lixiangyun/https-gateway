package util

import "github.com/astaxie/beego/logs"

func VersionGet() string {
	return "v0.3.0 20200913"
}

func init()  {
	logs.Info("version:", VersionGet())
}
