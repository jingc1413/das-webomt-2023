package model

import (
	"fmt"
	"regexp"
	"strings"
)

//  Settings,Radio Signal Information,Radio Module 1

var (
	reRadioSignalInformationRadioModule     = regexp.MustCompile(`^Settings,Radio Signal Information,Radio Module (\d+)$`)
	reRadioSignalInformationAmplifierModule = regexp.MustCompile(`^Settings,Radio Signal Information,Amplifier Module (\d+)$`)
	reInputSignalInformationAmplifierModule = regexp.MustCompile(`^Settings,Input Signal Information,Amplifier Module (\d+)$`)

	reBandConfigurationRadioModule     = regexp.MustCompile(`^Settings,Band Configuration,Radio Module (\d+)$`)
	reBandConfigurationAmplifierModule = regexp.MustCompile(`^Settings,Band Configuration,Amplifier Module (\d+)$`)
	reRadioInterfaceModuleRadioModule  = regexp.MustCompile(`^Settings,Radio Interface Modules,Radio Module\s*(\d+)$`)
	reRadioInterfaceModuleInputModule  = regexp.MustCompile(`^Settings,Input Module Information,Input Module\s*(\d+)$`)

	reCarrierConfigModuleConfiguration           = regexp.MustCompile(`^Carrier Config,Module\s*(\d+) Configuration,Module\s*(\d+) Configuration$`)
	reCarrierConfigModuleCarrierConfiguration    = regexp.MustCompile(`^Carrier Config,Module\s*(\d+) Configuration,Carrier\s*(\d+) Configuration$`)
	reCarrierConfigAmplifierChannelConfiguration = regexp.MustCompile(`^Carrier Config,Amplifier\s*(\d+) Configuration,Channel\s*(\d+) Configuration$`)
	reChannelConfigModuleChannelConfiguration    = regexp.MustCompile(`^Channel Config,Module\s*(\d+) Configuration,Channel\s*(\d+) Configuration$`)
	reChannelConfigModuleModuleConfiguration     = regexp.MustCompile(`^Channel Config,Module\s*(\d+) Configuration,Input Module\s*(\d+) Configuration$`)

	reCarrierPowerModuleCarrierPowerConfiguration = regexp.MustCompile(`^Carrier Power,Module\s*(\d+) Carrier Power Configuration,Carrier\s*(\d+) Power Configuration$`)
	reChannelPowerModuleChannelPowerConfiguration = regexp.MustCompile(`^Carrier Power,Module\s*(\d+) Channel Power Configuration,Channel\s*(\d+) Power Configuration$`)
	reManagementServiceWeekday                    = regexp.MustCompile(`^Management,Service Switch,(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday)$`)
	rePA                                          = regexp.MustCompile(`PA,PA\s*(\d+),.*`)

	reSmallSignalRadioModule = regexp.MustCompile(`^Small-Signal,(Gain|Frequency|Other Parameters),Radio Module\s*(\d+)`)
	reCombinersCombinerPort  = regexp.MustCompile(`^Combiners,Combiner\s*(\d+) Info,Port\s*(\d+)$`)
	reCombinersCombiner      = regexp.MustCompile(`^Combiners,Combiner\s*(\d+) Info,(Frequency Lower/Upper limit|Thoery Gain|Temperature|Module)$`)
	rePASwitchRxSwitch       = regexp.MustCompile(`RX Switch\s*(\d+)`)
	rePASwitchTxSwitch       = regexp.MustCompile(`TX Switch\s*(\d+)`)

	reSNMPUser   = regexp.MustCompile(`^Settings,SNMP User Info,User\s*(\d+) Info$`)
	reSNMPTrap   = regexp.MustCompile(`^Trap\s*(IP Address|IPv6 Address)\s*(\d+)$`)
	reSNMPTrapIp = regexp.MustCompile(`^Trap\s*IP Address\s*(\d+)\s*(.*)$`)

	reOpticalInfo             = regexp.MustCompile(`^OP\s*(\d+|P|S|M)\s*Transceiver\s*(.*)$`)
	reManagementServiceModule = regexp.MustCompile(`^(1st|2nd|)\s*(Primary A3|Secondary A3)\s*(Radio|Input)\s*Module\s*(\d+)\s*(UL|DL)$`)
)

func (m ParameterDefines) FixParameters() {
	for _, param := range m {
		fixParameter(param)
	}
}

