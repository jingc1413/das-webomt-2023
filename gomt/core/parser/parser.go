package parser

import (
	"fmt"
	"gomt/core/model"
	"gomt/core/parser/ini"
	"gomt/core/parser/xml"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func GetAlarmDefines() model.AlarmDefines {
	defs, err := model.LoadDefaultAlarmDefines()
	if err != nil {
		logrus.Error("load alarm defines")
	}
	return defs
}

type DeviceFileParser struct {
	DeviceTypeName string
	Filename       string
	Root           *xml.MachineRoot
	Combo          map[string]map[string]string
}

func NewDeviceFileParsers(xmlfiles []string, inifiles []string) ([]*DeviceFileParser, error) {
	out := []*DeviceFileParser{}
	for _, xmlfile := range xmlfiles {
		p, err := NewDeviceFileParser(xmlfile, inifiles)
		if err != nil {
			return out, errors.Wrapf(err, "load device file %v", xmlfile)
		}
		out = append(out, p)
	}
	return out, nil
}

func NewDeviceFileParser(xmlfile string, inifiles []string) (*DeviceFileParser, error) {
	m := &DeviceFileParser{}
	root, err := xml.LoadXmlFile(xmlfile)
	if err != nil {
		return nil, errors.Wrap(err, "load xml file")
	}
	_, filename := path.Split(xmlfile)
	m.Filename = filename[:len(filename)-len(path.Ext(filename))]
	m.DeviceTypeName = root.RepeaterType
	m.Root = root
	m.Combo = ini.Load(inifiles)
	return m, nil
}

func (m DeviceFileParser) GetDeviceTypename() string {
	return m.Root.RepeaterType
}

func (m DeviceFileParser) GetParameters() model.ParameterDefines {
	out := model.ParameterDefines{}
	root := m.Root
	if root == nil {
		return out
	}

	rangeRe := regexp.MustCompile(`^(.*)\s+[\(]{0,1}([a-zA-Z0-9\s]+)\s*[-~/]{1}\s*([a-zA-Z0-9\s]+)[\)]{0,1}(.*)\s*$`)
	range2Re := regexp.MustCompile(`^(.*)\(([a-zA-Z0-9\s]+)\s*[-~/]{1}\s*([a-zA-Z0-9\s]+)\)(.*)\s*$`)

	opRe := regexp.MustCompile(`^OP\s*[\dSMP]+\s*Transceiver$`)
	op2Re := regexp.MustCompile(`^OP\s*(Master|Slave)\s*Transceiver$`)
	groupRe := regexp.MustCompile(`^Group\s*(\d+)`)
	moduleRe := regexp.MustCompile(`^(Channel|Module)\s*(\d+)`)
	deviceModuleRe := regexp.MustCompile(`^(PAU|SAU1|SAU2)_module(\d+)`)

	for _, module := range root.Modules {
		for _, page := range module.Pages {
			if page.NodeType != "page" {
				continue
			}
			cid := page.CID()
			for _, page2 := range page.Pages {
				page.Groups = append(page.Groups, page2)
			}
			for _, group := range page.Groups {
				if group.NodeType != "group" {
					continue
				}
				if group.NameEn == "ON/OFF,disable compensate action of the relative Module before Testing the data again,and  enable it after finishing the test" {
					group.NameEn = "Flatness Switch"
				}
				paramGroup := []string{module.NameEn, page.NameEn, group.NameEn}
				for _, param := range group.Params {
					if param.NodeType != "param" {
						continue
					}
					oidString := strings.ToUpper(fmt.Sprintf("T%v.%v", cid, param.ID))

					if len(param.Params) == 1 {
						if strings.Contains(param.CaptionEn, "Alarm") || strings.HasSuffix(param.CaptionEn, "Device Loss") {
							if !strings.Contains(param.CaptionEn, "Alarm") {
								param.CaptionEn += " Alarm"
							}
							param.Params[0].CaptionEn = param.CaptionEn
						} else if subs := rangeRe.FindStringSubmatch(param.CaptionEn); len(subs) == 5 {
							// logrus.Errorf("match: %v, %v", param.CaptionEn, strings.Join(subs[2:], " ~ "))
							if group.NameEn == "RF Module Mapping Configuration" {
								if len(subs[4]) > 0 {
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[4]+" Mapping")
								} else {
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " Mapping"
								}
							} else if group.NameEn == "Module Mapping (with Access Unit Radio Module)" {
								param.CaptionEn = strings.TrimSpace(subs[1]) + " Module Mapping Frequency (UL/DL)"
								param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " Module Mapping"
							} else {
								if len(subs[4]) > 0 {
									param.CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2]) + " " + strings.TrimSpace(subs[4])
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3]) + " " + strings.TrimSpace(subs[4])
								} else {
									param.CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2])
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3])
								}
							}
						} else if subs := range2Re.FindStringSubmatch(param.CaptionEn); len(subs) == 5 {
							// logrus.Errorf("match: %v, %v", param.CaptionEn, strings.Join(subs[2:], " ~ "))
							if group.NameEn == "RF Module Mapping Configuration" {
								if len(subs[4]) > 0 {
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[4]+" Mapping")
								} else {
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " Mapping"
								}
							} else if group.NameEn == "Module Mapping (with Access Unit Radio Module)" {
								param.CaptionEn = strings.TrimSpace(subs[1]) + " Module Mapping Frequency (UL/DL)"
								param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " Module Mapping"
							} else {
								if len(subs[4]) > 0 {
									param.CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2]) + " " + strings.TrimSpace(subs[4])
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3]) + " " + strings.TrimSpace(subs[4])
								} else {
									param.CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[2])
									param.Params[0].CaptionEn = strings.TrimSpace(subs[1]) + " " + strings.TrimSpace(subs[3])
								}
							}

						} else if strings.HasSuffix(param.CaptionEn, "Digital Signal Bandwidth") {
							// logrus.Errorf("match: %v, Digital Signal Bandwidth", param.CaptionEn)
							prefix := strings.TrimSpace(param.CaptionEn)
							param.CaptionEn = prefix + " Select"
							param.Params[0].CaptionEn = prefix + " Custom"
						} else if subs := opRe.FindStringSubmatch(param.CaptionEn); len(subs) == 1 {
							prefix := strings.TrimSpace(param.CaptionEn)
							if group.NameEn == "Optical Module Serial Numbers And Vendor Name" || group.NameEn == "Optical Module Serial Number And Vendor Name" {
								// logrus.Errorf("match: %v, Serial Numbers And Vendor Name", param.CaptionEn)
								param.CaptionEn = prefix + " Serial Number"
								param.Params[0].CaptionEn = prefix + " Vendor Name"
							} else if group.NameEn == "Optical Module Tx Power And Rx Power" {
								// logrus.Errorf("match: %v, Tx Power And Rx Power", param.CaptionEn)
								param.CaptionEn = prefix + " Tx Power"
								param.Params[0].CaptionEn = prefix + " Rx Power"
							}
						} else if subs := op2Re.FindStringSubmatch(param.CaptionEn); len(subs) == 2 {
							prefix := strings.TrimSpace(param.CaptionEn)
							if group.NameEn == "Optical Module Serial Numbers And Vendor Name" || group.NameEn == "Optical Module Serial Number And Vendor Name" {
								// logrus.Errorf("match: %v, Serial Numbers And Vendor Name", param.CaptionEn)
								param.CaptionEn = prefix + " Serial Number"
								param.Params[0].CaptionEn = prefix + " Vendor Name"
							} else if group.NameEn == "Optical Module Tx Power And Rx Power" {
								// logrus.Errorf("match: %v, Tx Power And Rx Power", param.CaptionEn)
								param.CaptionEn = prefix + " Tx Power"
								param.Params[0].CaptionEn = prefix + " Rx Power"
							}
						} else if subs := moduleRe.FindStringSubmatch(param.CaptionEn); len(subs) == 3 && param.Params[0].CaptionEn == "Power Percent" {
							param.CaptionEn = fmt.Sprintf("Module %v Power Sharing Mode", subs[2])
							param.Params[0].CaptionEn = fmt.Sprintf("Module %v Power Percent", subs[2])
						} else if subs := deviceModuleRe.FindStringSubmatch(group.NameEn); len(subs) == 3 {
							param.CaptionEn = fmt.Sprintf("%v Module %v %v", subs[1], subs[2], param.CaptionEn)
							param.Params[0].CaptionEn = fmt.Sprintf("%v Module %v %v", subs[1], subs[2], param.Params[0].CaptionEn)
						} else if subs := groupRe.FindStringSubmatch(param.CaptionEn); len(subs) == 2 {
							// ignore
						} else if param.CaptionEn == "CPU SRX Control Data" {
							param.CaptionEn = "CPU SRX Control Data"
							param.Params[0].CaptionEn = "CPU SRX Control Module"
						} else {
							logrus.Errorf("miss: %v, %v", param.ID, param.CaptionEn)
						}
					}

					def, err := m.ParameterDefineFromParam(cid, param, paramGroup)
					if err != nil {
						logrus.Error(errors.Wrapf(err, "get parameter define of %v", oidString))
						continue
					}
					if param.DataType == "23" {
						child, child2, err := m.childParameterDefinesFromParam23(def, cid, param, paramGroup)
						if err != nil {
							logrus.Error(errors.Wrapf(err, "get child parameter defines of %v", oidString))
							continue
						}
						def.Child = child
						if def.DataType == model.DataTypeBinary {
							def.DataType = model.DataTypeObject
							def.Split = "^"
						}
						if err := out.AddParameterDefine(&def); err != nil {
							logrus.Error(errors.Wrapf(err, "add child parameter define of %v", oidString))
							continue
						}
						for _, v := range child2 {
							if err := out.AddParameterDefine(v); err != nil {
								logrus.Error(errors.Wrapf(err, "add extra parameter define of %v", oidString))
								continue
							}
						}
					} else if param.DataType == "38" || param.DataType == "39" || param.DataType == "25" {
						child, child2, err := m.childParameterDefinesFromStructParam(def, cid, param, paramGroup)
						if err != nil {
							logrus.Error(errors.Wrapf(err, "get child parameter defines of %v", oidString))
							continue
						}
						def.Child = child
						if def.DataType == model.DataTypeBinary {
							def.DataType = model.DataTypeObject
							def.Split = "^"
						}
						if err := out.AddParameterDefine(&def); err != nil {
							logrus.Error(errors.Wrapf(err, "add child parameter define of %v", oidString))
							continue
						}
						for _, v := range child2 {
							if err := out.AddParameterDefine(v); err != nil {
								logrus.Error(errors.Wrapf(err, "add extra parameter define of %v", oidString))
								continue
							}
						}
					} else {
						if err := out.AddParameterDefine(&def); err != nil {
							logrus.Error(errors.Wrapf(err, "add parameter define of %v", oidString))
							continue
						}
						for _, sub := range param.Params {
							if sub.NodeType != "param" {
								continue
							}
							def2, err := m.ParameterDefineFromParam(cid, sub, paramGroup)
							if err != nil {
								logrus.Error(errors.Wrapf(err, "get parameter define of cid=%v oid=%v", cid, sub.ID))
								continue
							}
							if err := out.AddParameterDefine(&def2); err != nil {
								logrus.Error(errors.Wrapf(err, "add parameter define of cid=%v oid=%v", cid, sub.ID))
								continue
							}

						}
					}
				}
			}
		}
	}
	return out
}

