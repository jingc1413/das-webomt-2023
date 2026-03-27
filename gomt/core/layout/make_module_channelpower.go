package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakeModuleForChannelPower(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reChannelGroup := regexp.MustCompile(`Channel\s*(\d+)`)
	reModuleChannel := regexp.MustCompile(`Module\s*(\d+) Channel(\d+)`)

	module := NewModule(m.Name)

	channelListTable := NewTable("Channel List")
	channelListTable.AddTableColumnUnique("Channel", -1)

	var channelNum int64 = 0
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			var channelIndex int64 = -1
			if subs := reChannelGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				if channelIndex, _ = strconv.ParseInt(subs[1], 10, 32); channelIndex > channelNum {
					channelNum = channelIndex
				}
			}
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModuleChannel.ReplaceAllString(obj.Name, ""))
				channelListTable.AddTableColumn(name, -1)
			}
		}
	}

	channelViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		var moduleIndex int64 = -1
		if subs := reModulePage.FindStringSubmatch(pg.Name); len(subs) == 2 {
			moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if moduleIndex < 1 {
			continue
		}
		for _, gp := range pg.Child {
			var channelIndex int64 = -1
			if subs := reChannelGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				channelIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if channelIndex < 1 {
				continue
			}
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)
			rowIndex := (moduleIndex-1)*channelNum + channelIndex - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModuleChannel.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				channelListTable.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v-%v", moduleIndex, channelIndex)
				viewItems, ok := channelViewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				channelViewItemsMap[key] = viewItems
			}
		}
	}
	for k, items := range channelViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		channelIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*channelNum + channelIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Module %v Channel %v", moduleIndex, channelIndex),
			items...)
		channelListTable.MustSetTableRowData(int64(rowIndex), "Channel",
			NewLabel("Channel", fmt.Sprintf("%v-%v", moduleIndex, channelIndex)))
		channelListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	moduleConfigurationPage := NewPageWithLayouts("Power Configuration",
		NewSingleColLayoutWithItems(channelListTable),
	)
	module.Items = append(module.Items, moduleConfigurationPage)
	return module, nil
}
