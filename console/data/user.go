package data

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/etcdsdk"
	"github.com/lixiangyun/https-gateway/util"
	"time"
)

func UserTableGet() etcdsdk.TableAPI {
	tab, err := DBAPI.NewTable("user")
	if err != nil {
		logs.Error("attach user table fail", err.Error())
		return nil
	}
	return tab
}

type UserInfo struct {
	Name  string
	Pwd   string
	Date  time.Time
}

func (this *UserInfo)CheckPasswod(pwd string) bool {
	return util.CheckPassword(this.Pwd, pwd)
}

func UserUpdate(name string, pwd string) error {
	user := UserInfo{
		Name: name, Pwd: util.CryptPassword(pwd), Date: time.Now(),
	}
	value, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	return UserTableGet().Update("admin", value)
}

func UserGet() (*UserInfo, error) {
	kv, err := UserTableGet().QueryKey("admin")
	if err != nil {
		return nil, err
	}
	user := new(UserInfo)
	err = json.Unmarshal(kv.Value, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func userInit() error {
	_, err := UserGet()
	if err == nil {
		return nil
	}
	// 初始化系统，默认账户和密码
	err = UserUpdate("admin", "admin")
	if err != nil {
		return err
	}
	return nil
}