func (m DeviceFileParser) childParameterDefinesFromParam23(
	parent model.ParameterDefine,
	cid string,
	param xml.Param,
	paramGroup []string,
) ([]*model.ParameterDefine, []*model.ParameterDefine, error) {
	out := []*model.ParameterDefine{}
	out2 := []*model.ParameterDefine{}
	if param.DataType != "23" {
		return out, out2, errors.New("invalid 0x23 data type of parameter")
	}

	i := 0
	for _, sub := range param.Params {
		if sub.NodeType != "param" {
			continue
		}

		if sub.ID == param.ID {
			def2, err := m.ParameterDefineFromParam(cid, sub, paramGroup)
			if err != nil {
				logrus.Error(errors.Wrapf(err, "get child parameter define of oid=T%v.P%v", cid, sub.ID))
				continue
			}
			switch i {
			case 0: // DDM
				def2.DataType = model.DataTypeBinary
				def2.ByteSize = 1
				def2.Options = model.Options{
					"00": "no",
					"01": "yes",
				}
			case 1: // Tx Power
				def2.DataType = model.DataTypeInt16
				def2.ByteSize = 2
				ratio := int64(100)
				def2.Ratio = &ratio
			case 2: // Rx Power
				def2.DataType = model.DataTypeInt16
				def2.ByteSize = 2
				ratio := int64(100)
				def2.Ratio = &ratio
			case 3: // Voltage
				def2.DataType = model.DataTypeInt16
				def2.ByteSize = 2
				ratio := int64(10)
				def2.Ratio = &ratio
			case 4: // Bias
				def2.DataType = model.DataTypeInt16
				def2.ByteSize = 2
				ratio := int64(100)
				def2.Ratio = &ratio
			case 5: //Temperature
				def2.DataType = model.DataTypeInt16
				def2.ByteSize = 2
				ratio := int64(10)
				def2.Ratio = &ratio
			case 6:
				def2.DataType = model.DataTypeUInt16
				def2.ByteSize = 2
			}
			def2.PrivOid = model.PrivObjectId(fmt.Sprintf("%v[%v]", parent.PrivOid, i))

			i += 1
			out = append(out, &def2)
		} else {
			def2, err := m.ParameterDefineFromParam(cid, sub, paramGroup)
			if err != nil {
				logrus.Error(errors.Wrapf(err, "get extra parameter define of oid=T%v.P%v", cid, sub.ID))
				continue
			}
			out2 = append(out2, &def2)
		}
	}

	return out, out2, nil
}

