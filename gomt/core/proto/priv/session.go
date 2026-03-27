package priv

import (
	"gomt/core/model"
)

type DeviceSession struct {
	Proto   string
	ApType  ApType
	VpType  VpType
	McpType McpType

	DeviceTypeName string
	MainID         []byte
	SubID          uint8
	RouteAddress   []byte
	Addr           string

	Params model.ParameterDefines
}

func NewDeviceSession(
	proto string,
	apType ApType,
	vpType VpType,
	mcpType McpType,
	deviceTypeName string,
	mainId []byte,
	subId uint8,
) *DeviceSession {
	sess := &DeviceSession{
		Proto:   proto,
		ApType:  apType,
		VpType:  vpType,
		McpType: mcpType,
		MainID:  mainId,
		SubID:   subId,
	}
	return sess
}

func (m *DeviceSession) SetDeviceID(mainID []byte, subID uint8) {
	m.MainID = mainID
	m.SubID = subID
}
