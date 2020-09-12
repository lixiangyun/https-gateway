package os

import (
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/os/cpu"
	diskpkg "github.com/lixiangyun/https-gateway/os/disk"
	"github.com/lixiangyun/https-gateway/os/memory"
	"github.com/lixiangyun/https-gateway/os/network"
	"github.com/lixiangyun/https-gateway/util"
	"sync"
)

type CpuInfo struct {
	Total int `json:"total"`
	Used  int `json:"used"`
	Usage int `json:"usage"`
}

type MemInfo struct {
	Total int `json:"total"`
	Used  int `json:"used"`
	Free  int `json:"free"`
	Usage int `json:"usage"`
}

type DiskInfo struct {
	Name    string `json:"name"`
	Mounted string `json:"mounted"`
	Total   int    `json:"total"`
	Used    int    `json:"used"`
	Usage   int    `json:"usage"`
}

type Network struct {
	Name     string `json:"name"`
	Downflow int64  `json:"downflow"`
	Upflow   int64  `json:"upflow"`
}

type SysInfo struct {
	sync.Mutex
	Date string     `json:"date"`
	Cpu  CpuInfo    `json:"cpu"`
	Mem  MemInfo    `json:"mem"`
	Net  []Network  `json:"network"`
	Disk []DiskInfo `json:"disk"`
}

func memInfoGet() MemInfo {
	var smem MemInfo
	mem, err := memory.Get()
	if err != nil {
		logs.Error("get memory info fail", err.Error())
		return smem
	}
	smem.Total = int(mem.Total)
	smem.Free = int(mem.Free)
	smem.Used = int(mem.Used)
	smem.Usage = util.TenThousandPercent(smem.Total, smem.Used)
	return smem
}

var lastCpuInfo *cpu.Stats

func cpuInfoGet() CpuInfo {
	var scpu CpuInfo
	cs, err := cpu.Get()
	if err != nil {
		logs.Error("get cpu info fail", err.Error())
		return scpu
	}
	if lastCpuInfo == nil {
		scpu.Total = int(cs.Total)
		scpu.Used = int(cs.Total - cs.Idle)
	} else {
		scpu.Total = int(cs.Total - lastCpuInfo.Total)
		scpu.Used = scpu.Total - int(cs.Idle - lastCpuInfo.Idle)
	}
	scpu.Usage = util.TenThousandPercent(scpu.Total, scpu.Used)
	lastCpuInfo = cs
	return scpu
}

func diskInfoGet() []DiskInfo {
	sdisk := make([]DiskInfo, 0)
	ds, err := diskpkg.GetDiskUsage()
	if err != nil {
		logs.Error("get disk info fail", err.Error())
		return sdisk
	}
	for _, v:= range ds {
		sdisk = append(sdisk, DiskInfo{
			Name: v.Name,
			Total: v.Total,
			Used: v.Used,
			Usage: util.TenThousandPercent(v.Total, v.Used),
			Mounted: v.Mount,
		})
	}
	return sdisk
}

var lastNetInfo map[string]network.Stats

func networkInfoGet() []Network {
	snet := make([]Network, 0)
	if lastNetInfo == nil {
		lastNetInfo = make(map[string]network.Stats, 0)
	}
	ns, err := network.Get()
	if err != nil {
		logs.Error("get network info fail", err.Error())
		return snet
	}
	for _, v := range ns {
		last, flag := lastNetInfo[v.Name]
		if flag == false {
			snet = append(snet, Network{
				Name: v.Name,
				Upflow: 0,
				Downflow: 0,
			})
		} else {
			snet = append(snet, Network{
				Name: v.Name,
				Upflow: int64(v.TxBytes - last.TxBytes),
				Downflow: int64(v.RxBytes - last.RxBytes),
			})
		}
		lastNetInfo[v.Name] = v
	}
	return snet
}

func SysInfoGet() *SysInfo {
	var sysall SysInfo

	sysall.Date = util.GetTimeStamp()
	sysall.Mem = memInfoGet()
	sysall.Cpu = cpuInfoGet()
	sysall.Disk = diskInfoGet()
	sysall.Net = networkInfoGet()

	return &sysall
}
