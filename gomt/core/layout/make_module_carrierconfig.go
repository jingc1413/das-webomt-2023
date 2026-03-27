package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func MakeModuleForAuCarrierConfig(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakeModuleForAuCarrierConfigUseChannel(m, info, path)
	} else {
		return MakeModuleForAuCarrierConfigUseCarrier(m, info, path)
	}
}

func MakeModuleForRuCarrierConfig(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakeModuleForRuCarrierConfigUseAmplifierModule(m, info, path)
	} else {
		return MakeModuleForRuCarrierConfigUseRadioModule(m, info, path)
	}
}

func MakeModuleForAuCarrierConfigUseCarrier(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reModuleGroup := regexp.MustCompile(`Module\s*(\d+)`)
	reCarrierGroup := regexp.MustCompile(`Carrier\s*(\d+)`)

	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)
	reModuleCarrier := regexp.MustCompile(`Radio Module\s*(\d+) Carrier\s*(\d+)`)

	moduleListTable := NewTable("Radio Module List")
	moduleListTable.AddTableColumnUnique("Module", -1)

	var configControl *Element
	childListTable := NewTable("Carrier List")
	childListTable.AddTableColumnUnique("Carrier", -1)

	tabPages := []*Element{}
	var childNum int64 = 0
	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Carrier Config Control" {
					groupPath := append([]string{}, path...)
					groupPath = append(groupPath, pg.Name, gp.Name)
					for _, obj := range gp.Child {
						if obj.Name == "Carrier Config Control" {
							configControl = MustMakeElementFromParam(obj, info, groupPath)
						}
					}
				}
			}
		} else {
			for _, gp := range pg.Child {
				if reModuleGroup.FindString(gp.Name) != "" {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
						moduleListTable.AddTableColumn(name, -1)
					}
				} else if reCarrierGroup.FindString(gp.Name) != "" {
					var index int64 = -1
					if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						if index, _ = strconv.ParseInt(subs[1], 10, 32); index > childNum {
							childNum = index
						}
					}
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
						if strings.HasPrefix(name, "Digital Signal Bandwidth") {
							childListTable.AddTableColumn("Digital Signal Bandwidth", -1)
						} else {
							childListTable.AddTableColumn(name, -1)
						}
					}
				}
			}
		}
	}

	childViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Carrier Config Control" {
				} else {
					page, err := MakePageForPageGroup(pg, gp, info, path)
					if err != nil {
						return nil, errors.Wrapf(err, "make page for %v", gp.Name)
					}
					if page != nil {
						tabPages = append(tabPages, page)
					}
				}
			}
		} else {
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
				} else if reCarrierGroup.FindString(gp.Name) != "" {
					var carrierIndex int64 = -1
					if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32)
					}
					if carrierIndex < 1 {
						continue
					}
					var (
						carrierDigitalSignalBandwidthSelect *Element
						carrierDigitalSignalBandwidthCustom *Element
					)
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
						elem := MustMakeElementFromParam(obj, info, groupPath)
						elem.SetName(name)

						if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Select") {
							carrierDigitalSignalBandwidthSelect = elem
						} else if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Custom") {
							carrierDigitalSignalBandwidthCustom = elem
						} else {
							rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
							childListTable.MustSetTableRowData(rowIndex, name, elem)
						}
						key := fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)
						viewItems, ok := childViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						viewItems = append(viewItems, elem)
						childViewItemsMap[key] = viewItems
					}
					if carrierDigitalSignalBandwidthCustom != nil && carrierDigitalSignalBandwidthSelect != nil {
						rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
						elem2 := NewParamGroupComponent("Digital Signal Bandwidth", 2, carrierDigitalSignalBandwidthSelect, carrierDigitalSignalBandwidthCustom)
						childListTable.MustSetTableRowData(rowIndex, elem2.Name, elem2)
					}
				}
			}

		}
	}
	for k, items := range childViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		carrierIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1

		name := fmt.Sprintf("Radio Module %v Carrier %v", moduleIndex, carrierIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		if configControl != nil {
			viewForm.SetStyle("disableParam", configControl.OID)
			viewForm.SetStyle("disableValue", "00")
			disableAlert := NewAlert("Disable Alert", "Carrier configuration is not started, please start configuration first.", "warning")
			disableAlert.SetStyle("visibleParam", configControl.OID)
			disableAlert.SetStyle("visibleValue", "00")
			viewPageItems = append(viewPageItems, disableAlert)
		}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		childListTable.MustSetTableRowData(int64(rowIndex), "Carrier", NewLabel("Carrier", fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)))
		childListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	if configControl != nil {
		action := NewAction("StateAction")
		action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(configControl, "01")))
		action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(configControl, "00")))
		configControl.SetAction("click", action)
		configControl.SetStyle("input", "button")
		configControl.SetStyle("inactiveValue", "00")
		configControl.SetStyle("inactiveText", "Start Config")
		configControl.SetStyle("inactiveType", "info")
		configControl.SetStyle("activeValue", "01")
		configControl.SetStyle("activeText", "Save Config")
		configControl.SetStyle("activeType", "primary")
		childListTable.SetAction("toolbar", NewToolbarWithItems("", configControl))
	}

	moduleListTable.SetStyle("height", "180px")

	module := NewModule(m.Name)
	tabsLayout := NewTabsLayout("System Signal", "top")
	tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Carrier List",
		NewColLayoutWithItems(0, moduleListTable, childListTable),
	))
	tabsLayout.Items = append(tabsLayout.Items, tabPages...)

	module.Items = append(module.Items, NewPageWithLayouts("System Signal", tabsLayout))
	return module, nil
}

