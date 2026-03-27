package model

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/spf13/cast"
)

type DeviceInfos []*DeviceInfo

func (a DeviceInfos) Len() int           { return len(a) }
func (a DeviceInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DeviceInfos) Less(i, j int) bool { return a[i].SubID < a[j].SubID }

type DeviceInfo struct {
	Schema         string
	DeviceTypeName string
	DeviceTypeID   int
	SubID          uint8

	RouteAddress []byte `json:"-"`
	IpAddress    []byte `json:"-"`

	RouteAddressString string `json:"RouteAddress,omitempty"`
	IpAddressString    string `json:"IpAddress,omitempty"`

	ConnectState int8
	AlarmState   int8
	OpticalState int8
	VersionState int8
	MixedState   int8

	DeviceName          string
	SiteName            string
	Version             string
	InstalledLocation   string
	ElementModelNumber  string
	ElementSerialNumber string

	//setup
	ProductModel     string
	ProductType      string
	OpticalPort      string
	OpticalInputPort string

	ForwardingPort uint16
	CascadingLevel int
	ParentAddress  string

	SystemTime int64
	UpTime     int64
	LifeTime   int64
}

func (m *DeviceInfo) Setup() {

	switch m.DeviceTypeName {
	case "M3RU-L":
		m.DeviceTypeName = "M3-RU-L"
	case "M3RU-H":
		m.DeviceTypeName = "M3-RU-H"
	}

	if len(m.RouteAddress) > 0 {
		m.RouteAddressString = fmt.Sprintf("%v.%v.%v.%v", m.RouteAddress[0], m.RouteAddress[1], m.RouteAddress[2], m.RouteAddress[3])
	}
	if len(m.IpAddress) > 0 {
		m.IpAddressString = fmt.Sprintf("%v.%v.%v.%v", m.IpAddress[0], m.IpAddress[1], m.IpAddress[2], m.IpAddress[3])
	}

	if m.SubID != 0 {
		m.ForwardingPort = RouteAddressToPort(m.RouteAddress)
		m.SetupOpticalPort()
	}

	if productModel := GetProductModelByDeviceTypeName(m.DeviceTypeName); productModel != nil {
		m.ProductModel = productModel.ID
		m.ProductType = productModel.ProductTypeName
	}
}

func (m *DeviceInfo) SetupOpticalPort() {
	m.OpticalPort = ""
	m.OpticalInputPort = ""
	m.CascadingLevel = 0
	if len(m.RouteAddress) != 4 {
		return
	}

	index := -1
	opIndex := 0
	opInputIndex := int(m.MixedState & 0xF)

	parent := make([]byte, len(m.RouteAddress))
	for i := 4 - 1; i >= 0; i-- {
		if opIndex == 0 && m.RouteAddress[i] != 0 {
			opIndex = int(m.RouteAddress[i])
			index = i
			parent[i] = 0
		} else {
			parent[i] = m.RouteAddress[i]
		}
	}
	if m.Schema == "corning" {
		switch m.DeviceTypeName {
		case "M3-RU-L", "M3-RU-H":
			if opInputIndex == 0 {
				m.OpticalInputPort = "OPS"
			} else {
				m.OpticalInputPort = "OPP"
			}
		case "E3-O":
			if opInputIndex == 0 {
				m.OpticalInputPort = "OP15"
			} else {
				m.OpticalInputPort = "OP16"
			}
		default:
			if opInputIndex == 0 {
				m.OpticalInputPort = "OP1"
			} else {
				m.OpticalInputPort = "OP2"
			}
		}
	}
	switch index {
	case 0:
		m.OpticalPort = fmt.Sprintf("OP%d", opIndex)
		m.CascadingLevel = 1
	case 1:
		if m.Schema == "corning" {
			switch m.DeviceTypeName {
			case "M3-RU-L", "M3-RU-H":
				if opInputIndex == 0 {
					m.OpticalPort = "OPP"
				} else {
					m.OpticalPort = "OPS"
				}
			case "E3-O":
				if opInputIndex == 0 {
					m.OpticalPort = "OP16"
				} else {
					m.OpticalPort = "OP15"
				}
			default:
				if opInputIndex == 0 {
					m.OpticalPort = "OP2"
				} else {
					m.OpticalPort = "OP1"
				}
			}
		} else {
			m.OpticalPort = "OP 1/2"
		}
		m.CascadingLevel = opIndex + 1
	case 2:
		m.OpticalPort = fmt.Sprintf("OP%d", opIndex)
		m.CascadingLevel = 1
	case 3:
		if m.Schema == "corning" {
			switch m.DeviceTypeName {
			case "M3-RU-L", "M3-RU-H":
				if opInputIndex == 0 {
					m.OpticalPort = "OPP"
				} else {
					m.OpticalPort = "OPS"
				}
			case "E3-O":
				if opInputIndex == 0 {
					m.OpticalPort = "OP16"
				} else {
					m.OpticalPort = "OP15"
				}
			default:
				if opInputIndex == 0 {
					m.OpticalPort = "OP2"
				} else {
					m.OpticalPort = "OP1"
				}
			}
		} else {
			m.OpticalPort = "OP 1/2"
		}
		m.CascadingLevel = opIndex + 1
	}
	m.ParentAddress = fmt.Sprintf("%d.%d.%d.%d", parent[0], parent[1], parent[2], parent[3])
	return
}

