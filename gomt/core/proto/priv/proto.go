package priv

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/howeyc/crc16"
	"github.com/pkg/errors"
)

var crcTable *crc16.Table = nil

func init() {
	//crcTable = crc16.MakeTable(0x1021)
	crcTable = crc16.MakeBitsReversedTable(0x1021)
}

const (
	APB_QUERY_TIMEOUT     = 20
	APB_SET_TIMEOUT       = 40
	APC_QUERY_IDS_TIMEOUT = 14
	APC_QUERY_ONE_TIMEOUT = 8
	APC_QUERY_ALL_TIMEOUT = 20
	APC_SET_ONE_TIMEOUT   = 10
)

var ErrInvalidLength = errors.New("invalid length")
var ErrInvalidPacketFlag = errors.New("invalid packet flag")
var ErrInvalidPacketLength = errors.New("invalid packet length")
var ErrInvalidCheckSum = errors.New("invalid checksum")
var ErrInvalidApType = errors.New("invalid access protocol type")
var ErrInvalidVpType = errors.New("invalid visiting protocol type")
var ErrInvalidMcpType = errors.New("invalid monitoring control protocol type")
var ErrInvalidMcpCommandID = errors.New("invalid mcp command id")

type ApType uint8  // access protocol type
type VpType uint8  // visiting protocol type
type McpType uint8 // monitoring control protocol type

const (
	AP_A ApType = 0x01 // RS-485, RS-232, MODEM(data)
	AP_B ApType = 0x02 // MODEM SMS
	AP_C ApType = 0x03 // ETH, MODEM ADSL, MODEM PS
)

const (
	VP_A  VpType = 0x01
	VP_A1 VpType = 0xA1
)

const (
	MCP_A McpType = 0x01
	MCP_B McpType = 0x02
	MCP_C McpType = 0x03
)

const APA_PACKET_MAX_LENGTH = 256 // before escapse without flags
const APA_PACKET_FLAG_BEGIN = 0x7e
const APA_PACKET_FLAG_END = 0x7e

const APB_PACKET_MAX_LENGTH = 256     // after escapse with flag
const APB_PACKET_SMS_MAX_LENGTH = 140 // after escapse with flag
const APB_PACKET_FLAG_BEGIN = 0x21
const APB_PACKET_FLAG_END = 0x21

const APC_PACKET_MAX_LENGTH = 1500 // before escapse without flags
const APC_PACKET_FLAG_BEGIN = 0x7e
const APC_PACKET_FLAG_END = 0x7e

const VPACKET_SERIAL_OMC_BEGIN = 0x0000
const VPACKET_SERIAL_OMC_END = 0x7FFF
const VPACKET_SERIAL_DEVICE_BEGIN = 0x8000
const VPACKET_SERIAL_DEVICE_END = 0x8FFF

type VpFlag uint8

const (
	VP_FLAG_RESPONSE_OK   VpFlag = 0x00
	VP_FLAG_RESPONSE_BUSY VpFlag = 0x01
	VP_FLAG_REQUSET       VpFlag = 0x80
)

type CmdID uint8
type CmdStatus uint8

const (
	CMD_ID_NONE           CmdID = 0x00
	CMD_ID_NOTIFY         CmdID = 0x01 // not support in MCPB
	CMD_ID_QUERY          CmdID = 0x02
	CMD_ID_SET            CmdID = 0x03
	CMD_ID_UPGRADE_MODE   CmdID = 0x10 // not support in MCPB
	CMD_ID_CHANGE_VERSION CmdID = 0x11 // not support in MCPB

	CMD_ID_QUERY_B0 CmdID = 0xb0
	CMD_ID_SET_B0   CmdID = 0xb1
	CMD_ID_QUERY_B2 CmdID = 0xb2
	CMD_ID_SET_B2   CmdID = 0xb3
	CMD_ID_QUERY_B4 CmdID = 0xb4
	CMD_ID_SET_B4   CmdID = 0xb5
)

const (
	CMD_STATUS_SUCC             CmdStatus = 0x00
	CMD_STATUS_REJECT           CmdStatus = 0x01
	CMD_STATUS_INVALID_ID       CmdStatus = 0x02
	CMD_STATUS_INVALID_LENGTH   CmdStatus = 0x03
	CMD_STATUS_INVALID_CRC      CmdStatus = 0x04
	CMD_STATUS_INVALID_MCP_TYPE CmdStatus = 0x05
	CMD_STATUS_OTHER            CmdStatus = 0xFE
	CMD_STATUS_COMMAND          CmdStatus = 0xFF
)