func MakeModuleForAuCarrierConfigUseChannel(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reModuleGroup := regexp.MustCompile(`Input Module\s*(\d+)`)
	reCarrierGroup := regexp.MustCompile(`Channel\s*(\d+)`)

	reModule := regexp.MustCompile(`Input Module\s*(\d+)`)
	reModuleCarrier := regexp.MustCompile(`Input Module\s*(\d+) Channel\s*(\d)`)

	moduleListTable := NewTable("Input Module List")
	moduleListTable.AddTableColumnUnique("Module", -1)

	var configControl *Element
	childListTable := NewTable("Channel List")
	childListTable.AddTableColumnUnique("Channel", -1)

	tabPages := []*Element{}
	var childNum int64 = 0
	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Channel Config Control" {
					groupPath := append([]string{}, path...)
					groupPath = append(groupPath, pg.Name, gp.Name)
					for _, obj := range gp.Child {
						if obj.Name == "Channel Config Control" {
							configControl = MustMakeElementFromParam(obj, info, groupPath)
						}
					}
				}
			}
		} else {
			for _, gp := range pg.Child {
				if reModuleGroup.FindString(gp.Name) != "" {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
						moduleListTable.AddTableColumn(name, -1)
					}
				} else if reCarrierGroup.FindString(gp.Name) != "" {
					var index int64 = -1
					if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						if index, _ = strconv.ParseInt(subs[1], 10, 32); index > childNum {
							childNum = index
						}
					}
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
						if strings.HasPrefix(name, "Digital Signal Bandwidth") {
							childListTable.AddTableColumn("Digital Signal Bandwidth", -1)
						} else {
							childListTable.AddTableColumn(name, -1)
						}
					}
				}
			}
		}
	}

	childViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		if pg.Name == "System Signal Info" {
			for _, gp := range pg.Child {
				if gp.Name == "Channel Config Control" {
				} else {
					page, err := MakePageForPageGroup(pg, gp, info, path)
					if err != nil {
						return nil, errors.Wrapf(err, "make page for %v", gp.Name)
					}
					if page != nil {
						tabPages = append(tabPages, page)
					}
				}
			}
		} else {
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
				} else if reCarrierGroup.FindString(gp.Name) != "" {
					var carrierIndex int64 = -1
					if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32)
					}
					if carrierIndex < 1 {
						continue
					}
					var (
						carrierDigitalSignalBandwidthSelect *Element
						carrierDigitalSignalBandwidthCustom *Element
					)
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
						elem := MustMakeElementFromParam(obj, info, groupPath)
						elem.SetName(name)

						if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Select") {
							carrierDigitalSignalBandwidthSelect = elem
						} else if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Custom") {
							carrierDigitalSignalBandwidthCustom = elem
						} else {
							rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
							childListTable.MustSetTableRowData(rowIndex, name, elem)
						}
						key := fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)
						viewItems, ok := childViewItemsMap[key]
						if !ok {
							viewItems = []*Element{}
						}
						viewItems = append(viewItems, elem)
						childViewItemsMap[key] = viewItems
					}
					if carrierDigitalSignalBandwidthCustom != nil && carrierDigitalSignalBandwidthSelect != nil {
						rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
						elem2 := NewParamGroupComponent("Digital Signal Bandwidth", 2, carrierDigitalSignalBandwidthSelect, carrierDigitalSignalBandwidthCustom)
						childListTable.MustSetTableRowData(rowIndex, elem2.Name, elem2)
					}
				}
			}

		}
	}
	for k, items := range childViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		carrierIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1

		name := fmt.Sprintf("Input Module %v Channel %v", moduleIndex, carrierIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		if configControl != nil {
			viewForm.SetStyle("disableParam", configControl.OID)
			viewForm.SetStyle("disableValue", "00")
			disableAlert := NewAlert("Disable Alert", "Channel configuration is not started, please start configuration first.", "warning")
			disableAlert.SetStyle("visibleParam", configControl.OID)
			disableAlert.SetStyle("visibleValue", "00")
			viewPageItems = append(viewPageItems, disableAlert)
		}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		childListTable.MustSetTableRowData(int64(rowIndex), "Channel", NewLabel("Channel", fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)))
		childListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	if configControl != nil {
		action := NewAction("StateAction")
		action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(configControl, "01")))
		action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(configControl, "00")))
		configControl.SetAction("click", action)
		configControl.SetStyle("input", "button")
		configControl.SetStyle("inactiveValue", "00")
		configControl.SetStyle("inactiveText", "Start Config")
		configControl.SetStyle("inactiveType", "info")
		configControl.SetStyle("activeValue", "01")
		configControl.SetStyle("activeText", "Save Config")
		configControl.SetStyle("activeType", "primary")
		childListTable.SetAction("toolbar", NewToolbarWithItems("", configControl))
	}

	moduleListTable.SetStyle("height", "180px")

	module := NewModule(m.Name)
	tabsLayout := NewTabsLayout("System Signal", "top")
	tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Channel List",
		NewColLayoutWithItems(0, moduleListTable, childListTable),
	))
	tabsLayout.Items = append(tabsLayout.Items, tabPages...)

	module.Items = append(module.Items, NewPageWithLayouts("System Signal", tabsLayout))
	return module, nil
}

