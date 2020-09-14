package controller

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/lixiangyun/https-gateway/console/data"
	"github.com/lixiangyun/https-gateway/console/nginx"
	"github.com/lixiangyun/https-gateway/os"
	"github.com/lixiangyun/https-gateway/util"
	"time"
)

type Console struct {
	TodayCnt  int    `json:"today_request"`
	TotalCnt  int    `json:"total_request"`
	TodaySize string `json:"today_flow"`
	TotalSize string `json:"total_flow"`
	Certs     int    `json:"cert_number"`
	Proxys    int    `json:"proxy_number"`
	Cpu       int    `json:"node_cpu"`
	Memory    int    `json:"node_memory"`
	
	Version   string `json:"version"`
}

var sysinfo *os.SysInfo

var FlowTodaySize int64
var FlowTotalSize int64

var RequestTodayCnt  int
var RequestTotalCnt  int

func GetNetwork(network []os.Network) (int64,int64) {
	var down, up int64
	for _,v := range network {
		down += v.Downflow
		up += v.Upflow
	}
	return up, down
}

var lastCnt int
func GetRequest() int {
	totalcnt := nginx.AccessCntGet()
	var retcnt int
	if lastCnt > totalcnt {
		retcnt = totalcnt
	} else {
		retcnt = totalcnt - lastCnt
	}
	lastCnt = totalcnt
	return retcnt
}

func init()  {
	sysinfo = os.SysInfoGet()
	go func() {
		for  {
			before := time.Now()
			time.Sleep(time.Minute)
			after := time.Now()

			if after.Day() != before.Day() {
				FlowTodaySize = 0
				RequestTodayCnt = 0
			}

			sysinfo = os.SysInfoGet()

			up, down := GetNetwork(sysinfo.Net)

			sys := data.SysStatGet()
			sys.DownFlowSize += down
			sys.UpFlowSize += up

			FlowTotalSize = sys.DownFlowSize + sys.UpFlowSize
			FlowTodaySize += up + down

			cnt := GetRequest()
			sys.ReqeustCnt += cnt

			RequestTodayCnt += cnt
			RequestTotalCnt = sys.ReqeustCnt

			data.SysStatUpdate(sys)
		}
	}()
}

func ConsoleInfoControllerGet(ctx *context.Context)  {
	proxys, _ := data.ProxyQueryAll()
	certs, _ := data.CertQueryAll()

	var rsp Console

	rsp.Proxys = len(proxys)
	rsp.Certs = len(certs)
	rsp.Cpu = sysinfo.Cpu.Usage/100
	rsp.Memory = sysinfo.Mem.Usage/100
	rsp.TotalSize = util.ByteView(FlowTotalSize)
	rsp.TodaySize = util.ByteView(FlowTodaySize)
	rsp.TotalCnt = RequestTotalCnt
	rsp.TodayCnt = RequestTodayCnt
	rsp.Version = util.VersionGet()

	result, _ := json.Marshal(&rsp)
	ctx.WriteString(string(result))
}