func GetQueryCommandId(key string) CmdID {
	key = strings.ToUpper(key)
	switch key {
	case "B0":
		return CMD_ID_QUERY_B0
	case "B2":
		return CMD_ID_QUERY_B2
	case "B4":
		return CMD_ID_QUERY_B4
	default:
		return CMD_ID_QUERY
	}
}

func GetSetCommandId(key string) CmdID {
	key = strings.ToUpper(key)
	switch key {
	case "B0":
		return CMD_ID_SET_B0
	case "B2":
		return CMD_ID_SET_B2
	case "B4":
		return CMD_ID_SET_B4
	default:
		return CMD_ID_SET
	}
}

type ApPacket struct {
	ApType   ApType
	VpType   VpType
	VpPacket VpPacket
	// CRC16: X16+X12+X5+1, XMODEM
}

func (m ApPacket) Dump() string {
	return fmt.Sprintf(
		"AP=%02X, VP=%02X, MCP=%02X, MainID=%08X, SubID=%v, Serial=%d, VpFlag=0x%02X, CmdID=0x%02X, CmdStatus=0x%02X, Data=%s",
		m.ApType, m.VpType, m.VpPacket.CPID,
		m.VpPacket.MainID, m.VpPacket.SubID, m.VpPacket.Serial, m.VpPacket.VpFlag,
		m.VpPacket.McpPacket.CmdID,
		m.VpPacket.McpPacket.CmdStatus,
		hex.EncodeToString(m.VpPacket.McpPacket.Data),
	)
}

func (m *ApPacket) UnmarshalBinary(in []byte) error {
	length := len(in)
	if length < 4 {
		return ErrInvalidPacketLength
	}

	beginFlag := in[0]
	endFlag := in[length-1]

	in = in[1 : length-1]
	m.ApType = ApType(in[0])
	switch m.ApType {
	case AP_A:
		if beginFlag != APA_PACKET_FLAG_BEGIN || endFlag != APA_PACKET_FLAG_END {
			return ErrInvalidPacketFlag
		}
		in = Unescapse(in)
	case AP_B:
		if beginFlag != APB_PACKET_FLAG_BEGIN || endFlag != APC_PACKET_FLAG_END {
			return ErrInvalidPacketFlag
		}
		out, err := hex.DecodeString(string(in))
		if err != nil {
			return errors.Wrap(err, "decode ascii")
		}
		in = out
	case AP_C:
		if beginFlag != APC_PACKET_FLAG_BEGIN || endFlag != APC_PACKET_FLAG_END {
			return ErrInvalidPacketFlag
		}
		in = Unescapse(in)
	default:
		return ErrInvalidApType
	}
	m.VpType = VpType(in[1])
	switch m.VpType {
	case VP_A:
	case VP_A1:
	default:
		return ErrInvalidVpType
	}

	var checksum uint16
	length = len(in)
	r := bytes.NewReader(in[length-2 : length])
	if err := binary.Read(r, binary.LittleEndian, &checksum); err != nil {
		return errors.Wrap(err, "read checksum")
	}
	in = in[:length-2]
	checksum2 := crc16.Checksum(in[:], crcTable)
	if checksum != checksum2 {
		return errors.New(fmt.Sprintf("checksum error, %04X != %04X", checksum, checksum2))
	}

	if err := m.VpPacket.UnmarshalBinary(in[2:]); err != nil {
		return errors.Wrap(err, "unmarshal vp packet, "+hex.EncodeToString(in[2:]))
	}
	return nil
}

