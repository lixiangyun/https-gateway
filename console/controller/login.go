package controller

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/weberr"
	"net/http"
	"time"
)

type ChangePasswdRequest struct {
	OldUser  string `json:"oldusername"`
	OldPwd   string `json:"oldpassword"`
	User     string `json:"username"`
	Pwd      string `json:"password"`
}

func ChangePwdControllerPut(ctx *context.Context)  {
	var req ChangePasswdRequest

	werr := weberr.WEB_ERR_OK
	defer func() {
		ctx.WriteString(weberr.WebErr(werr))
		if werr.Code != 0 {
			time.Sleep(5 * time.Second)
		}
	}()

	body := ctx.Input.RequestBody
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Error("json unmarshal fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_JSON_UCODER)
		return
	}

	userinfo, err := data.UserGet()
	if err != nil {
		logs.Error("get login user fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_NO_USER)
		return
	}

	if userinfo.Name != req.OldUser {
		logs.Error("check old username fail", req.User)
		werr = weberr.WebErrMake(weberr.WEB_ERR_LOGIN_FAIL)
		return
	}

	if userinfo.CheckPasswod(req.OldPwd) == false {
		logs.Error("check old password fail", req.User)
		werr = weberr.WebErrMake(weberr.WEB_ERR_LOGIN_FAIL)
		return
	}

	err = data.UserUpdate(req.User, req.Pwd)
	if err != nil {
		logs.Error("update user and password fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_USER_UP_FAIL)
		return
	}

	LoginSessionInit(ctx, req.User)
}

type loginRequest struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

func LoginControllerPost(ctx *context.Context)  {
	var req loginRequest

	werr := weberr.WEB_ERR_OK
	defer func() {
		ctx.WriteString(weberr.WebErr(werr))
		if werr.Code != 0 {
			time.Sleep(5 * time.Second)
		}
	}()

	body := ctx.Input.RequestBody
	err := json.Unmarshal(body, &req)
	if err != nil {
		logs.Error("json unmarshal fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_JSON_UCODER)
		return
	}

	userinfo, err := data.UserGet()
	if err != nil {
		logs.Error("get login user fail", err.Error())
		werr = weberr.WebErrMake(weberr.WEB_ERR_NO_USER)
		return
	}

	if userinfo.Name != req.User {
		logs.Error("check username fail", req.User)
		werr = weberr.WebErrMake(weberr.WEB_ERR_LOGIN_FAIL)
		return
	}

	if userinfo.CheckPasswod(req.Password) == false {
		logs.Error("check password fail", req.User)
		werr = weberr.WebErrMake(weberr.WEB_ERR_LOGIN_FAIL)
		return
	}

	LoginSessionInit(ctx, req.User)
}

func LoginSessionInit(ctx *context.Context, user string) {
	sess, _ := globalSessions.SessionStart(ctx.ResponseWriter,ctx.Request)
	defer sess.SessionRelease(ctx.ResponseWriter)
	sess.Set("username",user)
}

func LoginSessionGet(ctx *context.Context) string {
	sess, _ := globalSessions.SessionStart(ctx.ResponseWriter,ctx.Request)
	defer sess.SessionRelease(ctx.ResponseWriter)
	infe := sess.Get("username")
	if infe == nil {
		return ""
	}
	value , ok:=infe.(string)
	if ok {
		return value
	}
	return ""
}

func NoLogin(ctx *context.Context) {
	ctx.WriteString(weberr.WebErr(weberr.WEB_ERR_NO_LOGIN))
	ctx.Output.Status = http.StatusUnauthorized
}

func LogoutControllerGet(ctx *context.Context)  {
	user := LoginSessionGet(ctx)
	if user != "" {
		LoginSessionInit(ctx, "")
		logs.Info("user logout success", user)
	}
	ctx.Redirect(302, "/login.html")
}

var globalSessions * session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:"HttpsGatewaySessionID",
		EnableSetCookie: true,
		Gclifetime: 3600,
		Secure: true,
	}
	globalSessions, _ = session.NewManager("memory",sessionConfig)
	go globalSessions.GC()
}
