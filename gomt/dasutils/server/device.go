package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gomt/core/model"
	"gomt/das"
	"gomt/das/agent"

	"github.com/pkg/errors"
)

type QueryDevicesProgress struct {
	Loading  bool
	Finished bool

	Total   uint8
	Index   uint8
	Success bool
	Message string
}

func (s *DasUtilsServer) updateDasDevices(force bool) {
	s.queryDevicesLock.RLock()
	if s.queryDevicesRunning {
		s.queryDevicesLock.RUnlock()
		return
	}
	s.queryDevicesLock.RUnlock()

	s.queryDevicesLock.Lock()
	s.queryDevicesRunning = true
	s.queryDevicesLock.Unlock()

	defer func() {
		s.queryDevicesLock.Lock()
		s.queryDevicesRunning = false
		s.queryDevicesLock.Unlock()
	}()

	if s.queryDevicesInterval < time.Second*30 {
		s.queryDevicesInterval = time.Second * 30
	}

	t := time.Now()
	if !force && t.After(s.queryDevicesTime) && t.Before(s.queryDevicesTime.Add(s.queryDevicesInterval)) {
		return
	}

	log.Println("query devices started")
	defer log.Println("query devices finished")

	lastDevieInfos := s.dassys.GetAllDeviceInfos()

	localAgent := s.getDasDeviceAgent("local")
	queryString := "01020304060708090B0C0D0E0F"
	deviceInfos, err := localAgent.QueryDevices(s.opts.Schema, queryString, s.queryDeviceCallback)
	if err != nil {
		//s.log.Error(errors.Wrap(err, "query devices"))
		log.Println("")
		return
	}
	var wg sync.WaitGroup
	for _, v := range deviceInfos {
		info := v

		wg.Add(1)
		go func(info *model.DeviceInfo) {
			defer wg.Done()
			if _, err := s.setupOrUpdateDevice(info); err != nil {
				log.Printf("setup device, %d, %v\n", info.SubID, err)
			}
		}(info)
	}
	wg.Wait()

	for _, last := range lastDevieInfos {
		found := false
		for _, info := range deviceInfos {
			if info.SubID == last.SubID && info.DeviceTypeName == last.DeviceTypeName {
				found = true
				break
			}
		}
		if !found {
			s.dassys.DeleteAgent(fmt.Sprintf("%v", last.SubID))
		}
	}

	if deviceTopo := model.GetDeviceTopo(deviceInfos); deviceTopo != nil {
		deviceTopo.Dump(0)
	}
	s.queryDevicesTime = time.Now()
	s.queryDevicesInterval = time.Second * time.Duration(len(deviceInfos))
}

func (s *DasUtilsServer) getDasDeviceAgent(deviceSub string) *agent.DasDeviceAgent {
	return s.dassys.GetAgent(deviceSub)
}

func (s *DasUtilsServer) queryDeviceCallback(total uint8, index uint8, info *model.DeviceInfo) {
	if info == nil {
		return
	}
	deviceAgent, err := s.setupOrUpdateDevice(info)
	if err != nil {
		log.Printf("setup device call, %d, %v\n", info.SubID, err)
	}
	if deviceAgent == nil {
		return
	}
}

func (s *DasUtilsServer) setupDasDeviceAgent(info *model.DeviceInfo) (*agent.DasDeviceAgent, error) {
	if info == nil {
		return nil, errors.New("invalid device info")
	}
	opts := agent.DasDeviceAgentOptions{
		PrivServerPort: s.opts.DevicePort,
		CGIBinUrlPath:  "/cgi-bin",
		CGIBinFilePath: s.opts.CgiBinDir,
		FileTypes:      s.opts.FileTypes,
		ConfigPath:     s.opts.ConfigDir,
	}
	if s.localhostMode {
		opts.DeviceAddr = info.IpAddressString
		opts.HttpServerAddr = das.MakeDeviceUrlBase(info.IpAddressString, 443, true)
	} else {
		opts.DeviceAddr = s.opts.DeviceAddr
		opts.HttpServerAddr = das.MakeDeviceUrlBase(s.opts.DeviceAddr, info.ForwardingPort, true)
	}
	deviceSub := fmt.Sprintf("%v", info.SubID)
	a, err := s.dassys.SetupDasDeviceAgent(s.opts.Schema, deviceSub, info, opts)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *DasUtilsServer) setupOrUpdateDevice(info *model.DeviceInfo) (*agent.DasDeviceAgent, error) {
	deviceSub := fmt.Sprintf("%v", info.SubID)
	if info.ProductModel == "" || info.DeviceTypeName == "" {
		return nil, errors.Errorf("unknow device type %v", info.DeviceTypeName)
	}
	if productDefine := model.GetProductDefine(s.opts.Schema, info.DeviceTypeName); productDefine == nil {
		log.Printf("no product define for device type, deviceTypeName=%v version=%v\n", info.DeviceTypeName, info.Version)
	} else if !productDefine.MatcheVersion(info.Version) {
		log.Printf("no match version for device type, deviceTypeName=%v version=%v\n", info.DeviceTypeName, info.Version)
	}

	deviceAgent := s.getDasDeviceAgent(deviceSub)
	if deviceAgent != nil {
		info2 := deviceAgent.GetDeviceInfo()
		if info2.SubID != info.SubID ||
			info2.DeviceTypeName != info.DeviceTypeName ||
			info2.IpAddressString != info.IpAddressString ||
			info2.RouteAddressString != info.RouteAddressString {
			s.dassys.DeleteAgent(deviceSub)
			deviceAgent = nil
		}
	}
	if deviceAgent == nil {
		a, err := s.setupDasDeviceAgent(info)
		if err != nil {
			return nil, errors.Wrapf(err, "setup device agent")
		}
		deviceAgent = a
	}

	needUpdate := false
	if ok := deviceAgent.IsServiceAvailable(true); ok {
		needUpdate = true
	}
	if err := deviceAgent.UpdateDeviceInfoByQuery(info, needUpdate); err != nil {
		return deviceAgent, errors.Wrapf(err, "update device info")
	}
	return deviceAgent, nil
}
