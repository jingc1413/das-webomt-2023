package xml

import (
	"encoding/xml"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (m *MachineRoot) UnmarshalXmlFile(filename string, fix bool) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "read xml file, "+filename)
	}
	return m.UnmarshalXmlContent(content, fix)
}

func (m *MachineRoot) UnmarshalXmlContent(content []byte, fix bool) error {
	tmp := string(content)
	tmp = strings.ReplaceAll(tmp, " CarrierC ", " Carrier ")
	if fix {
		// tmp = strings.ReplaceAll(tmp, "</p_uplink>", "</g_uplink>")
		tmp = strings.ReplaceAll(tmp, "（", "(")
		tmp = strings.ReplaceAll(tmp, "）", ")")
		tmp = strings.ReplaceAll(tmp, "(SLAVE AU1)", "")
		tmp = strings.ReplaceAll(tmp, "(SLAVE AU2)", "")
		// tmp = strings.ReplaceAll(tmp, "<p_uplink ", "<g_uplink ")
		// tmp = strings.ReplaceAll(tmp, "</p_uplink>", "</g_uplink>")
		tmp = strings.ReplaceAll(tmp, "（", "(")
		tmp = strings.ReplaceAll(tmp, "）", ")")
		tmp = strings.ReplaceAll(tmp, "(SLAVE AU1)", "")
		tmp = strings.ReplaceAll(tmp, "(SLAVE AU2)", "")
		tmp = strings.ReplaceAll(tmp, "(MP-PSAU)", "")
		tmp = strings.ReplaceAll(tmp, "(HP-RUAU)", "")
		tmp = strings.ReplaceAll(tmp, "TX Power", "Tx Power")
		tmp = strings.ReplaceAll(tmp, "RX Power", "Rx Power")
		tmp = strings.ReplaceAll(tmp, "TX POWER", "Tx Power")
		tmp = strings.ReplaceAll(tmp, "RX POWER", "Rx Power")
		tmp = strings.ReplaceAll(tmp, "TX power", "Tx Power")
		tmp = strings.ReplaceAll(tmp, "RX power", "Rx Power")
		tmp = strings.ReplaceAll(tmp, " and ", " And ")
		tmp = strings.ReplaceAll(tmp, "Site Name Label", "Site Name")
		tmp = strings.ReplaceAll(tmp, "Radio Module Port", "Port")
		tmp = strings.ReplaceAll(tmp, "Interface Module ", "Interface Modules ")
		tmp = strings.ReplaceAll(tmp, "settings", "Settings")
		tmp = strings.ReplaceAll(tmp, "Interface Module\"", "Interface Modules\"")
		tmp = strings.ReplaceAll(tmp, "Module Port", "Port")
		tmp = strings.ReplaceAll(tmp, "MOV", "Movement")
		tmp = strings.ReplaceAll(tmp, "temperature", "Temperature")
		tmp = strings.ReplaceAll(tmp, "Serial Numbers", "Serial Number")
		tmp = strings.ReplaceAll(tmp, "Serial Numbes", "Serial Number")
		tmp = strings.ReplaceAll(tmp, "External Alarm\"", "External Alarms\"")
		tmp = strings.ReplaceAll(tmp, "Vender", "Vendor")
		tmp = strings.ReplaceAll(tmp, "Manufacturer", "Vendor Name")
		tmp = strings.ReplaceAll(tmp, "Configration", "Configuration")
		tmp = strings.ReplaceAll(tmp, "Infomation", "Information")
		tmp = strings.ReplaceAll(tmp, "Carrier Config Switch", "Carrier Config Global Switch")
		tmp = strings.ReplaceAll(tmp, "-power", " Power")
		tmp = strings.ReplaceAll(tmp, " power ", " Power ")
		tmp = strings.ReplaceAll(tmp, " power\"", " Power\"")
		tmp = strings.ReplaceAll(tmp, " Freq ", " Frequency ")
		tmp = strings.ReplaceAll(tmp, " Freq\"", " Frequency\"")
		tmp = strings.ReplaceAll(tmp, " Freq_low", " Frequency Start")
		tmp = strings.ReplaceAll(tmp, " Freq_high", " Frequency End")
		tmp = strings.ReplaceAll(tmp, " Frequency_low", " Frequency Start")
		tmp = strings.ReplaceAll(tmp, " Frequency_high", " Frequency End")
		tmp = strings.ReplaceAll(tmp, " Freq Low", " Frequency Start")
		tmp = strings.ReplaceAll(tmp, " Freq High", " Frequency End")
		tmp = strings.ReplaceAll(tmp, " Frequency Low", " Frequency Start")
		tmp = strings.ReplaceAll(tmp, " Frequency High", " Frequency End")
		tmp = strings.ReplaceAll(tmp, " Thr ", " Threshold ")
		tmp = strings.ReplaceAll(tmp, " Thr\"", " Threshold\"")
		tmp = strings.ReplaceAll(tmp, "slave", "Slave")
		tmp = strings.ReplaceAll(tmp, "SLAVE", "Slave")
		tmp = strings.ReplaceAll(tmp, "master", "Master")
		tmp = strings.ReplaceAll(tmp, "MASTER", "Master")
		tmp = strings.ReplaceAll(tmp, " over ", " Over ")
		tmp = strings.ReplaceAll(tmp, "failure", "Failure")
		tmp = strings.ReplaceAll(tmp, "(max)", "(Max)")
		tmp = strings.ReplaceAll(tmp, "(average)", "(Average)")
		tmp = strings.ReplaceAll(tmp, "(warning)", "(Warning)")
		tmp = strings.ReplaceAll(tmp, "(major)", "(Major)")
		tmp = strings.ReplaceAll(tmp, "(minor)", "(Minor)")
		tmp = strings.ReplaceAll(tmp, "(critical)", "(Critical)")
		tmp = strings.ReplaceAll(tmp, "Device Routing Address", "Route Address")
		tmp = strings.ReplaceAll(tmp, "OP-AU1", "OP7")
		tmp = strings.ReplaceAll(tmp, "OP-AU2", "OP8")
		tmp = strings.ReplaceAll(tmp, "OP-Master", "Optical Module Master")
		tmp = strings.ReplaceAll(tmp, "OP Master", "Optical Module Master")
		tmp = strings.ReplaceAll(tmp, "Optical Module Master Transceiver", "OP Master Transceiver")
		tmp = strings.ReplaceAll(tmp, "Optical Module Master", "OP Master Transceiver")
		tmp = strings.ReplaceAll(tmp, "OP-Slave", "Optical Module Slave")
		tmp = strings.ReplaceAll(tmp, "OP Slave", "Optical Module Slave")
		tmp = strings.ReplaceAll(tmp, "Optical Module Slave Transceiver", "OP Slave Transceiver")
		tmp = strings.ReplaceAll(tmp, "Optical Module Slave", "OP Slave Transceiver")
		tmp = strings.ReplaceAll(tmp, "Transceiver Transceiver", "Transceiver")
		tmp = strings.ReplaceAll(tmp, "Transceiver  Transceiver", "Transceiver")
		tmp = strings.ReplaceAll(tmp, "Starting(24h)/Ending", "Start - End (24h)")
		tmp = strings.ReplaceAll(tmp, "WEBOMT", "OMT")
		tmp = strings.ReplaceAll(tmp, "Offset(%)", "Distribution")
		tmp = strings.ReplaceAll(tmp, "\"Time Zone", "\"NTP Time Zone")
		tmp = strings.ReplaceAll(tmp, "Operator\"", "(Operator/Service)\"")
		tmp = strings.ReplaceAll(tmp, "Label (Operator/Service)\"", "(Operator/Service)\"")
		tmp = strings.ReplaceAll(tmp, "(Operator/Service)\"", "Label (Operator/Service)\"")
		tmp = strings.ReplaceAll(tmp, "TDD Manual Mode", "4G TDD")
		tmp = strings.ReplaceAll(tmp, "SFTP Username", "SFTP Account Username")
		tmp = strings.ReplaceAll(tmp, "SFTP Password", "SFTP Account Password")
		tmp = strings.ReplaceAll(tmp, "IP setting", "IP Setting")
		tmp = strings.ReplaceAll(tmp, "Setting ", "Settings ")
		tmp = strings.ReplaceAll(tmp, "IP Addr", "IP Address")
		tmp = strings.ReplaceAll(tmp, "Addressess", "Address")
		tmp = strings.ReplaceAll(tmp, "\"NMS ", "\"Primary NMS ")
		tmp = strings.ReplaceAll(tmp, "\"Primary NMS Port\"", "\"Primary NMS Port Number\"")
		tmp = strings.ReplaceAll(tmp, "\"Server Port (SFTP)\"", "\"Server Port Number(SFTP)\"")
		tmp = strings.ReplaceAll(tmp, "Device Port", "Device Recv Port")
		tmp = strings.ReplaceAll(tmp, "\"Upgrade ", "\"Firmware Upgrade ")
		tmp = strings.ReplaceAll(tmp, "\"Device Recv Port\"", "\"Device Recv Port(UDP)\"")
		tmp = strings.ReplaceAll(tmp, "\"Heartbeat Interval Time\"", "\"Heartbeat Clock\"")
		tmp = strings.ReplaceAll(tmp, "\"Alarm Threshold\"", "\"Alarm Thresholds\"")
		tmp = strings.ReplaceAll(tmp, "\"Link Alarm", "\"Optical Link Alarm")
		tmp = strings.ReplaceAll(tmp, "EngineID", "Engine ID")
		tmp = strings.ReplaceAll(tmp, "Over-Temperature", "Over Temperature")
		tmp = strings.ReplaceAll(tmp, "UL Carrier Power", "Carrier UL Power")
		tmp = strings.ReplaceAll(tmp, "DL Carrier Power", "Carrier DL Power")
	}
	content = []byte(tmp)

	re1 := regexp.MustCompile(`<[pP]([0-9a-f]*) `)
	re2 := regexp.MustCompile(`</[pP]([0-9a-f]*)>`)
	re3 := regexp.MustCompile(`<[pP]_(\S*) `)
	re4 := regexp.MustCompile(`</[pP]_([0-9a-zA-Z_]*)>`)
	re5 := regexp.MustCompile(`<[gG]_(\S*) `)
	re6 := regexp.MustCompile(`</[gG]_([0-9a-zA-Z_]*)\s*>`)
	re7 := regexp.MustCompile(`<m_(\S*) `)
	re8 := regexp.MustCompile(`</m_([0-9a-zA-Z_]*)>`)

	content = re1.ReplaceAll(content, []byte("<param id=\"p$1\" "))
	content = re2.ReplaceAll(content, []byte("</param>"))
	content = re3.ReplaceAll(content, []byte("<page id=\"$1\" "))
	content = re4.ReplaceAll(content, []byte("</page>"))
	content = re5.ReplaceAll(content, []byte("<group id=\"$1\" "))
	content = re6.ReplaceAll(content, []byte("</group>"))
	content = re7.ReplaceAll(content, []byte("<module id=\"$1\" "))
	content = re8.ReplaceAll(content, []byte("</module>"))

	if fix {
		re9 := regexp.MustCompile(`External Alarm\s*([0-9]+)`)
		re10 := regexp.MustCompile(`External Input Alarm\s*([0-9]+)`)
		re12 := regexp.MustCompile(`(Fan|FAN)([0-9]*) Alarm`)
		re13 := regexp.MustCompile(`(Fan|FAN)([0-9]*) Failure`)
		re14 := regexp.MustCompile(`Power([0-9]*) Failure`)
		re15 := regexp.MustCompile(`Port([0-9]+)`)
		re16 := regexp.MustCompile(`OP([0-9]+)/SAU([0-9]+)`)
		re17 := regexp.MustCompile(`(Port\s*[0-9]+) Operator`)
		re18 := regexp.MustCompile(`"Optical Module\s*([0-9]+) Transceiver`)
		re19 := regexp.MustCompile(`"Optical Module\s*([0-9]+)`)
		re20 := regexp.MustCompile(`"OP\s*([0-9]+)\s*Transceiver`)
		re21 := regexp.MustCompile(`"OP\s*([0-9]+)`)
		// re22 := regexp.MustCompile(`Module([0-9]+)`)
		re23 := regexp.MustCompile(`Address([0-9]+)`)
		re24 := regexp.MustCompile(`unit="[0-9]+\s*characters"`)

		content = re9.ReplaceAll(content, []byte("External Input $1 Alarm"))
		content = re10.ReplaceAll(content, []byte("External Input $1 Alarm"))
		content = re12.ReplaceAll(content, []byte("Fan$2 Device Loss"))
		content = re13.ReplaceAll(content, []byte("Fan$2 Device Loss"))
		content = re14.ReplaceAll(content, []byte("Power$1 Interruption"))
		content = re15.ReplaceAll(content, []byte("Port $1"))
		content = re16.ReplaceAll(content, []byte("OP$1"))
		content = re17.ReplaceAll(content, []byte("$1 Label Operator"))
		content = re18.ReplaceAll(content, []byte("\"OP$1 Transceiver"))
		content = re19.ReplaceAll(content, []byte("\"OP$1 Transceiver"))
		content = re20.ReplaceAll(content, []byte("\"OP$1"))
		content = re21.ReplaceAll(content, []byte("\"OP$1 Transceiver"))
		// content = re22.ReplaceAll(content, []byte("Module $1"))
		content = re23.ReplaceAll(content, []byte("Address $1"))
		content = re24.ReplaceAll(content, []byte("unit=\"\""))
	}

	if err := xml.Unmarshal(content, m); err != nil {
		return errors.Wrap(err, "unmarshal xml")
	}
	return nil
}

