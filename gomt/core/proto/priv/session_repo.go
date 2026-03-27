package priv

import (
	"fmt"
	"sync"
)

type deviceSessionRepository struct {
	sync.Mutex
	repo map[string]*DeviceSession
}

var defaultDeviceSessionRepository *deviceSessionRepository
var defaultDeviceSessionRepositoryOnce sync.Once

func GetDefaultDeviceSessionRepository() *deviceSessionRepository {
	defaultDeviceSessionRepositoryOnce.Do(func() {
		defaultDeviceSessionRepository = NewDeviceSessionRepository()
	})
	return defaultDeviceSessionRepository
}

func NewDeviceSessionRepository() *deviceSessionRepository {
	s := &deviceSessionRepository{}
	s.repo = map[string]*DeviceSession{}
	return s
}

func (s *deviceSessionRepository) RegiterDeviceSession(
	proto string,
	apType ApType,
	vpType VpType,
	mcpType McpType,
	deviceTypeName string,
	mainId []byte,
	subId uint8,
) *DeviceSession {
	s.Lock()
	defer s.Unlock()
	sess := NewDeviceSession(proto, apType, vpType, mcpType, deviceTypeName, mainId, subId)
	key := fmt.Sprintf("%v", sess.SubID)
	s.repo[key] = sess
	return sess
}

func (s *deviceSessionRepository) GetDeviceSession(
	subId uint8,
) *DeviceSession {
	s.Lock()
	defer s.Unlock()
	key := fmt.Sprintf("%v", subId)
	sess, ok := s.repo[key]
	if !ok || sess == nil {
		return nil
	}
	return sess
}
