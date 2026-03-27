package priv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gomt/core/model"
	"strings"

	"github.com/spf13/cast"
)

const (
	QUERY_DEVICE_INPUT_SUB_ID            uint8 = 0x01
	QUERY_DEVICE_INPUT_ROUTE             uint8 = 0x02
	QUERY_DEVICE_INPUT_IPADDR            uint8 = 0x03
	QUERY_DEVICE_INPUT_CONNECT_STATE     uint8 = 0x04
	QUERY_DEVICE_INPUT_SUB_CODE          uint8 = 0x05
	QUERY_DEVICE_INPUT_DEVICE_TYPE_ID    uint8 = 0x06
	QUERY_DEVICE_INPUT_LOCATION_INFO     uint8 = 0x07
	QUERY_DEVICE_INPUT_ALARM_STATE       uint8 = 0x08 // WEBOMT only
	QUERY_DEVICE_INPUT_VERSION_STATE     uint8 = 0x09 // WEBOMT only
	QUERY_DEVICE_INPUT_MIXED_STATE       uint8 = 0x0A // WEBOMT only
	QUERY_DEVICE_INPUT_OPTICAL_STATE     uint8 = 0x0B // WEBOMT only
	QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME  uint8 = 0x0C
	QUERY_DEVICE_INPUT_DEVICE_NAME       uint8 = 0x0D
	QUERY_DEVICE_INPUT_VERSION           uint8 = 0x0E
	QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM uint8 = 0x0F
)

type QueryDeviceInputDefine struct {
	ID   uint8
	Size uint8
}

var QueryDeviceInputDefineMap map[uint8]QueryDeviceInputDefine = map[uint8]QueryDeviceInputDefine{
	QUERY_DEVICE_INPUT_SUB_ID:            {ID: QUERY_DEVICE_INPUT_SUB_ID, Size: 1},
	QUERY_DEVICE_INPUT_ROUTE:             {ID: QUERY_DEVICE_INPUT_ROUTE, Size: 4},
	QUERY_DEVICE_INPUT_IPADDR:            {ID: QUERY_DEVICE_INPUT_IPADDR, Size: 4},
	QUERY_DEVICE_INPUT_CONNECT_STATE:     {ID: QUERY_DEVICE_INPUT_CONNECT_STATE, Size: 1},
	QUERY_DEVICE_INPUT_SUB_CODE:          {ID: QUERY_DEVICE_INPUT_SUB_CODE, Size: 10},
	QUERY_DEVICE_INPUT_DEVICE_TYPE_ID:    {ID: QUERY_DEVICE_INPUT_DEVICE_TYPE_ID, Size: 1},
	QUERY_DEVICE_INPUT_LOCATION_INFO:     {ID: QUERY_DEVICE_INPUT_LOCATION_INFO, Size: 40},
	QUERY_DEVICE_INPUT_ALARM_STATE:       {ID: QUERY_DEVICE_INPUT_ALARM_STATE, Size: 1},
	QUERY_DEVICE_INPUT_VERSION_STATE:     {ID: QUERY_DEVICE_INPUT_VERSION_STATE, Size: 1},
	QUERY_DEVICE_INPUT_MIXED_STATE:       {ID: QUERY_DEVICE_INPUT_MIXED_STATE, Size: 1},
	QUERY_DEVICE_INPUT_OPTICAL_STATE:     {ID: QUERY_DEVICE_INPUT_OPTICAL_STATE, Size: 1},
	QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME:  {ID: QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME, Size: 20},
	QUERY_DEVICE_INPUT_DEVICE_NAME:       {ID: QUERY_DEVICE_INPUT_DEVICE_NAME, Size: 40},
	QUERY_DEVICE_INPUT_VERSION:           {ID: QUERY_DEVICE_INPUT_VERSION, Size: 20},
	QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM: {ID: QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM, Size: 20},
}

type QueryDevcieOutput struct {
	SubID              uint8
	RouteAddress       []byte
	IpAddress          []byte
	ConnectState       int8
	AlarmState         int8
	OpticalState       int8
	MixedState         int8
	VersionState       int8
	SubCode            string
	DeviceTypeID       uint8
	DeviceTypeName     string
	DeviceName         string
	Version            string
	LocationInfo       string
	ElementModelNumber string
}

func (m QueryDevcieOutput) String() string {
	return fmt.Sprintf(
		`{SubID=%d, RouteAddress=%v, IpAddr=%v, DeviceTypeName=%s, DeviceTypeID=%d, LocationInfo=%s, ConnStatus=%v. ModelNum=%v}`,
		m.SubID, m.RouteAddress, m.IpAddress, m.DeviceTypeName, m.DeviceTypeID, m.LocationInfo, m.ConnectState, m.ElementModelNumber,
	)
}

func (m QueryDevcieOutput) DeviceInfo(schema string, mainID []byte) model.DeviceInfo {
	info := model.DeviceInfo{
		SubID:              m.SubID,
		DeviceTypeName:     m.DeviceTypeName,
		DeviceTypeID:       int(m.DeviceTypeID),
		RouteAddress:       m.RouteAddress,
		IpAddress:          m.IpAddress,
		ConnectState:       m.ConnectState,
		AlarmState:         m.AlarmState,
		OpticalState:       m.OpticalState,
		VersionState:       m.VersionState,
		MixedState:         m.MixedState,
		Version:            m.Version,
		InstalledLocation:  m.LocationInfo,
		DeviceName:         m.DeviceName,
		ElementModelNumber: m.ElementModelNumber,
	}
	info.Setup()
	return info
}