func MakeModuleForRuCarrierConfigUseRadioModule(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reModuleGroup := regexp.MustCompile(`Module\s*(\d+)`)
	reCarrierGroup := regexp.MustCompile(`Carrier\s*(\d+)`)

	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)
	reModuleCarrier := regexp.MustCompile(`Radio Module\s*(\d+) Carrier\s*(\d)`)

	moduleListTable := NewTable("Radio Module List")
	moduleListTable.AddTableColumnUnique("Module", -1)

	childListTable := NewTable("Carrier List")
	childListTable.AddTableColumnUnique("Carrier", -1)

	tabPages := []*Element{}
	var childNum int64 = 0
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			if reModuleGroup.FindString(gp.Name) != "" {
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
					moduleListTable.AddTableColumn(name, -1)
				}
			} else if reCarrierGroup.FindString(gp.Name) != "" {
				var index int64 = -1
				if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
					if index, _ = strconv.ParseInt(subs[1], 10, 32); index > childNum {
						childNum = index
					}
				}
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
					if strings.HasPrefix(name, "Digital Signal Bandwidth") {
						childListTable.AddTableColumn("Digital Signal Bandwidth", -1)
					} else {
						childListTable.AddTableColumn(name, -1)
					}
				}
			}
		}
	}

	childViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {

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
			} else if reCarrierGroup.FindString(gp.Name) != "" {
				var carrierIndex int64 = -1
				if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
					carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if carrierIndex < 1 {
					continue
				}
				var (
					carrierDigitalSignalBandwidthSelect *Element
					carrierDigitalSignalBandwidthCustom *Element
				)
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)

					if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Select") {
						carrierDigitalSignalBandwidthSelect = elem
					} else if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Custom") {
						carrierDigitalSignalBandwidthCustom = elem
					} else {
						rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
						childListTable.MustSetTableRowData(rowIndex, name, elem)
					}
					key := fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)
					viewItems, ok := childViewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					childViewItemsMap[key] = viewItems
				}
				if carrierDigitalSignalBandwidthCustom != nil && carrierDigitalSignalBandwidthSelect != nil {
					rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
					elem2 := NewParamGroupComponent("Digital Signal Bandwidth", 2, carrierDigitalSignalBandwidthSelect, carrierDigitalSignalBandwidthCustom)
					childListTable.MustSetTableRowData(rowIndex, elem2.Name, elem2)
				}
			}
		}
	}
	for k, items := range childViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		carrierIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1

		name := fmt.Sprintf("Radio Module %v Carrier %v", moduleIndex, carrierIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		childListTable.MustSetTableRowData(int64(rowIndex), "Carrier",
			NewLabel("Carrier", fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)))
		childListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	moduleListTable.SetStyle("height", "180px")

	module := NewModule(m.Name)
	tabsLayout := NewTabsLayout("System Signal", "top")
	tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Carrier List",
		NewColLayoutWithItems(0, moduleListTable, childListTable),
	))
	tabsLayout.Items = append(tabsLayout.Items, tabPages...)

	module.Items = append(module.Items, NewPageWithLayouts("System Signal", tabsLayout))
	return module, nil
}

