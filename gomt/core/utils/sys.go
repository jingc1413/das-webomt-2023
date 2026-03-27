package utils

import (
	"time"

	"github.com/jaypipes/ghw"
	"github.com/pkg/errors"
	gocpu "github.com/shirou/gopsutil/v3/cpu"
	godisk "github.com/shirou/gopsutil/v3/disk"
	gomem "github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

func init() {
	// GetSystemInfo()
	GetSystemStatus(1 * time.Second)
}

type Product struct {
	Family       string
	Name         string
	SerialNumber string
	UUID         string
	Vendor       string
}

type ProcessorCore struct {
	ID         int
	NumThreads uint32
}

type Processor struct {
	ID         int
	Model      string
	Vendor     string
	NumThreads uint32
	NumCores   uint32
	Cores      []*ProcessorCore
}

type CPUInfo struct {
	ID           uint32
	TotalThreads uint32
	TotalCores   uint32
	Processors   []*Processor
}

type NIC struct {
	Name       string
	MacAddress string
}

type SystemInfo struct {
	Product  *Product
	CPU      *CPUInfo
	Networks []*NIC
}

type CPUStat struct {
	Percent float64
}

type MemoryStat struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

type PartitionStat struct {
	Device      string
	FsType      string
	MountPoint  string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type MountPointStat struct {
	MountPoint  string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type SystemStatus struct {
	CPU        CPUStat
	Memory     MemoryStat
	Partitions []PartitionStat
}

type DirectoryDiskStatus struct {
	MountPoints []MountPointStat
}

func GetSystemStatus(interval time.Duration) (*SystemStatus, error) {
	status := &SystemStatus{}

	if v, err := gocpu.Percent(interval, false); err == nil {
		if len(v) > 0 {
			status.CPU.Percent = v[0]
		}
	} else {
		logrus.Error(errors.Wrap(err, "get cpu usage percent"))
	}
	if v, err := gomem.VirtualMemory(); v != nil {
		status.Memory.Total = v.Total
		status.Memory.Available = v.Available
		status.Memory.UsedPercent = v.UsedPercent
		status.Memory.Used = v.Used
	} else {
		logrus.Error(errors.Wrap(err, "get memory info"))
	}

	status.Partitions = []PartitionStat{}
	if v, err := godisk.Partitions(false); err == nil {
		for _, v2 := range v {
			stat := PartitionStat{
				Device:     v2.Device,
				MountPoint: v2.Mountpoint,
				FsType:     v2.Fstype,
			}
			if v3, err := godisk.Usage(stat.MountPoint); err == nil {
				stat.Total = v3.Total
				stat.Free = v3.Free
				stat.Used = v3.Used
				stat.UsedPercent = v3.UsedPercent
			} else {
				logrus.Error(errors.Wrap(err, "get partition usage"))
			}
			status.Partitions = append(status.Partitions, stat)
		}
	} else {
		logrus.Error(errors.Wrap(err, "get partitions info"))
	}
	return status, nil
}

func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}
	product, err := ghw.Product()
	if err != nil {
		return nil, errors.Wrap(err, "get product")
	}
	info.Product = &Product{
		Family:       product.Family,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
		UUID:         product.UUID,
		Vendor:       product.Vendor,
	}

	cpu, err := ghw.CPU()
	if err != nil {
		return nil, errors.Wrap(err, "get cpu info")
	}
	info.CPU = &CPUInfo{
		TotalCores:   cpu.TotalCores,
		TotalThreads: cpu.TotalThreads,
		Processors:   []*Processor{},
	}
	for _, processor := range cpu.Processors {
		p := &Processor{
			ID:         processor.ID,
			Model:      processor.Model,
			Vendor:     processor.Vendor,
			NumCores:   processor.NumCores,
			NumThreads: processor.NumThreads,
			//Cores:      []*ProcessorCore{},
		}
		//for _, core := range processor.Cores {
		//c := &ProcessorCore{
		//ID:         core.ID,
		//NumThreads: core.NumThreads,
		//}
		//p.Cores = append(p.Cores, c)
		//}
		info.CPU.Processors = append(info.CPU.Processors, p)
	}

	net, err := ghw.Network()
	if err != nil {
		return info, errors.Wrap(err, "get network info")
	}
	info.Networks = []*NIC{}
	for _, nic := range net.NICs {
		if nic.IsVirtual {
			continue
		}
		n := &NIC{
			Name:       nic.Name,
			MacAddress: nic.MacAddress,
		}
		info.Networks = append(info.Networks, n)
	}

	return info, nil
}

func GetDirectoryDiskStatus(paths []string) (*DirectoryDiskStatus, error) {
	var status DirectoryDiskStatus
	status.MountPoints = []MountPointStat{}
	for _, path := range paths {
		if ExistsDir(path) == false {
			continue
		}
		stat := MountPointStat{
			MountPoint: path,
		}
		if v3, err := godisk.Usage(path); err == nil {
			stat.Total = v3.Total
			stat.Free = v3.Free
			stat.Used = v3.Used
			stat.UsedPercent = v3.UsedPercent
		} else {
			logrus.Error(errors.Wrap(err, "get path usage"))
		}
		status.MountPoints = append(status.MountPoints, stat)
	}
	return &status, nil
}
