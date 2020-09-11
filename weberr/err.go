package weberr

import (
	"encoding/json"
	"errors"
	"github.com/lixiangyun/https-gateway/util"
	"os"
)

// 系统内部错误码
var (
	ERR_EXIST = errors.New("source exist")
	ERR_NO_EXIST = errors.New("source not exist")
	ERR_DB = errors.New("db connect fail")
	ERR_DATA = errors.New("data invaild")
	ERR_PARM = errors.New("param invaild")
	ERR_STATUS_EN = errors.New("status is enable")
	ERR_STATUS_DIS = errors.New("status is disable")
)

// 请求响应错误信息
type PublicRsponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Extend  string `json:"ext"`
}

// 预定义web错误码
var (
	WEB_ERR_OK            = PublicRsponse{Code: 0, Message: "success"}

	WEB_ERR_JSON_CODER    = PublicRsponse{Code: 1001, Message: "json marshal fail"}
	WEB_ERR_JSON_UCODER   = PublicRsponse{Code: 1002, Message: "json unmarshal fail"}

	WEB_ERR_PARAM         = PublicRsponse{Code: 1004, Message: "input param invalid"}
	WEB_ERR_HTTPS_PORT    = PublicRsponse{Code: 1005, Message: "port has beed used"}
	WEB_ERR_BACKED_FAIL   = PublicRsponse{Code: 1006, Message: "backend http url invalid"}

	WEB_ERR_NOT_CERT      = PublicRsponse{Code: 1010, Message: "domain cert not exist"}
	WEB_ERR_HAS_CERT      = PublicRsponse{Code: 1011, Message: "domain cert has been exist"}
	WEB_ERR_ADD_CERT      = PublicRsponse{Code: 1012, Message: "add proxy fail"}
	WEB_ERR_DEL_CERT      = PublicRsponse{Code: 1013, Message: "delete proxy fail"}
	WEB_ERR_UP_CERT       = PublicRsponse{Code: 1014, Message: "update proxy fail"}
	WEB_ERR_USED_CERT     = PublicRsponse{Code: 1015, Message: "cert using"}

	WEB_ERR_NOT_PROXY     = PublicRsponse{Code: 1020, Message: "proxy no found"}
	WEB_ERR_ADD_PROXY     = PublicRsponse{Code: 1021, Message: "add proxy fail"}
	WEB_ERR_DEL_PROXY     = PublicRsponse{Code: 1022, Message: "delete proxy fail"}
	WEB_ERR_UP_PROXY      = PublicRsponse{Code: 1023, Message: "update proxy fail"}
	WEB_ERR_NG_PROXY      = PublicRsponse{Code: 1024, Message: "sync proxy to nginx fail"}
)

func WebOK() string {
	return WebErr(WEB_ERR_OK)
}

func WebErr(weberr PublicRsponse,ext... string) string {
	if len(ext) > 0 {
		weberr.Extend = util.StringList(ext)
	}
	errs, _ := json.Marshal(&weberr)
	return string(errs)
}

func WebErrMake(weberr PublicRsponse,ext... string) PublicRsponse {
	if len(ext) > 0 {
		weberr.Extend = util.StringList(ext)
	}
	return weberr
}

func consoleLoad(errs PublicRsponse)  {
	tempList = append(tempList, errs)
}

var tempList []PublicRsponse

func init()  {
	consoleLoad(WEB_ERR_OK)

	errs, err := json.Marshal(tempList)
	if err != nil {
		panic(err.Error())
	}

	_, err = os.Stat("./static")
	if err != nil {
		return
	}

	err = util.SaveToFile("./static/errors.json", errs)
	if err != nil {
		panic(err.Error())
	}
}