func MakeModuleForRuCarrierConfigUseAmplifierModule(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Amplifier\s*(\d+)`)
	reModuleGroup := regexp.MustCompile(`Amplifier\s*(\d+|)\s*Configuration`)
	reCarrierGroup := regexp.MustCompile(`Channel\s*(\d+)`)

	reModule := regexp.MustCompile(`Amplifier Module\s*(\d+)`)
	reModuleCarrier := regexp.MustCompile(`Amplifier Module\s*(\d+) Channel\s*(\d)`)

	moduleListTable := NewTable("Amplifier Module List")
	moduleListTable.AddTableColumnUnique("Module", -1)

	childListTable := NewTable("Channel List")
	childListTable.AddTableColumnUnique("Channel", -1)

	tabPages := []*Element{}
	var childNum int64 = 0
	for _, pg := range m.Child {
		if reModulePage.FindString(pg.Name) != "" {
			for _, gp := range pg.Child {
				if reModuleGroup.FindString(gp.Name) != "" {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
						moduleListTable.AddTableColumn(name, -1)
					}
				} else if reCarrierGroup.FindString(gp.Name) != "" {
					var index int64 = -1
					if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
						if index, _ = strconv.ParseInt(subs[1], 10, 32); index > childNum {
							childNum = index
						}
					}
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
						if strings.HasPrefix(name, "Digital Signal Bandwidth") {
							childListTable.AddTableColumn("Digital Signal Bandwidth", -1)
						} else {
							childListTable.AddTableColumn(name, -1)
						}
					}
				}
			}
		}

	}

	childViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {

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
			} else if reCarrierGroup.FindString(gp.Name) != "" {
				var carrierIndex int64 = -1
				if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
					carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if carrierIndex < 1 {
					continue
				}
				var (
					carrierDigitalSignalBandwidthSelect *Element
					carrierDigitalSignalBandwidthCustom *Element
				)
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)

					if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Select") {
						carrierDigitalSignalBandwidthSelect = elem
					} else if strings.HasSuffix(elem.Name, "Digital Signal Bandwidth Custom") {
						carrierDigitalSignalBandwidthCustom = elem
					} else {
						rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
						childListTable.MustSetTableRowData(rowIndex, name, elem)
					}
					key := fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)
					viewItems, ok := childViewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					childViewItemsMap[key] = viewItems
				}
				if carrierDigitalSignalBandwidthCustom != nil && carrierDigitalSignalBandwidthSelect != nil {
					rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1
					elem2 := NewParamGroupComponent("Digital Signal Bandwidth", 2, carrierDigitalSignalBandwidthSelect, carrierDigitalSignalBandwidthCustom)
					childListTable.MustSetTableRowData(rowIndex, elem2.Name, elem2)
				}
			}
		}
	}
	for k, items := range childViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		carrierIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*childNum + carrierIndex - 1

		name := fmt.Sprintf("Amplifier Module %v Channel %v", moduleIndex, carrierIndex)
		viewForm := NewSetParameterValuesForm(name, items...)
		viewPageItems := []*Element{viewForm}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		childListTable.MustSetTableRowData(int64(rowIndex), "Channel",
			NewLabel("Channel", fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)))
		childListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	moduleListTable.SetStyle("height", "180px")

	module := NewModule(m.Name)
	tabsLayout := NewTabsLayout("System Signal", "top")
	tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Carrier List",
		NewColLayoutWithItems(0, moduleListTable, childListTable),
	))
	tabsLayout.Items = append(tabsLayout.Items, tabPages...)

	module.Items = append(module.Items, NewPageWithLayouts("System Signal", tabsLayout))
	return module, nil
}
