package agent

import (
	"gomt/core/proto/priv"
	"time"
)

func (s *DasDeviceAgent) IsConnectState() bool {
	if s.info.ConnectState > 0 && s.info.ConnectState < 6 {
		return true
	}
	return false
}

func (s *DasDeviceAgent) IsServiceAvailable(force bool) bool {
	if !s.IsConnectState() {
		return false
	}
	available, checkTime := s.getServiceAvailable()
	now := time.Now()
	diff := now.Sub(checkTime)
	if force || diff < 0 || diff > 15*time.Second {
		s.doCheckServiceAvailable()
	} else if !available {
		go s.doCheckServiceAvailable()
	}
	return available
}

func (s *DasDeviceAgent) setServiceAvailable(available bool) {
	s.serviceAvailableLock.Lock()
	defer s.serviceAvailableLock.Unlock()
	s.serviceAvailable = available
	s.serviceAvailableCheckTime = time.Now()
}

func (s *DasDeviceAgent) getServiceAvailable() (bool, time.Time) {
	s.serviceAvailableLock.RLock()
	defer s.serviceAvailableLock.RUnlock()
	available := s.serviceAvailable
	checkTime := s.serviceAvailableCheckTime
	return available, checkTime
}

func (s *DasDeviceAgent) doCheckServiceAvailable() {
	if ok := s.doCheckConnection(); !ok {
		s.setServiceAvailable(false)
		return
	}
	if ok := s.doCheckGetParameters(); !ok {
		s.setServiceAvailable(false)
		return
	}
	s.setServiceAvailable(true)
	return
}

func (s *DasDeviceAgent) doCheckConnection() bool {
	if s.info.ConnectState > 0 && s.info.ConnectState < 6 {
		return true
	}
	return false
}

func (s *DasDeviceAgent) doCheckGetParameters() bool {
	inputValues := []priv.Parameter{
		{ID: "T02.P0102"},
	}
	values, err := s.doGetParameterValues(inputValues)
	if err == nil && len(values) == 1 {
		return true
	}
	return false
}