func (m MachineRoot) Dump() {
	for _, module := range m.Modules {
		if module.NodeType != "module" {
			continue
		}
		if module.ID == "digital" {
			break
		}
		logrus.Tracef("module: %v, %v", module.NameEn, module.Name)
		for _, page := range module.Pages {
			if page.NodeType != "page" {
				continue
			}
			logrus.Tracef("\tpage: %v, %v", page.NameEn, page.Name)
			for _, group := range page.Groups {
				if group.NodeType != "group" {
					continue
				}
				logrus.Tracef("\t\tgroup: %v, %v", group.NameEn, group.Name)
				for _, param := range group.Params {
					if param.NodeType != "param" {
						continue
					}
					logrus.Tracef("\t\t\tparam: %v", param)
					for _, sub := range param.Params {
						if param.NodeType != "param" {
							continue
						}
						logrus.Tracef("\t\t\t\tparam: %v", sub)
					}
				}
			}
		}
	}
	for _, module := range m.Modules {
		if module.NodeType != "module" {
			continue
		}
		if module.ID == "digital" {
			break
		}
		logrus.Tracef("module: %v, %v", module.NameEn, module.Name)
		for _, page := range module.Pages {
			if page.NodeType != "page" {
				continue
			}
			logrus.Tracef("\tpage: %v, %v", page.NameEn, page.Name)
		}
	}
}