func (m DeviceFileParser) childParameterDefinesFromStructParam(
	parent model.ParameterDefine,
	cid string,
	param xml.Param,
	paramGroup []string,
) ([]*model.ParameterDefine, []*model.ParameterDefine, error) {
	out := []*model.ParameterDefine{}
	out2 := []*model.ParameterDefine{}
	if param.DataType != "38" && param.DataType != "39" && param.DataType != "25" {
		return out, out2, errors.New("invalid data type of struct parameter")
	}

	i := 0
	for _, sub := range param.Params {
		if sub.NodeType != "param" {
			continue
		}

		if sub.ID == param.ID {
			def2, err := m.ParameterDefineFromParam(cid, sub, paramGroup)
			if err != nil {
				logrus.Error(errors.Wrapf(err, "get child parameter define of oid=T%v.P%v", cid, sub.ID))
				continue
			}
			def2.PrivOid = model.PrivObjectId(fmt.Sprintf("%v[%v]", parent.PrivOid, i))

			i += 1
			out = append(out, &def2)
		} else {
			def2, err := m.ParameterDefineFromParam(cid, sub, paramGroup)
			if err != nil {
				logrus.Error(errors.Wrapf(err, "get extra parameter define of oid=T%v.P%v", cid, sub.ID))
				continue
			}
			out2 = append(out2, &def2)
		}
	}

	return out, out2, nil
}

