package util

import "github.com/astaxie/beego/logs"

func VersionGet() string {
	return "v0.3.1 20200914"
}

func init()  {
	logs.Info("version:", VersionGet())
}
