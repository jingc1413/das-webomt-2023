package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func MakeModuleForChannelConfig(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reModuleGroup := regexp.MustCompile(`Input Module\s*(\d+)`)
	reChannelGroup := regexp.MustCompile(`Channel\s*(\d+)`)
	reModule2Group := regexp.MustCompile(`(PAU|SAU)(\d*)_module(\d+)`)

	reModuleChannel := regexp.MustCompile(`Input Module\s*(\d+) Channel\s*(\d+)`)
	reCarrier := regexp.MustCompile(`Carrier\s*(\d+)`)
	reModule := regexp.MustCompile(`Module\s*(\d+)`)
	reDeviceModule := regexp.MustCompile(`(PAU|SAU1|SAU2)\s*Module\s*(\d+)`)

	moduleListTable := NewTable("Input Module List")
	moduleListTable.AddTableColumnUnique("Module", -1)

	var channelConfigControl *Element
	channelListTable := NewTable("Channel List")
	channelListTable.AddTableColumnUnique("Channel", -1)

	supportCarrierConfig := false
	carrierConfigListTable := NewTable("Channel Configuration")
	carrierConfigListTable.AddTableColumnUnique("Carrier", -1)

	addButton := NewButtonWithAction("Add", "Add", "primary", NewAction("AddItem"))
	carrierConfigListTable.SetAction("toolbar", NewToolbarWithItems("", addButton))
	carrierConfigListTable.SetStyle("viewKey", "carrier")
	carrierConfigListTable.SetStyle("invalidKey", "operator_name")
	carrierConfigListTable.SetStyle("invalidValue", "")

	supportCarrierStatus := false
	carrierStatusListTable := NewTable("Carrier Status")
	carrierStatusListTable.AddTableColumnUnique("Carrier", -1)
	carrierStatusListTable.SetStyle("invalidKey", "operator_name")
	carrierStatusListTable.SetStyle("invalidValue", "")

	supportPowerSharing := false
	powerSharingListTable := NewTable("Power Sharing")
	powerSharingListTable.AddTableColumnUnique("Module", -1)

	supportModule2List := false
	module2ListModuleNumber := int64(0)
	module2ListTable := NewTable("Input Module List")
	module2ListTable.AddTableColumnUnique("Module", -1)

	tabPages := []*Element{}
	var channelNum int64 = 0
	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Channel Config Control" {
					groupPath := append([]string{}, path...)
					groupPath = append(groupPath, pg.Name, gp.Name)
					for _, obj := range gp.Child {
						if obj.Name == "Channel Config Control" {
							channelConfigControl = MustMakeElementFromParam(obj, info, groupPath)
						}
					}
				}
			}
		} else if pg.Name == "Carrier Config" {
			supportCarrierConfig = true
			for _, gp := range pg.Child {
				if gp.Name == "Carrier Config Control" {
					for _, obj := range gp.Child {
						p := obj.Params[0]
						for _, sub := range p.Child {
							switch sub.Name {
							case "Input Power":
							case "UL/DL Configuration", "Spectial Subframe Configuration":
							case "SCS Type", "SSB Arfcn", "Spectial Slot Symbol Format", "Slot 1 Configuration", "Slot 2 Configuration":
							default:
								elem := carrierConfigListTable.AddTableColumn(sub.Name, -1)
								if sub.Name == "Carrier Index" || sub.Name == "Carrier Mask" {
									elem.SetStyle("hidden", true)
								}
							}
						}
					}
				}
			}
		} else if pg.Name == "Carrier Status" {
			supportCarrierStatus = true
			for _, gp := range pg.Child {
				if gp.Name == "Carrier Status" {
					for _, obj := range gp.Child {
						p := obj.Params[0]
						for _, sub := range p.Child {
							elem := carrierStatusListTable.AddTableColumn(sub.Name, -1)
							if sub.Name == "Carrier Index" {
								elem.SetStyle("hidden", true)
							}
						}
					}
				}
			}
		} else if pg.Name == "Power Sharing Mode" {
			supportPowerSharing = true
			for _, gp := range pg.Child {
				if gp.Name == "Power Sharing Mode" {
					for _, obj := range gp.Child {
						if strings.HasSuffix(obj.Name, "Power Percent") {
							powerSharingListTable.AddTableColumn("Power Percent", -1)
						} else if reModule.FindString(obj.Name) != "" {
							powerSharingListTable.AddTableColumn("Power Mode", -1)
						}
					}
				}
			}
		} else if pg.Name == "Input Module List" {
			supportModule2List = true
			for _, gp := range pg.Child {
				if subs := reModule2Group.FindStringSubmatch(gp.Name); len(subs) == 4 {
					moduleIndex := cast.ToInt64(subs[3])
					if moduleIndex > module2ListModuleNumber {
						module2ListModuleNumber = moduleIndex
					}
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reDeviceModule.ReplaceAllString(obj.Name, ""))
						module2ListTable.AddTableColumn(name, -1)
					}
				}
			}
		} else if reModulePage.FindString(pg.Name) != "" {
			for _, gp := range pg.Child {
				if reModuleGroup.FindString(gp.Name) != "" {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
						moduleListTable.AddTableColumn(name, -1)
					}
				} else if reChannelGroup.FindString(gp.Name) != "" {
					var channelIndex int64 = -1
					if subs := reChannelGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						if channelIndex, _ = strconv.ParseInt(subs[1], 10, 32); channelIndex > channelNum {
							channelNum = channelIndex
						}
					}
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleChannel.ReplaceAllString(obj.Name, ""))
						if strings.HasPrefix(name, "Digital Signal Bandwidth") {
							channelListTable.AddTableColumn("Digital Signal Bandwidth", -1)
						} else {
							channelListTable.AddTableColumn(name, -1)
						}
					}
				}
			}
		}
	}
	carrierConfigListTable.AddTableColumn("Action", -1)

	channelViewItemsMap := map[string][]*Element{}
	carrierConfigViewItemsMap := map[string][]*Element{}
	carrierStatusViewItemsMap := map[string][]*Element{}
	powerSharingViewItemsMap := map[string][]*Element{}

	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Channel Config Control" {
					continue
				}
				page, err := MakePageForPageGroup(pg, gp, info, path)
				if err != nil {
					return nil, errors.Wrapf(err, "make page for %v", gp.Name)
				}
				if page != nil {
					tabPages = append(tabPages, page)
				}
			}
		} else if pg.Name == "Carrier Config" {
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				if gp.Name == "Carrier Config Control" {
					for _, obj := range gp.Child {
						p := obj.Params[0]
						p.SetPath(groupPath)

						var index int64 = -1
						if subs := reCarrier.FindStringSubmatch(p.Name); len(subs) == 2 {
							index, _ = strconv.ParseInt(subs[1], 10, 32)
						}
						if index < 1 {
							continue
						}
						rowIndex := index - 1

						key := fmt.Sprintf("%v", index)
						viewItems, ok := carrierConfigViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						for _, v := range p.Child {
							elem := makeParameter(v, info)
							switch elem.Name {
							case "Input Power":
								elem.SetStyle("visibleParam", p.Child[12].PrivOid)
								elem.SetStyle("visibleValue", "01")
							case "UL/DL Configuration", "Spectial Subframe Configuration":
								elem.SetStyle("visibleParam", p.Child[15].PrivOid)
								elem.SetStyle("visibleValue", "00")
							case "SCS Type", "SSB Arfcn", "Spectial Slot Symbol Format", "Slot 1 Configuration", "Slot 2 Configuration":
								elem.SetStyle("visibleParam", p.Child[15].PrivOid)
								elem.SetStyle("visibleValue", "01")
							case "Carrier Index":
								continue
							default:
								carrierConfigListTable.MustSetTableRowData(rowIndex, v.Name, elem)
							}
							viewItems = append(viewItems, elem)
						}
						carrierConfigViewItemsMap[key] = viewItems
					}
				}
			}
		} else if pg.Name == "Carrier Status" {
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				if gp.Name == "Carrier Status" {
					for _, obj := range gp.Child {
						p := obj.Params[0]
						p.SetPath(groupPath)

						var index int64 = -1
						if subs := reCarrier.FindStringSubmatch(p.Name); len(subs) == 2 {
							index, _ = strconv.ParseInt(subs[1], 10, 32)
						}
						if index < 1 {
							continue
						}
						rowIndex := index - 1

						key := fmt.Sprintf("%v", index)
						viewItems, ok := carrierStatusViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						for _, v := range p.Child {
							elem := makeParameter(v, info)
							carrierStatusListTable.MustSetTableRowData(rowIndex, v.Name, elem)
							viewItems = append(viewItems, elem)
						}
						carrierStatusViewItemsMap[key] = viewItems
					}
				}
			}
		} else if pg.Name == "Power Sharing Mode" {
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				if gp.Name == "Power Sharing Mode" {
					for _, obj := range gp.Child {
						p := obj.Params[0]
						p.SetPath(groupPath)

						var moduleIndex int64 = -1
						if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
							moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
						}
						if moduleIndex < 1 {
							continue
						}
						rowIndex := moduleIndex - 1

						key := fmt.Sprintf("%v", moduleIndex)
						viewItems, ok := powerSharingViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						if strings.HasSuffix(obj.Name, "Power Percent") {
							name := "Power Percent"
							num := 0
							labelItems := map[int]*Element{}
							valueItems := map[int]*Element{}
							reLabel := regexp.MustCompile(`Operator\s*(\d+)`)
							reValue := regexp.MustCompile(`Power Percent\s*(\d+)`)
							for _, v := range p.Child {
								if subs := reLabel.FindStringSubmatch(v.Name); len(subs) == 2 {
									index := cast.ToInt(subs[1])
									if index < 0 {
										continue
									}
									if num < index {
										num = index
									}
									labelItems[index] = makeParameter(v, info)
								} else if subs := reValue.FindStringSubmatch(v.Name); len(subs) == 2 {
									index := cast.ToInt(subs[1])
									if index < 0 {
										continue
									}
									if num < index {
										num = index
									}
									valueItems[index] = makeParameter(v, info)
								}
							}

							statsItems := []*Element{}
							for i := 1; i <= num; i++ {
								labelItem := labelItems[i]
								valueItem := valueItems[i]
								if labelItem != nil && valueItem != nil {
									statsItems = append(statsItems,
										NewStatisticComponent(fmt.Sprintf("Port %v", i), labelItem, valueItem),
									)
								}
								viewItems = append(viewItems, labelItem, valueItem)
							}
							elem := NewStatisticGroupComponent("Power Percent", statsItems...)
							powerSharingListTable.MustSetTableRowData(rowIndex, name, elem)
						} else if reModule.FindString(obj.Name) != "" {
							name := "Power Mode"
							elem := MustMakeElementFromParam(obj, info, groupPath)
							elem.SetName(name)
							viewItems = append(viewItems, elem)
							powerSharingListTable.MustSetTableRowData(rowIndex, name, elem)
						}
						powerSharingViewItemsMap[key] = viewItems
					}
				}
			}
		} else if pg.Name == "Input Module List" {
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				var (
					deviceType  string = ""
					deviceIndex int64  = -1
					moduleIndex int64  = -1
				)
				if subs := reModule2Group.FindStringSubmatch(gp.Name); len(subs) == 4 {
					deviceType = subs[1]
					deviceIndex = cast.ToInt64(subs[2])
					moduleIndex = cast.ToInt64(subs[3])
				}
				if deviceType == "" || moduleIndex < 1 {
					continue
				}
				rowIndex := int64(-1)
				labelString := ""
				if deviceType == "PAU" {
					rowIndex = moduleIndex - 1
					labelString = fmt.Sprintf("%v Module %v", deviceType, moduleIndex)
				} else if deviceType == "SAU" {
					rowIndex = module2ListModuleNumber*deviceIndex + moduleIndex - 1
					labelString = fmt.Sprintf("%v%v Module %v", deviceType, deviceIndex, moduleIndex)
				}
				if rowIndex < 0 {
					continue
				}

				for _, obj := range gp.Child {
					name := strings.TrimSpace(reDeviceModule.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)
					module2ListTable.MustSetTableRowData(rowIndex, name, elem)
					module2ListTable.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", labelString))
				}
			}
		} else if reModulePage.FindString(pg.Name) != "" {
			var moduleIndex int64 = -1
			if subs := reModulePage.FindStringSubmatch(pg.Name); len(subs) == 2 {
				moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if moduleIndex < 1 {
				continue
			}

			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				if reModuleGroup.FindString(gp.Name) != "" {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
						elem := MustMakeElementFromParam(obj, info, groupPath)
						elem.SetName(name)
						rowIndex := moduleIndex - 1
						moduleListTable.MustSetTableRowData(rowIndex, name, elem)
						moduleListTable.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
					}
				} else if reChannelGroup.FindString(gp.Name) != "" {
					var channelIndex int64 = -1
					if subs := reChannelGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						channelIndex, _ = strconv.ParseInt(subs[1], 10, 32)
					}
					if channelIndex < 1 {
						continue
					}
					var (
						digitalSignalBandwidthSelect *Element
						digitalSignalBandwidthCustom *Element
					)
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleChannel.ReplaceAllString(obj.Name, ""))
						elem := MustMakeElementFromParam(obj, info, groupPath)
						elem.SetName(name)
						if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Select") {
							digitalSignalBandwidthSelect = elem
						} else if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Custom") {
							digitalSignalBandwidthCustom = elem
						} else {
							rowIndex := (moduleIndex-1)*channelNum + channelIndex - 1
							name := strings.TrimSpace(reModuleChannel.ReplaceAllString(obj.Name, ""))
							channelListTable.MustSetTableRowData(rowIndex, name, elem)
						}
						key := fmt.Sprintf("%v-%v", moduleIndex, channelIndex)
						viewItems, ok := channelViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						viewItems = append(viewItems, elem)
						channelViewItemsMap[key] = viewItems
					}
					if digitalSignalBandwidthCustom != nil && digitalSignalBandwidthSelect != nil {
						rowIndex := (moduleIndex-1)*channelNum + channelIndex - 1
						elem2 := NewParamGroupComponent("Digital Signal Bandwidth", 2, digitalSignalBandwidthSelect, digitalSignalBandwidthCustom)
						channelListTable.MustSetTableRowData(rowIndex, elem2.Name, elem2)
					}
				}
			}
		}
	}
	for k, items := range channelViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		channelIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*channelNum + channelIndex - 1

		name := fmt.Sprintf("Module %v Channel %v", moduleIndex, channelIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		if channelConfigControl != nil {
			viewForm.SetStyle("disableParam", channelConfigControl.OID)
			viewForm.SetStyle("disableValue", "00")
			disableAlert := NewAlert("Disable Alert", "Channel configuration is not started, please start configuration first.", "warning")
			disableAlert.SetStyle("visibleParam", channelConfigControl.OID)
			disableAlert.SetStyle("visibleValue", "00")
			viewPageItems = append(viewPageItems, disableAlert)
		}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		channelListTable.MustSetTableRowData(int64(rowIndex), "Channel",
			NewLabel("Channel", fmt.Sprintf("%v-%v", moduleIndex, channelIndex)))
		channelListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	if channelConfigControl != nil {
		action := NewAction("StateAction")
		action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(channelConfigControl, "01")))
		action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(channelConfigControl, "00")))
		channelConfigControl.SetAction("click", action)
		channelConfigControl.SetStyle("input", "button")
		channelConfigControl.SetStyle("inactiveValue", "00")
		channelConfigControl.SetStyle("inactiveText", "Start Config")
		channelConfigControl.SetStyle("inactiveType", "info")
		channelConfigControl.SetStyle("activeValue", "01")
		channelConfigControl.SetStyle("activeText", "Save Config")
		channelConfigControl.SetStyle("activeType", "primary")
		channelListTable.SetAction("toolbar", NewToolbarWithItems("", channelConfigControl))
	}

	for k, items := range carrierConfigViewItemsMap {
		index, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		name := fmt.Sprintf("Carrier %v", index)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		carrierConfigListTable.MustSetTableRowData(int64(rowIndex), "Carrier",
			NewLabel("Carrier", fmt.Sprintf("%v", index)))
		carrierConfigListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))

		var carrierInvalidElem *Element
		for _, v := range items {
			if v.Key == "operator_name" {
				carrierInvalidElem = v
				break
			}
		}
		if carrierInvalidElem != nil {
			param := *carrierInvalidElem
			param.SetValue("^^^00^00^0000000000000000^0^0^0^00^0^0^00^0^00^00^0^0^00^00000000^^^")
			action := NewSetParameterValuesAction(&param)
			btn := NewButtonWithAction("Delete", "Delete", "text", action)
			btn.SetStyle("confirmTitle", "WARNING")
			btn.SetStyle("confirmMessage", "Confirm to delete channel configuration?")
			btn.SetStyle("confirmType", "warning")
			carrierConfigListTable.MustSetTableRowData(rowIndex, "Action", btn)
		}
	}

	for k, items := range carrierStatusViewItemsMap {
		index, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1

		name := fmt.Sprintf("Carrier %v", index)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		carrierStatusListTable.MustSetTableRowData(int64(rowIndex), "Carrier",
			NewLabel("Carrier", fmt.Sprintf("%v", index)))
		carrierStatusListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	for k, items := range powerSharingViewItemsMap {
		moduleIndex, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := moduleIndex - 1

		name := fmt.Sprintf("Module %v", moduleIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		powerSharingListTable.MustSetTableRowData(int64(rowIndex), "Module",
			NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
		powerSharingListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	moduleListTable.SetStyle("height", "180px")

	module := NewModule(m.Name)

	tabsLayout := NewTabsLayout("System Signal", "top")
	tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Channel List",
		NewColLayoutWithItems(0, moduleListTable, channelListTable),
	))
	if supportCarrierConfig {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Channel Configuration",
			NewColLayoutWithItems(0, carrierConfigListTable),
		))
	}
	if supportCarrierStatus {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Carrier Status",
			NewColLayoutWithItems(0, carrierStatusListTable),
		))
	}
	if supportPowerSharing {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Power Sharing",
			NewColLayoutWithItems(0, powerSharingListTable),
		))
	}
	if supportModule2List {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Input Module List",
			NewColLayoutWithItems(0, module2ListTable),
		))
	}
	tabsLayout.Items = append(tabsLayout.Items, tabPages...)

	module.Items = append(module.Items,
		NewPageWithLayouts("System Signal", tabsLayout),
	)
	return module, nil
}
