package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakePageForRuBandConfiguration(pg *model.ParameterTreeNode, pg2 *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakePageForRuBandConfigurationUseAmplifierModule(pg, pg2, info, path)
	} else {
		return MakePageForRuBandConfigurationUseRadioModule(pg, pg2, info, path)
	}
}
func MakePageForRuBandConfigurationUseRadioModule(pg *model.ParameterTreeNode, pg2 *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Radio\s*Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Radio\s*Module\s*(\d+)`)

	table := NewTable("Radio Module List")
	table.AddTableColumnUnique("Module", -1)

	for _, gp := range pg.Child {
		if strings.HasPrefix(gp.Name, "Module Mapping") {
			table.AddTableColumn("Module Mapping", -1)
		}
	}
	for _, gp := range pg.Child {
		if !strings.HasPrefix(gp.Name, "Module Mapping") {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				table.AddTableColumn(name, -1)
			}
		}
	}

	for _, gp := range pg2.Child {
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			table.AddTableColumn(gp.Name, -1)
		} else {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				table.AddTableColumn(name, -1)
			}
		}
	}

	viewItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		if strings.HasPrefix(gp.Name, "Module Mapping") {
			for _, obj := range gp.Child {
				var moduleIndex int64 = -1
				if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
					moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if moduleIndex < 1 {
					continue
				}
				rowIndex := moduleIndex - 1
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				if strings.HasSuffix(name, "Module Mapping") {
					table.MustSetTableRowData(rowIndex, "Module Mapping", elem)
					key := fmt.Sprintf("%v", moduleIndex)
					viewItems, ok := viewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					viewItemsMap[key] = viewItems
				}
			}
		}
	}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		if !strings.HasPrefix(gp.Name, "Module Mapping") {
			var moduleIndex int64 = -1
			if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if moduleIndex < 1 {
				continue
			}
			rowIndex := moduleIndex - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}

				table.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", moduleIndex)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				if name == "RF Signal Active" {
					elem2 := *elem
					elem2.SetStyle("input", "button")
					action := NewAction("StateAction")
					action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(&elem2, "01")))
					action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(&elem2, "00")))

					elem2.SetAction("click", action)
					elem2.SetStyle("input", "button")
					elem2.SetStyle("inactiveValue", "00")
					elem2.SetStyle("inactiveText", "Signal Inactive")
					elem2.SetStyle("inactiveType", "info")
					elem2.SetStyle("activeValue", "01")
					elem2.SetStyle("activeText", "Signal Active")
					elem2.SetStyle("activeType", "success")
					viewItems = append(viewItems, &elem2)
				} else {
					viewItems = append(viewItems, elem)
				}
				viewItemsMap[key] = viewItems
			}
		}
	}

	for _, gp := range pg2.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg2.Name, gp.Name)
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			for _, obj := range gp.Child {
				var index int64 = -1
				if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
					index, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if index < 1 {
					continue
				}
				rowIndex := index - 1
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(gp.Name)
				table.MustSetTableRowData(rowIndex, gp.Name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		} else {
			var index int64 = -1
			if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				index, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if index < 1 {
				continue
			}
			rowIndex := index - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}
				if strings.HasSuffix(name, "Module Mapping") {
					table.MustSetTableRowData(rowIndex, name, elem)
					key := fmt.Sprintf("%v", index)
					viewItems, ok := viewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					viewItemsMap[key] = viewItems
				}
			}
		}
	}
	for k, items := range viewItemsMap {
		moduleIndex, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := moduleIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Radio Module %v", moduleIndex),
			items...)
		table.MustSetTableRowData(int64(rowIndex), "Module", NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
		table.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}

func MakePageForRuBandConfigurationUseAmplifierModule(pg *model.ParameterTreeNode, pg2 *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Amplifier\s*Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Amplifier\s*Module\s*(\d+)`)

	table := NewTable("Amplifier Module List")
	table.AddTableColumnUnique("Module", -1)

	for _, gp := range pg.Child {
		if strings.HasPrefix(gp.Name, "Module Mapping") {
			table.AddTableColumn("Module Mapping", -1)
		}
	}
	for _, gp := range pg.Child {
		if !strings.HasPrefix(gp.Name, "Module Mapping") {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				table.AddTableColumn(name, -1)
			}
		}
	}

	for _, gp := range pg2.Child {
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			table.AddTableColumn(gp.Name, -1)
		} else {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				table.AddTableColumn(name, -1)
			}
		}
	}

	viewItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		if strings.HasPrefix(gp.Name, "Module Mapping") {
			for _, obj := range gp.Child {
				var moduleIndex int64 = -1
				if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
					moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if moduleIndex < 1 {
					continue
				}
				rowIndex := moduleIndex - 1
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				if strings.HasSuffix(name, "Module Mapping") {
					table.MustSetTableRowData(rowIndex, "Module Mapping", elem)
					key := fmt.Sprintf("%v", moduleIndex)
					viewItems, ok := viewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					viewItemsMap[key] = viewItems
				}
			}
		}
	}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		if !strings.HasPrefix(gp.Name, "Module Mapping") {
			var moduleIndex int64 = -1
			if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if moduleIndex < 1 {
				continue
			}
			rowIndex := moduleIndex - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}

				table.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", moduleIndex)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				if name == "RF Signal Active" {
					elem2 := MustMakeElementFromParam(obj, info, groupPath)
					elem2.SetName(name)
					elem2.SetStyle("input", "button")
					action := NewAction("StateAction")
					action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(elem2, "01")))
					action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(elem2, "00")))

					elem2.SetAction("click", action)
					elem2.SetStyle("input", "button")
					elem2.SetStyle("inactiveValue", "00")
					elem2.SetStyle("inactiveText", "Signal Inactive")
					elem2.SetStyle("inactiveType", "info")
					elem2.SetStyle("activeValue", "01")
					elem2.SetStyle("activeText", "Signal Active")
					elem2.SetStyle("activeType", "success")
					viewItems = append(viewItems, elem2)
				} else {
					viewItems = append(viewItems, elem)
				}
				viewItemsMap[key] = viewItems
			}
		}
	}

	for _, gp := range pg2.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg2.Name, gp.Name)
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			for _, obj := range gp.Child {
				var index int64 = -1
				if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
					index, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if index < 1 {
					continue
				}
				rowIndex := index - 1
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(gp.Name)
				table.MustSetTableRowData(rowIndex, gp.Name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		} else {
			var index int64 = -1
			if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				index, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if index < 1 {
				continue
			}
			rowIndex := index - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}

				table.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		}
	}
	for k, items := range viewItemsMap {
		moduleIndex, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := moduleIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Amplifier Module %v", moduleIndex),
			items...)
		table.MustSetTableRowData(int64(rowIndex), "Module", NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
		table.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}