func (m *MachineRoot) CleanUp() {
	root := *m
	modules := []Module{}
	for _, module := range root.Modules {
		if module.NodeType != "module" {
			continue
		}
		module2 := module
		pages := []Page{}
		for _, page := range module2.Pages {
			if page.NodeType == "page1" && page.NameEn == "Engineering CMD" {
				page.NodeType = "page"
			} else if page.NodeType == "page1" && page.NameEn == "Engineering CMD1" {
				page.NodeType = "page"
			} else if page.NodeType == "page1" && page.NameEn == "Engineering CMD2" {
				page.NodeType = "page"
			}
			if page.NodeType != "page" {
				continue
			}
			page2 := page
			groups := []Group{}
			for _, group := range page2.Groups {
				if group.NodeType == "group1" && group.NameEn == "PA InitFile" {
					group.NodeType = "group"
				} else if group.NodeType == "group1" && group.NameEn == "RF Module Mapping" {
					group.NodeType = "group"
				}
				if group.NodeType != "group" {
					continue
				}
				group2 := group
				params := []Param{}
				for _, param := range group2.Params {
					if param.NodeType == "param1" && param.CaptionEn == "Carrinfo export" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "Load PA InitFile" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "Carrinfo Upload" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "FlatnessCSV" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "Customcarrier import" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "Customcarrier export" {
						param.NodeType = "param"
					} else if param.NodeType == "param1" && param.CaptionEn == "Port Jump" {
						param.NodeType = "param"
					}

					if param.NodeType != "param" {
						continue
					}
					param2 := param
					if len(param2.Params) > 0 {
						subParams := []Param{}
						for _, sub := range param2.Params {
							if sub.NodeType != "param" {
								continue
							}
							sub2 := sub
							subParams = append(subParams, sub2)
						}
						param2.Params = subParams
					}
					params = append(params, param2)
				}
				group2.Params = params
				groups = append(groups, group2)
			}
			page2.Groups = groups
			pages = append(pages, page2)
		}
		module2.Pages = pages
		modules = append(modules, module2)
	}
	root.Modules = modules
	*m = root
}
