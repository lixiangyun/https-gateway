package util

import "github.com/astaxie/beego/logs"

func VersionGet() string {
	return "v0.2.0 20200912"
}

func init()  {
	logs.Info("version:", VersionGet())
}