func fixParameter(param *ParameterDefine) {
	for _, group := range param.Groups {
		if subs := reRadioSignalInformationRadioModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Radio Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reRadioSignalInformationAmplifierModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Amplifier Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reInputSignalInformationAmplifierModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Amplifier Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reBandConfigurationRadioModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Radio Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reBandConfigurationAmplifierModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Amplifier Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reCarrierConfigModuleConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Radio Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					fmt.Sprintf("Amplifier %v", index),
				},
			)
		} else if subs := reCarrierConfigModuleCarrierConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Radio Module %v Carrier %v", index, index2),
				[]string{
					fmt.Sprintf("Module%v-C%v", index, index2),
				},
			)
		} else if subs := reCarrierConfigAmplifierChannelConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Amplifier Module %v Channel %v", index, index2),
				[]string{
					fmt.Sprintf("Amplifier%v-C%v", index, index2),
				},
			)
		} else if subs := reChannelConfigModuleChannelConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Input Module %v Channel %v", index, index2),
				[]string{
					fmt.Sprintf("Module%v-C%v", index, index2),
				},
			)
		} else if subs := reChannelConfigModuleModuleConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			// index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Input Module %v ", index2),
				[]string{
					fmt.Sprintf("Module %v", index2),
				},
			)
		} else if subs := reCarrierPowerModuleCarrierPowerConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Radio Module %v Carrier %v", index, index2),
				[]string{
					fmt.Sprintf("Module%v-C%v", index, index2),
				},
			)
		} else if subs := reChannelPowerModuleChannelPowerConfiguration.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Input Module %v Carrier %v", index, index2),
				[]string{
					fmt.Sprintf("Module%v-C%v", index, index2),
				},
			)
		} else if subs := reRadioInterfaceModuleRadioModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Radio Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
				},
			)
		} else if subs := reRadioInterfaceModuleInputModule.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Input Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
				},
			)
		} else if subs := reSmallSignalRadioModule.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[2]
			param.FixName(fmt.Sprintf("Radio Module %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
					"Radio Module",
				},
			)
		} else if group == "Small-Signal,Gain,UpLink DSA Attention" || group == "Small-Signal,Gain,DownLink VOP Attention" {
			re := regexp.MustCompile(`^Module\s*(\d+)`)
			if subs2 := re.FindStringSubmatch(param.Name); len(subs2) == 2 {
				index := subs2[1]
				param.FixName(fmt.Sprintf("Radio Module %v", index),
					[]string{
						fmt.Sprintf("Module %v", index),
						"Radio Module",
					},
				)
				// if strings.HasSuffix(param.Name, " Value") {
				// 	param.SetName(strings.ReplaceAll(param.Name, " Value", ""))
				// }
			}
		} else if subs := reCombinersCombinerPort.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			index2 := subs[2]
			param.FixName(fmt.Sprintf("Combiner %v Port %v", index, index2),
				[]string{
					fmt.Sprintf("Port %v", index2),
					"Combiner",
					"Port",
				},
			)
		} else if subs := reCombinersCombiner.FindStringSubmatch(group); len(subs) == 3 {
			index := subs[1]
			param.FixName(fmt.Sprintf("Combiner %v", index),
				[]string{
					"Combiner",
				},
			)
		} else if group == "Settings,SNMP Configuration,Trap Settings" {
			if subs2 := reSNMPTrap.FindStringSubmatch(param.Name); len(subs2) == 3 {
				index := subs2[2]
				param.SetName(fmt.Sprintf("Trap %v %v", index, subs2[1]))
			} else if subs2 := reSNMPTrapIp.FindStringSubmatch(param.Name); len(subs2) == 3 {
				index := subs2[1]
				if subs2[2] == "" {
					param.SetName(fmt.Sprintf("Trap %v IP Address", index))
				} else {
					param.SetName(fmt.Sprintf("Trap %v %v", index, subs2[2]))
				}
			}
		} else if subs := reSNMPUser.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("User %v", index),
				[]string{
					fmt.Sprintf("Module %v", index),
				},
			)
		} else if group == "Maintenance,Optical Info,Optical Module Info" {
			if subs2 := reOpticalInfo.FindStringSubmatch(param.Name); len(subs2) == 3 {
				index := subs2[1]
				suffix := strings.TrimSpace(subs2[2])
				param.SetName(fmt.Sprintf("OP%v Transceiver %v", index, suffix))
			}
		} else if subs := reManagementServiceWeekday.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("%v", index), []string{})
		} else if group == "Management,Service Configuration,RF Module Mapping Configuration" {
			if param.Name == "Update" {
				param.SetName("Capacity Group Update")
			}
		} else if group == "Management,Service Configuration,A3 Operator/Service Configuration" {
			if subs2 := reManagementServiceModule.FindStringSubmatch(param.Name); len(subs2) == 5 {
				if subs2[1] != "" {
					param.SetName(fmt.Sprintf("%v %v %v Module %v %v", subs2[1], subs2[2], subs2[3], subs2[4], subs2[5]))
				} else {
					param.SetName(fmt.Sprintf("%v %v Module %v %v", subs2[2], subs2[3], subs2[4], subs2[5]))
				}
			}
		} else if group == "PA,PA SWITCH,RX Switch" {
			if subs2 := rePASwitchRxSwitch.FindStringSubmatch(param.Name); len(subs2) == 2 {
				param.SetName(fmt.Sprintf("RX Switch %v", subs2[1]))

			}
		} else if group == "PA,PA SWITCH,TX Switch" {
			if subs2 := rePASwitchTxSwitch.FindStringSubmatch(param.Name); len(subs2) == 2 {
				param.SetName(fmt.Sprintf("TX Switch %v", subs2[1]))

			}
		} else if subs := rePA.FindStringSubmatch(group); len(subs) == 2 {
			index := subs[1]
			param.FixName(fmt.Sprintf("PA %v", index),
				[]string{
					fmt.Sprintf("PA%v", index),
					"PA",
				},
			)
		}

		// for _, child := range param.Child {
		// 	fixParameter(child)
		// }
	}
}