var VaildTipsMap = map[string]string{
	"Reserve":                   "Reserve",
	"Turn off when not in use!": "Turn off when not in use.",
	"hexadecimal":               "Hexadecimal",
	"decimalism":                "Decimalism",
	"decimal":                   "Decimalism",
	"Mode 1: freq>3G reference value = -15 Mode 2: freq>3G reference value = -20": "Mode 1: freq>3G reference value = -15 Mode 2: freq>3G reference value = -20",
	"Time synchronizes to LMS when NTP is off":                                    "Time synchronizes to LMS when NTP is off",
	"For SNMPv3,default engine ID will be used if not set":                        "For SNMPv3, default engine ID will be used if not set",
}
var VaildUnitNameMap = map[string]string{
	"%":                  "%",
	"℃":                  "℃",
	"°C":                 "℃",
	"us":                 "us",
	"ms":                 "ms",
	"s":                  "s",
	"min":                "min",
	"hours":              "hours",
	"hour":               "hour",
	"days:hours:minutes": "days:hours:minutes",
	"V":                  "V",
	"mV":                 "mV",
	"A":                  "A",
	"mA":                 "mA",
	"nm":                 "nm",
	"cm":                 "cm",
	"m":                  "m",
	"km":                 "km",

	"MHz":       "MHz",
	"dB":        "dB",
	"dBm":       "dBm",
	"db/℃":      "dB/℃",
	"0.001db/°": "0.001dB/℃",
}