func (m QueryDevcieOutput) MarshalBinaryWithInputs(inputs []uint8) ([]byte, error) {
	w := new(bytes.Buffer)
	for _, input := range inputs {
		def := QueryDeviceInputDefineMap[input]
		switch input {
		case QUERY_DEVICE_INPUT_SUB_ID:
			if err := binary.Write(w, binary.LittleEndian, m.SubID); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_ROUTE:
			if err := binary.Write(w, binary.LittleEndian, m.RouteAddress); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_IPADDR:
			if err := binary.Write(w, binary.LittleEndian, m.IpAddress); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_CONNECT_STATE:
			if err := binary.Write(w, binary.LittleEndian, m.ConnectState); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_SUB_CODE:
			tmp := []byte(m.SubCode)
			b := make([]byte, def.Size)
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_DEVICE_TYPE_ID:
			if err := binary.Write(w, binary.LittleEndian, m.DeviceTypeID); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_LOCATION_INFO:
			tmp := []byte(m.LocationInfo)
			b := make([]byte, def.Size)
			for i := 0; i < int(def.Size); i++ {
				b[i] = 0
			}
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_ALARM_STATE:
			if err := binary.Write(w, binary.LittleEndian, m.AlarmState); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_VERSION_STATE:
			if err := binary.Write(w, binary.LittleEndian, m.VersionState); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_MIXED_STATE:
			if err := binary.Write(w, binary.LittleEndian, m.MixedState); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_OPTICAL_STATE:
			if err := binary.Write(w, binary.LittleEndian, m.OpticalState); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME:
			tmp := []byte(m.DeviceTypeName)
			b := make([]byte, def.Size)
			for i := 0; i < int(def.Size); i++ {
				b[i] = 0
			}
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_DEVICE_NAME:
			tmp := []byte(m.DeviceName)
			b := make([]byte, def.Size)
			for i := 0; i < int(def.Size); i++ {
				b[i] = 0
			}
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_VERSION:
			tmp := []byte(m.Version)
			b := make([]byte, def.Size)
			for i := 0; i < int(def.Size); i++ {
				b[i] = 0
			}
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		case QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM:
			tmp := []byte(m.ElementModelNumber)
			b := make([]byte, def.Size)
			for i := 0; i < int(def.Size); i++ {
				b[i] = 0
			}
			end := len(tmp)
			if end > int(def.Size) {
				end = int(def.Size)
			}
			copy(b[0:end], tmp[0:end])
			if err := binary.Write(w, binary.LittleEndian, b); err != nil {
				return w.Bytes(), err
			}
		}
	}
	return w.Bytes(), nil
}

func (m *QueryDevcieOutput) UnmarshalBinaryWithInputs(inputs []uint8, in []byte) error {
	offset := 0
	for _, input := range inputs {
		def := QueryDeviceInputDefineMap[input]
		switch input {
		case QUERY_DEVICE_INPUT_SUB_ID:
			r := bytes.NewReader(in[offset : offset+1])
			if err := binary.Read(r, binary.LittleEndian, &m.SubID); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_ROUTE:
			m.RouteAddress = make([]byte, def.Size)
			copy(m.RouteAddress, in[offset:offset+int(def.Size)])
		case QUERY_DEVICE_INPUT_IPADDR:
			m.IpAddress = make([]byte, def.Size)
			copy(m.IpAddress, in[offset:offset+int(def.Size)])
		case QUERY_DEVICE_INPUT_CONNECT_STATE:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.ConnectState); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_ALARM_STATE:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.AlarmState); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_VERSION_STATE:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.VersionState); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_MIXED_STATE:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.MixedState); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_OPTICAL_STATE:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.OpticalState); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_SUB_CODE:
			m.SubCode = string(in[offset : offset+int(def.Size)])
		case QUERY_DEVICE_INPUT_DEVICE_TYPE_ID:
			r := bytes.NewReader(in[offset : offset+int(def.Size)])
			if err := binary.Read(r, binary.LittleEndian, &m.DeviceTypeID); err != nil {
				return err
			}
		case QUERY_DEVICE_INPUT_LOCATION_INFO:
			m.LocationInfo = strings.Trim(cast.ToString(in[offset:offset+int(def.Size)]), "\x00")
		case QUERY_DEVICE_INPUT_DEVICE_TYPE_NAME:
			m.DeviceTypeName = strings.Trim(cast.ToString(in[offset:offset+int(def.Size)]), "\x00")
		case QUERY_DEVICE_INPUT_VERSION:
			m.Version = strings.Trim(cast.ToString(in[offset:offset+int(def.Size)]), "\x00")
		case QUERY_DEVICE_INPUT_DEVICE_NAME:
			m.DeviceName = strings.Trim(cast.ToString(in[offset:offset+int(def.Size)]), "\x00")
		case QUERY_DEVICE_INPUT_ELEMENT_MODEL_NUM:
			m.ElementModelNumber = strings.Trim(cast.ToString(in[offset:offset+int(def.Size)]), "\x00")
		}
		offset += int(def.Size)
	}
	return nil
}
