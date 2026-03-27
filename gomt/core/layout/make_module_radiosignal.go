package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakePageForAuRadioSignalInformation(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakePageForAuRadioSignalInformationUseInputModule(pg, info, path)
	} else {
		return MakePageForAuRadioSignalInformationUseRadioModule(pg, info, path)
	}
}

func MakePageForRuRadioSignalInformation(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakePageForRuRadioSignalInformationUseAmplifierModule(pg, info, path)
	} else {
		return MakePageForRuRadioSignalInformationUseRadioModule(pg, info, path)
	}
}

func MakePageForAuRadioSignalInformationUseRadioModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Radio Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)

	table := NewTable("Radio Module List")
	table.AddTableColumnUnique("Module", -1)
	for _, gp := range pg.Child {
		if reModuleGroup.FindString(gp.Name) != "" {
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

		if reModuleGroup.FindString(gp.Name) != "" {
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
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v", index), items...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", index)))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}

func MakePageForRuRadioSignalInformationUseRadioModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Radio Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)

	table := NewTable("Radio Module List")
	table.AddTableColumnUnique("Module", -1)
	for _, gp := range pg.Child {
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			table.AddTableColumn(gp.Name, -1)
		} else if reModuleGroup.FindString(gp.Name) != "" {
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
				name := gp.Name
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				table.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		} else if reModuleGroup.FindString(gp.Name) != "" {
			var index int64 = -1
			if subs := reModule.FindStringSubmatch(gp.Name); len(subs) == 2 {
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
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v", index), items...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", index)))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}

func MakePageForAuRadioSignalInformationUseInputModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Input Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Input Module\s*(\d+)`)

	table := NewTable("Input Module List")
	table.AddTableColumnUnique("Module", -1)
	for _, gp := range pg.Child {
		if reModuleGroup.FindString(gp.Name) != "" {
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

		if reModuleGroup.FindString(gp.Name) != "" {
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
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Input Module %v", index), items...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", index)))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}

func MakePageForRuRadioSignalInformationUseAmplifierModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Amplifier Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Amplifier Module\s*(\d+)`)

	table := NewTable("Amplifier Module List")
	table.AddTableColumnUnique("Module", -1)
	for _, gp := range pg.Child {
		if gp.Name == "Uplink Baseband Power" || gp.Name == "Downlink Baseband Power" {
			table.AddTableColumn(gp.Name, -1)
		} else if reModuleGroup.FindString(gp.Name) != "" {
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
				name := gp.Name
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				table.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		} else if reModuleGroup.FindString(gp.Name) != "" {
			var index int64 = -1
			if subs := reModule.FindStringSubmatch(gp.Name); len(subs) == 2 {
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
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Amplifier Module %v", index), items...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", index)))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}
