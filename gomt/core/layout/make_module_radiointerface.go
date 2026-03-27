package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakePageForAuRadioInterfaceModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	if info.Device.Schema == "corning" {
		return MakePageForAuRadioInterfaceModuleUseInputModule(pg, info, path)
	} else {
		return MakePageForAuRadioInterfaceModuleUseRadioModule(pg, info, path)
	}
}

func MakePageForAuRadioInterfaceModuleUseRadioModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Radio Module\s*(\d+)`)

	reModulePort := regexp.MustCompile(`Radio Module\s*(\d+) Port\s*(\d+)`)
	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)

	generalFormItems := []*Element{}

	table := NewTable("Radio Module List")
	table.AddTableColumnUnique("Module", -1)
	table.AddTableColumn("Serial Number", -1)

	var portNum int64 = 0
	table2 := NewTable("Radio Interface List")
	table2.AddTableColumnUnique("Interface", -1)

	for _, gp := range pg.Child {
		if gp.Name == "General" {
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				generalFormItems = append(generalFormItems, elem)
			}
		} else if reModuleGroup.FindString(gp.Name) != "" {
			for _, obj := range gp.Child {
				if subs := reModulePort.FindStringSubmatch(obj.Name); len(subs) == 3 {
					if v, _ := strconv.ParseInt(subs[2], 10, 32); v > portNum {
						portNum = v
					}
					name := strings.TrimSpace(reModulePort.ReplaceAllString(obj.Name, "Port"))
					table2.AddTableColumn(name, -1)
				} else if reModule.FindString(obj.Name) != "" {
					name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
					table.AddTableColumn(name, -1)
				}
			}
		}
	}

	viewFormItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		var moduleIndex int64 = -1
		if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
			moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if moduleIndex < 1 {
			continue
		}
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		for _, obj := range gp.Child {
			if reModulePort.FindString(obj.Name) != "" {
				var portIndex int64 = -1
				if subs := reModulePort.FindStringSubmatch(obj.Name); len(subs) == 3 {
					portIndex, _ = strconv.ParseInt(subs[2], 10, 32)
				}
				if portIndex < 1 {
					continue
				}
				rowIndex := (moduleIndex-1)*portNum + portIndex - 1
				name := strings.TrimSpace(reModulePort.ReplaceAllString(obj.Name, "Port"))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				key := fmt.Sprintf("%v-%v", moduleIndex, portIndex)
				viewFormItems, ok := viewFormItemsMap[key]
				if !ok {
					viewFormItems = []*Element{}
				}
				table2.AddTableColumn(name, -1)
				table2.MustSetTableRowData(rowIndex, name, elem)
				viewFormItems = append(viewFormItems, elem)
				viewFormItemsMap[key] = viewFormItems
			} else if reModule.FindString(obj.Name) != "" {
				rowIndex := moduleIndex - 1
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}

				table.MustSetTableRowData(rowIndex, name, elem)
				table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
			}
		}
	}

	for k, items := range viewFormItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		portIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex2 := (moduleIndex-1)*portNum + portIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v Port %v", moduleIndex, portIndex), items...)
		table2.SetTableRowClickAction(rowIndex2, NewViewPageAction(viewPage))
		table2.MustSetTableRowData(rowIndex2, "Interface", NewLabel("Interface", fmt.Sprintf("%v - %v", moduleIndex, portIndex)))
	}

	table.SetStyle("height", "180px")

	tabs := NewTabsLayout(pg.Name, "top")
	tabs.Items = append(tabs.Items,
		NewPageWithLayouts("Radio Interface", NewColLayoutWithItems(0, table, table2)),
		NewPageWithSetParameterValuesFormItems("General", generalFormItems...),
	)
	page := NewPageWithLayouts(pg.Name, tabs)
	return page, nil
}

func MakePageForAuRadioInterfaceModuleUseInputModule(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Input Module\s*(\d+)`)

	reModulePort := regexp.MustCompile(`Input Module\s*(\d+) Port\s*(\d+)`)
	reModule := regexp.MustCompile(`Input Module\s*(\d+)`)

	generalFormItems := []*Element{}

	table := NewTable("Input Module List")
	table.AddTableColumnUnique("Module", -1)
	table.AddTableColumn("Serial Number", -1)

	var portNum int64 = 0
	table2 := NewTable("Input Interface List")
	table2.AddTableColumnUnique("Interface", -1)

	for _, obj := range pg.Child[1].Child {
		if subs := reModulePort.FindStringSubmatch(obj.Name); len(subs) == 3 {
			if v, _ := strconv.ParseInt(subs[2], 10, 32); v > portNum {
				portNum = v
			}
			name := strings.TrimSpace(reModulePort.ReplaceAllString(obj.Name, "Port"))
			table2.AddTableColumn(name, -1)
		} else if reModule.FindString(obj.Name) != "" {
			name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
			table.AddTableColumn(name, -1)
		}
	}

	for _, gp := range pg.Child {
		if gp.Name == "General" {
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				generalFormItems = append(generalFormItems, elem)
			}
		}
	}

	viewFormItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		var moduleIndex int64 = -1
		if subs := reModuleGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
			moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if moduleIndex < 1 {
			continue
		}
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		for _, obj := range gp.Child {
			if reModulePort.FindString(obj.Name) != "" {
				var portIndex int64 = -1
				if subs := reModulePort.FindStringSubmatch(obj.Name); len(subs) == 3 {
					portIndex, _ = strconv.ParseInt(subs[2], 10, 32)
				}
				if portIndex < 1 {
					continue
				}
				rowIndex := (moduleIndex-1)*portNum + portIndex - 1
				name := strings.TrimSpace(reModulePort.ReplaceAllString(obj.Name, "Port"))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				for _, item := range elem.Items {
					n := strings.TrimSpace(reModule.ReplaceAllString(item.Name, ""))
					item.SetName(n)
				}
				key := fmt.Sprintf("%v-%v", moduleIndex, portIndex)
				viewFormItems, ok := viewFormItemsMap[key]
				if !ok {
					viewFormItems = []*Element{}
				}
				table2.AddTableColumn(name, -1)
				table2.MustSetTableRowData(rowIndex, name, elem)
				viewFormItems = append(viewFormItems, elem)
				viewFormItemsMap[key] = viewFormItems
			} else if reModule.FindString(obj.Name) != "" {
				rowIndex := moduleIndex - 1
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				table.MustSetTableRowData(rowIndex, name, elem)
				table.MustSetTableRowData(rowIndex, "Module", NewLabel("Module", fmt.Sprintf("%v", moduleIndex)))
			}
		}
	}

	for k, items := range viewFormItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		portIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex2 := (moduleIndex-1)*portNum + portIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Input Module %v Port %v", moduleIndex, portIndex), items...)
		table2.SetTableRowClickAction(rowIndex2, NewViewPageAction(viewPage))
		table2.MustSetTableRowData(rowIndex2, "Interface", NewLabel("Interface", fmt.Sprintf("%v - %v", moduleIndex, portIndex)))
	}

	table.SetStyle("height", "180px")

	tabs := NewTabsLayout(pg.Name, "top")
	tabs.Items = append(tabs.Items,
		NewPageWithLayouts("Input Interface", NewColLayoutWithItems(0, table, table2)),
		NewPageWithSetParameterValuesFormItems("General", generalFormItems...),
	)
	page := NewPageWithLayouts(pg.Name, tabs)
	return page, nil
}