func MatchTips(name string) string {
	if name == "" {
		return ""
	}
	name = strings.TrimSpace(name)
	for k, v := range VaildTipsMap {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	for k := range VaildUnitNameMap {
		if strings.EqualFold(k, name) {
			return ""
		}
	}
	// logrus.Errorf("invalid tip: %v", name)
	return ""
}

func MatchUnitName(name string) string {
	if name == "" {
		return ""
	}
	for k, v := range VaildUnitNameMap {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	for k := range VaildTipsMap {
		if strings.EqualFold(k, name) {
			return ""
		}
	}
	logrus.Errorf("invalid unit name: %v", name)
	return ""
}

func (m DeviceFileParser) ParameterDefineFromParam(cid string, param xml.Param, paramGroup []string) (model.ParameterDefine, error) {
	def := model.ParameterDefine{}
	def.Groups = []string{
		strings.Join(paramGroup, ","),
	}
	poid, err := model.MakePrivObjectId(cid, param.ID)
	if err != nil {
		return def, errors.Wrap(err, "parse private object id")
	}
	def.PrivOid = poid
	if param.CaptionEn != "" {
		def.SetName(param.CaptionEn)
	} else {
		def.SetName(param.Caption)
	}

	if param.UnitEn != "" {
		def.UnitName = MatchUnitName(param.UnitEn)
		def.Tips = MatchTips(param.UnitEn)
	} else if param.Unit != "" {
		def.UnitName = MatchUnitName(param.Unit)
		def.Tips = MatchTips(param.Unit)
	} else if param.Range != "" {
		def.Tips = MatchTips(param.Range)
	}

	def.ByteSize = param.Len
	if param.Rate != 1 {
		def.Ratio = &param.Rate
	}
	switch param.Limit {
	case "orw":
		def.Access = model.AccessReadWrite
	case "nrw":
		def.Access = model.AccessReadWrite
	case "rw":
		def.Access = model.AccessReadWrite
	case "ro":
		if param.EditType == "04" {
			def.Access = model.AccessReadWrite
		} else {
			def.Access = model.AccessReadOnly
		}
	case "wo":
		def.Access = model.AccessWriteOnly
	default:
		def.Access = model.AccessNotAccessible
	}
	if def.Name == "Factory Mode" {
		def.Access = model.AccessReadWrite
	} else if def.Name == "Alarm Initialization" {
		def.Access = model.AccessWriteOnly
	} else if def.Name == "Alarm Mode Select" {
		def.Access = model.AccessReadWrite
	}

	if param.Range != "" {
		re := regexp.MustCompile(`\s*([-]{0,1}\d+)\s*[-~]{1}\s*([-]{0,1}\d+)\s*.*`)
		if subs := re.FindStringSubmatch(param.Range); len(subs) == 3 {
			if tmp, err := strconv.ParseInt(strings.TrimSpace(subs[1]), 10, 32); err == nil {
				if def.Ratio != nil {
					tmp = tmp * *def.Ratio
				}
				def.Min = &tmp
			}
			if tmp, err := strconv.ParseInt(strings.TrimSpace(subs[2]), 10, 32); err == nil {
				if def.Ratio != nil {
					tmp = tmp * *def.Ratio
				}
				def.Max = &tmp
			}
		} else if def.PrivOid == "T02.P0D4D" {
			def.Tips = param.Range
		}
	} else if param.CaptionEn == "Privacy Password" {
		min := int64(8)
		max := int64(15)
		def.Min = &min
		def.Max = &max
	}

	switch param.DataType {
	case "0":
		fallthrough
	case "00":
		if def.ByteSize == 1 {
			def.DataType = model.DataTypeUInt8
		} else if def.ByteSize == 2 {
			def.DataType = model.DataTypeUInt16
		} else if def.ByteSize <= 4 {
			def.DataType = model.DataTypeUInt32
		} else if def.ByteSize <= 8 {
			def.DataType = model.DataTypeUInt64
		} else {
			def.DataType = model.DataTypeUInt64
		}
	case "1":
		fallthrough
	case "01":
		if def.ByteSize == 1 {
			def.DataType = model.DataTypeInt8
		} else if def.ByteSize == 2 {
			def.DataType = model.DataTypeInt16
		} else if def.ByteSize <= 4 {
			def.DataType = model.DataTypeInt32
		} else if def.ByteSize <= 8 {
			def.DataType = model.DataTypeInt64
		} else {
			def.DataType = model.DataTypeInt64
		}
	case "2":
		fallthrough
	case "02":
		def.DataType = model.DataTypeString
	case "3":
		fallthrough
	case "03":
		def.DataType = model.DataTypeBinary
	case "4":
		fallthrough
	case "04":
		def.DataType = model.DataTypeIPV4
	case "5":
		fallthrough
	case "05":
		def.DataType = model.DataTypeDateTime
	case "6":
		fallthrough
	case "06":
		def.DataType = model.DataTypeBinary
	case "9":
		fallthrough
	case "09":
		def.DataType = model.DataTypeIPV4
	case "33":
		def.ByteSize = 8
		def.DataType = model.DataTypeUint32Array
	case "34":
		def.DataType = model.DataTypeBinary
	default:
		def.DataType = model.DataTypeBinary
	}

	if def.PrivOid == "TB2.P0CCC" {
		key := ""
		switch m.DeviceTypeName {
		case "Primary A3", "Secondary A3":
			key = "register_au"
		case "N2-RU":
			key = "register_n2ru"
		case "N3-RU":
			key = "register_n3ru"
		case "M3-RU-L", "M3-RU-H":
			key = "register_m3ru"
		case "X3-RU":
			key = "register_x3ru"
		case "E3-O":
			key = "register_au"
		case "RU":
			key = "register_ru"
		}
		if key != "" {
			if section, ok := m.Combo[key]; ok {
				def.Options = model.NewOptionsFromMap(section)
			}
		}
	}

	if param.SectionEn != "" {
		if section, ok := m.Combo[param.SectionEn]; ok {
			def.Options = model.NewOptionsFromMap(section)
		}
	} else if param.Section != "" {
		if section, ok := m.Combo[param.Section]; ok {
			def.Options = model.NewOptionsFromMap(section)
		}
	} else {
		switch param.EditType {
		case "03": // switch 1:On 0:Off
			if def.DataType == model.DataTypeBinary {
				def.Options = model.NewOptionsFromMap(map[string]string{
					"00": "Off",
					"01": "On",
				})
			} else {
				def.Options = model.NewOptionsFromMap(map[string]string{
					"0": "Off",
					"1": "On",
				})
			}
		case "05": // Alaram Status
			if def.DataType == model.DataTypeBinary {
				if strings.Contains(def.Name, "Sync") {
					def.Options = model.NewOptionsFromMap(map[string]string{
						"01": "Unsync",
						"00": "Sync",
					})
				} else {
					def.Options = model.NewOptionsFromMap(map[string]string{
						"00": "Normal",
						"01": "Alarm",
					})
				}
			} else {
				if strings.Contains(def.Name, "Sync") {
					def.Options = model.NewOptionsFromMap(map[string]string{
						"01": "Unsync",
						"00": "Sync",
					})
				} else {
					def.Options = model.NewOptionsFromMap(map[string]string{
						"0": "Normal",
						"1": "Alarm",
					})
				}
			}
		}
	}

	if def.DataType == model.DataTypeString {
		for _, v := range def.Groups {
			if v == "Management,Service Configuration/Operator,Service Configuration" {
				def.SetupArrayWithNumber(4, ",")
			} else if v == "Management,Service Configuration,A3 Operator/Service Configuration" {
				def.SetupArrayWithNumber(4, ",")
			} else if v == "Management,Service Configuration,RF Module Mapping Configuration" {
				def.SetupArrayWithNumber(2, "/")
			} else if v == "Management,Service Configuration,RF Module Mapping" {
				num := 0
				re := regexp.MustCompile(`Module\s*(\d+)`)
				for _, module := range m.Root.Modules {
					for _, page := range module.Pages {
						if page.NodeType == "page" && page.NameEn == "Service Configuration" {
							for _, group := range page.Groups {
								if group.NodeType == "group" && group.NameEn == "RF Module Mapping Configuration" {
									for _, p := range group.Params {
										if subs := re.FindStringSubmatch(p.CaptionEn); len(subs) == 2 {
											index := cast.ToInt(subs[1])
											if index > num {
												num = index
											}
										}
									}
								}
							}
						}
					}
				}
				def.SetupArrayWithNumber(num, ":")
			}
		}
		// } else if def.DataType == model.DataTypeBinary {
		// for _, v := range def.Groups {
		// 	if v == "Channel Config,Carrier Config,Carrier Config Control" ||
		// 		v == "Channel Config,Power Sharing Mode,Power Sharing Mode" ||
		// 		v == "Channel Config,Carrier Status,Carrier Status" ||
		// 		v == "Maintenance,Optical Info,Optical Module Info" {
		// 		def.DataType = model.DataTypeObject
		// 	}
		// }
	}

	def.InputType = "default"
	switch param.EditType {
	case "00", "0":
		if def.Access == model.AccessReadOnly {
			def.InputType = "default"
		} else if def.DataType.IsNumeric() {
			def.InputType = "number"
		}
	case "01", "1":
		def.InputType = "button"
	case "02", "2":
		if def.Access == model.AccessReadOnly {
			def.InputType = "default"
		} else if def.DataType.IsNumeric() {
			def.InputType = "number"
		} else if def.DataType == model.DataTypeBinary {
			def.InputType = "binary"
		} else if strings.Contains(strings.ToLower(def.Name), "password") {
			def.InputType = "password"
		}
	case "03", "3":
		activeValue := ""
		inactiveValue := ""
		if len(def.Options) == 2 {
			for k, v := range def.Options {
				text := strings.ToLower(fmt.Sprintf("%v", v))
				switch text {
				case "on":
					activeValue = k
				case "enable":
					activeValue = k
				case "off":
					inactiveValue = k
				case "disable":
					inactiveValue = k
				case "factory pattern":
					activeValue = k
				case "quit":
					inactiveValue = k
				}
			}
		}
		if activeValue != "" && inactiveValue != "" {
			def.InputType = "switch"
		} else if def.Access == model.AccessReadOnly {
			def.InputType = "default"
		} else {
			def.InputType = "select"
		}
	case "04", "4":
		if def.Access == model.AccessReadOnly {
			def.InputType = "default"
		} else if def.Name == "Alarm Initialization" {
			def.InputType = "buttonGroup"
		} else {
			def.InputType = "select"
		}
	case "05", "5":
		if !strings.Contains(strings.ToLower(def.Name), "alarm") && strings.Contains(strings.ToLower(def.Name), "sync") {
			def.InputType = "status:sync"
		} else {
			def.InputType = "status:alarm"
		}
	case "06", "6":
		def.InputType = "datetime"
	case "07", "7":
		if strings.Contains(strings.ToLower(def.Name), "password") {
			def.InputType = "password"
		}
	case "08", "8":
		def.InputType = "password"
	case "09", "9":
		if def.Access == model.AccessReadOnly {
			def.InputType = "default"
		} else {
			def.InputType = "select"
		}
	case "11":
		def.InputType = "percent"
	default:
		return def, errors.Errorf("invalid edit type %v for %v", param.EditType, def.PrivOid)
	}
	if def.DataType == model.DataTypeIPV4 {
		def.InputType = "ipv4"
	} else if def.DataType == model.DataTypeIPV6 {
		def.InputType = "ipv6"
	}

	if param.EditType == "09" {
		def.MultipleOption = true
	}
	if def.Name == "Carrier Mask" {
		def.MultipleOption = true
		def.InputType = "treeSelect"
	}

	if def.PrivOid == "TB2.P0CCC" {
		def.Child = []*model.ParameterDefine{}
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[0]", def.PrivOid)),
			Name:      "Module ID",
			Access:    model.AccessReadWrite,
			DataType:  model.DataTypeBinary,
			ByteSize:  1,
			Options:   def.Options,
			InputType: "select",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		minAddr := int64(0)
		maxAddr := int64(0xFFFF)
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[1]", def.PrivOid)),
			Name:      "Register Address",
			Access:    model.AccessReadWrite,
			DataType:  model.DataTypeUInt16,
			ByteSize:  2,
			InputType: "number",
			Max:       &maxAddr,
			Min:       &minAddr,
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[2]", def.PrivOid)),
			Name:      "Byte Size",
			Access:    model.AccessReadWrite,
			DataType:  model.DataTypeBinary,
			ByteSize:  1,
			InputType: "radio",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[3]", def.PrivOid)),
			Name:      "Byte Buffer",
			Access:    model.AccessReadWrite,
			DataType:  model.DataTypeBinary,
			ByteSize:  4,
			InputType: "binary",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.DataType = model.DataTypeObject
		def.Split = ""
		def.ByteSize = 8
		def.Options = nil
		def.InputType = "default"
	} else if def.PrivOid == "TB4.P0B23" {
		def.Child = []*model.ParameterDefine{}
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[0]", def.PrivOid)),
			Name:      "Timeout Device Number",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[1]", def.PrivOid)),
			Name:      "Failed Device Number",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[2]", def.PrivOid)),
			Name:      "Upgrade Progress",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[3]", def.PrivOid)),
			Name:      "Error Code",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[4]", def.PrivOid)),
			Name:      "Upgrade Status",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[5]", def.PrivOid)),
			Name:      "Download Status",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[6]", def.PrivOid)),
			Name:      "Total Device Number",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.Child = append(def.Child, &model.ParameterDefine{
			PrivOid:   model.PrivObjectId(fmt.Sprintf("%v[7]", def.PrivOid)),
			Name:      "Reserve",
			Access:    model.AccessReadOnly,
			DataType:  model.DataTypeUInt8,
			ByteSize:  1,
			InputType: "number",
			Groups:    def.Groups,
			Paths:     def.Paths,
		})
		def.DataType = model.DataTypeObject
		def.Split = ""
		def.ByteSize = 8
		def.InputType = "default"
	}
	return def, nil
}