func (m ApPacket) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	if err := binary.Write(w, binary.LittleEndian, m.ApType); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.VpType); err != nil {
		return w.Bytes(), err
	}
	if b, err := m.VpPacket.MarshalBinary(); err != nil {
		return w.Bytes(), err
	} else {
		if err := binary.Write(w, binary.LittleEndian, b); err != nil {
			return w.Bytes(), err
		}
	}
	checksum := crc16.Checksum(w.Bytes(), crcTable)
	if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
		return w.Bytes(), err
	}
	out := w.Bytes()
	switch m.ApType {
	case AP_A:
		if len(out) > APA_PACKET_MAX_LENGTH {
			return out, ErrInvalidPacketLength
		}
		out = Escapse(out)
		out = append([]byte{APA_PACKET_FLAG_BEGIN}, out...)
		out = append(out, APA_PACKET_FLAG_END)
	case AP_B:
		out = []byte(hex.EncodeToString(out[:]))
		out = append([]byte{APB_PACKET_FLAG_BEGIN}, out...)
		out = append(out, APB_PACKET_FLAG_END)
		if len(out) > APB_PACKET_MAX_LENGTH {
			return out, ErrInvalidPacketLength
		}
	case AP_C:
		if len(out) > APC_PACKET_MAX_LENGTH {
			return out, ErrInvalidPacketLength
		}
		out = Escapse(out)
		out = append([]byte{APC_PACKET_FLAG_BEGIN}, out...)
		out = append(out, APC_PACKET_FLAG_END)
	default:
		return out, ErrInvalidApType
	}

	return out, nil
}

type VpPacket struct {
	MainID    []byte
	SubID     uint8
	Serial    uint16
	VpFlag    VpFlag
	CPID      McpType
	McpPacket McpPacket
}

func reverseBytes(s []byte) []byte {
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = s[len(s)-1-i]
	}
	return r
}

func (m *VpPacket) UnmarshalBinary(in []byte) error {
	length := len(in)
	if length < 9 {
		return ErrInvalidPacketLength
	}
	m.MainID = reverseBytes(in[0:4])
	m.SubID = in[4]
	r := bytes.NewReader(in[5:7])
	if err := binary.Read(r, binary.LittleEndian, &m.Serial); err != nil {
		return errors.Wrap(err, "unmarshal serail")
	}
	m.VpFlag = VpFlag(in[7])
	m.CPID = McpType(in[8])
	if err := m.McpPacket.UnmarshalBinary(in[9:]); err != nil {
		return errors.Wrap(err, "unmarshal mcp packet, "+hex.EncodeToString(in[8:]))
	}
	return nil
}

func (m VpPacket) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	if err := binary.Write(w, binary.LittleEndian, reverseBytes(m.MainID)); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.SubID); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.Serial); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.VpFlag); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.CPID); err != nil {
		return w.Bytes(), err
	}
	if b, err := m.McpPacket.MarshalBinary(); err != nil {
		return w.Bytes(), err
	} else {
		if err := binary.Write(w, binary.LittleEndian, b); err != nil {
			return w.Bytes(), err
		}
	}

	return w.Bytes(), nil
}

type McpPacket struct {
	CmdID     CmdID
	CmdStatus CmdStatus // only in MCPA, MCPC
	Data      McpData
}

func (m *McpPacket) UnmarshalBinary(in []byte) error {
	length := len(in)
	if length < 2 {
		return ErrInvalidPacketLength
	}
	m.CmdID = CmdID(in[0])
	m.CmdStatus = CmdStatus(in[1])

	if err := m.Data.UnmarshalBinary(in[2:]); err != nil {
		return errors.Wrap(err, "unmarshal mcp data")
	}
	return nil
}

func (m McpPacket) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)

	if err := binary.Write(w, binary.LittleEndian, m.CmdID); err != nil {
		return w.Bytes(), err
	}
	if err := binary.Write(w, binary.LittleEndian, m.CmdStatus); err != nil {
		return w.Bytes(), err
	}
	if b, err := m.Data.MarshalBinary(); err != nil {
		return w.Bytes(), err
	} else {
		if err := binary.Write(w, binary.LittleEndian, b); err != nil {
			return w.Bytes(), err
		}
	}
	return w.Bytes(), nil
}

func Unescapse(data []byte) []byte {
	out := bytes.Replace(data, []byte{0x5e, 0x7d}, []byte{0x7e}, -1)
	out = bytes.Replace(out, []byte{0x5e, 0x5d}, []byte{0x5e}, -1)
	return out
}

func Escapse(data []byte) []byte {
	out := bytes.Replace(data, []byte{0x5e}, []byte{0x5e, 0x5d}, -1)
	out = bytes.Replace(out, []byte{0x7e}, []byte{0x5e, 0x7d}, -1)
	return out
}
