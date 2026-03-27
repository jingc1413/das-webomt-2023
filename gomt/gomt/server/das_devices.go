package server

import (
	"fmt"
	"gomt/core/model"
	"gomt/das/agent"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type QueryDevicesProgress struct {
	Loading  bool
	Finished bool

	Total   uint8
	Index   uint8
	Success bool
	Message string
}

func (s *OMTServer) updateDasDevices(force bool) {
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

	s.log.Tracef("query devices started")
	defer s.log.Tracef("query devices finished")

	s.queryDevicesLock.Lock()
	s.queryDevicesProgress.Loading = true
	s.queryDevicesProgress.Finished = false
	s.queryDevicesProgress.Total = 0
	s.queryDevicesProgress.Index = 0
	s.queryDevicesLock.Unlock()

	lastDevieInfos := s.dassys.GetAllDeviceInfos()

	localAgent := s.getDasDeviceAgent("local")

	queryString := "01020304060708090B0C0D0E0F"
	if s.opts.Schema == "corning" {
		queryString = "01020304060708090A0B0C0D0E0F"
	}
	deviceInfos, err := localAgent.QueryDevices(s.opts.Schema, queryString, s.queryDeviceCallback)
	if err != nil {
		s.log.Error(errors.Wrap(err, "query devices"))

		s.queryDevicesLock.Lock()
		s.queryDevicesProgress.Loading = false
		s.queryDevicesProgress.Finished = true
		s.queryDevicesProgress.Success = false
		s.queryDevicesProgress.Message = "query devices error"
		s.queryDevicesLock.Unlock()

		s.sendDasNotifyMessage("queryDeviceStatus", s.queryDevicesProgress)
		return
	}
	var wg sync.WaitGroup
	for _, v := range deviceInfos {
		info := v

		wg.Add(1)
		go func(info *model.DeviceInfo) {
			defer wg.Done()
			if _, err := s.setupOrUpdateDevice(info); err != nil {
				s.log.Error(errors.Wrapf(err, "setup device %v", info.SubID))
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
		if logrus.GetLevel() == logrus.TraceLevel {
			deviceTopo.Dump(0)
		}
	}
	s.queryDevicesTime = time.Now()
	s.queryDevicesInterval = time.Second * time.Duration(len(deviceInfos))

	s.queryDevicesLock.Lock()
	s.queryDevicesProgress.Loading = false
	s.queryDevicesProgress.Finished = true
	s.queryDevicesProgress.Success = true
	s.queryDevicesProgress.Message = ""
	s.queryDevicesLock.Unlock()

	s.sendDasNotifyMessage("queryDeviceStatus", s.queryDevicesProgress)
}

func (s *OMTServer) queryDeviceCallback(total uint8, index uint8, info *model.DeviceInfo) {
	if info == nil {
		return
	}
	deviceAgent, err := s.setupOrUpdateDevice(info)
	if err != nil {
		s.log.Error(errors.Wrapf(err, "setup device %v", info.SubID))
	}
	if deviceAgent == nil {
		return
	}

	s.queryDevicesLock.Lock()
	s.queryDevicesProgress.Total = total
	s.queryDevicesProgress.Index += 1
	s.queryDevicesLock.Unlock()

	s.sendDasNotifyMessage("queryDeviceStatus", s.queryDevicesProgress)
	s.sendDasNotifyMessage("queryDeviceData", deviceAgent.GetDeviceInfo())
}

func (s *OMTServer) setupOrUpdateDevice(info *model.DeviceInfo) (*agent.DasDeviceAgent, error) {
	deviceSub := fmt.Sprintf("%v", info.SubID)
	if info.ProductModel == "" || info.DeviceTypeName == "" {
		return nil, errors.Errorf("unknow device type %v", info.DeviceTypeName)
	}

	if productDefine := model.GetProductDefine(s.opts.Schema, info.DeviceTypeName); productDefine == nil {
		s.log.Warnf("no product define for device type, deviceTypeName=%v version=%v\n", info.DeviceTypeName, info.Version)
	} else if !productDefine.MatcheVersion(info.Version) {
		s.log.Warnf("no match version for device type, deviceTypeName=%v version=%v\n", info.DeviceTypeName, info.Version)
	}

	deviceAgent := s.getDasDeviceAgent(deviceSub)
	if deviceAgent != nil {
		info2 := deviceAgent.GetDeviceInfo()
		if info2.SubID != info.SubID ||
			info2.ConnectState != info.ConnectState ||
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
