package data

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/etcdsdk"
)

func SysTableGet() etcdsdk.TableAPI {
	tab, err := DBAPI.NewTable("sys")
	if err != nil {
		logs.Error("attach sys table fail", err.Error())
		return nil
	}
	return tab
}

type SysStat struct {
	UpFlowSize   int64
	DownFlowSize int64
	ReqeustCnt   int
}

func SysStatUpdate(sys SysStat) error {
	value, err := json.Marshal(&sys)
	if err != nil {
		return err
	}
	return UserTableGet().Update("stat", value)
}

func SysStatGet() SysStat {
	kv, err := UserTableGet().QueryKey("stat")
	if err != nil {
		return SysStat{}
	}
	sys := new(SysStat)
	err = json.Unmarshal(kv.Value, sys)
	if err != nil {
		return SysStat{}
	}
	return *sys
}
