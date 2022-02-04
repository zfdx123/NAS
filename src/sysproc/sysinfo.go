package sysproc

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"time"
)

type (
	Info struct {
		Cpu  interface{}
		Mem  interface{}
		Disk interface{}
	}

	cp struct {
		CpuInfo    interface{}
		CpuPercent float64 `json:"CpuPercent"`
	}

	men struct {
		Percent float64 `json:"Percent"`
		Free    uint64  `json:"Free"`
		Total   uint64  `json:"Total"`
	}

	dis struct {
		DiskPercent float64 `json:"DiskPercent"`
		DiskTotal   uint64  `json:"DiskTotal"`
		DiskFree    uint64  `json:"DiskFree"`
		DiskFs      string  `json:"DiskFs"`
		DiskPath    string  `json:"DiskPath"`
	}
)

func GetSysInfo(c *gin.Context) {
	c.JSON(http.StatusOK, Info{
		Cpu:  getCpuPercent(),
		Mem:  getMemPercent(),
		Disk: getDiskPercent(),
	})
}

func GetWs() Info {
	var inf Info
	inf = Info{
		Cpu:  getCpuPercent(),
		Mem:  getMemPercent(),
		Disk: getDiskPercent(),
	}
	return inf
}

func getCpuPercent() cp {
	var cpuPercent cp
	percent, _ := cpu.Percent(time.Second, false)
	stats, _ := cpu.Info()
	cpuPercent = cp{
		CpuInfo:    stats,
		CpuPercent: percent[0],
	}
	return cpuPercent
}

func getMemPercent() men {
	var memPercent men
	memInfo, _ := mem.VirtualMemory()
	memPercent = men{
		Percent: memInfo.UsedPercent,
		Free:    memInfo.Free,
		Total:   memInfo.Total,
	}
	return memPercent
}

func getDiskPercent() dis {
	var diskPercent dis
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	diskPercent = dis{
		DiskPercent: diskInfo.UsedPercent,
		DiskTotal:   diskInfo.Total,
		DiskFree:    diskInfo.Free,
		DiskFs:      diskInfo.Fstype,
		DiskPath:    diskInfo.Path,
	}
	return diskPercent
}