func RouteAddressToPort(routeAddress []byte) uint16 {
	if len(routeAddress) != 4 {
		return 0
	}
	prefix := 1
	if routeAddress[0] < 10 && routeAddress[2] < 10 {
		prefix = 3
	} else if routeAddress[0] > 9 {
		prefix += 2
	} else if routeAddress[2] > 9 {
		prefix += 1
	}
	base := cast.ToUint16(fmt.Sprintf("%v%v%v%v", routeAddress[0]%10, routeAddress[1]%10, routeAddress[2]%10, routeAddress[3]%10))
	port := base + uint16(prefix)*10000 + 1
	return port
}

type SortByRouteAddress []*DeviceInfo

func (a SortByRouteAddress) Len() int { return len(a) }
func (a SortByRouteAddress) Less(i, j int) bool {
	v1 := hex.EncodeToString(a[i].RouteAddress)
	v2 := hex.EncodeToString(a[j].RouteAddress)
	return v1 < v2
}
func (a SortByRouteAddress) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type DeviceTopoNode struct {
	ID       string            `json:"id"`
	Info     *DeviceInfo       `json:"info"`
	Children []*DeviceTopoNode `json:"children"`
}

func (m *DeviceTopoNode) Dump(level int) {
	for i := 0; i < level; i++ {
		fmt.Printf("-")
	}
	fmt.Printf(" %v: %v, %v, %v, %v\n", m.Info.SubID, m.Info.DeviceTypeName, m.Info.RouteAddress, m.Info.IpAddress, m.Info.ConnectState)
	for _, child := range m.Children {
		child.Dump(level + 1)
	}
}

func (m *DeviceTopoNode) AddChild(child *DeviceTopoNode) {
	if m.Children == nil {
		m.Children = []*DeviceTopoNode{}
	}
	m.Children = append(m.Children, child)
}

func GetDeviceTopo(all []*DeviceInfo) *DeviceTopoNode {
	var root *DeviceTopoNode = nil
	sort.Sort(SortByRouteAddress(all))

	if len(all) < 1 {
		return nil
	}
	root = &DeviceTopoNode{ID: "0", Info: all[0]}
	getNodeChild(root, all[1:])
	return root
}

func getRouteAddressLevel(addr []byte) (int, int) {
	size := len(addr)
	n := size - 1
	c := int(addr[size-1])
	for i := 0; i < size-1; i++ {
		if addr[size-i-2] == 0 {
			n -= 1
		} else {
			break
		}
	}
	return n, c
}
func matchRouteAddressLevel(a, b []byte, level int) bool {
	for i := 0; i < level; i++ {
		if a[i] != b[i] {
			// logrus.Warnf("%v: %v, %v, %v != %v", level, hex.EncodeToString(a), hex.EncodeToString(b), a[i], b[i])
			return false
		}
	}
	// logrus.Warnf("%v: %v == %v", level, hex.EncodeToString(a), hex.EncodeToString(b))
	return true
}

func getNodeChild(node *DeviceTopoNode, infos []*DeviceInfo) *DeviceTopoNode {
	if len(infos) == 0 {
		return node
	}

	for _, info := range infos {
		child := &DeviceTopoNode{ID: fmt.Sprintf("%v", info.SubID), Info: info}
		if child.Info.ParentAddress == node.Info.RouteAddressString {
			getNodeChild(child, infos)
			node.AddChild(child)
		}
	}
	return node
}
