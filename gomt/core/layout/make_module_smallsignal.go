package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func MakeModuleForSmallSignal(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	module := NewModule(m.Name)
	for _, pg := range m.Child {
		if pg.Name == "Frequency" {
			page, err := makePageForSmallSignalFrequency(pg, info, path)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
		} else if pg.Name == "Gain" {
			page, err := makePageForSmallSignalGain(pg, info, path)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
		} else if pg.Name == "Other Parameters" {
			pages, err := makePageForOtherParameters(pg, info, path)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			for _, page := range pages {
				addPageToModule(page, module)
			}
			// module.Items = append(module.Items, pages...)
		} else {
			page, err := MakePage(m, pg, info)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
		}
	}

	return module, nil
}

func makePageForSmallSignalFrequency(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	re := regexp.MustCompile(`Radio\s*Module\s*(\d+)`)

	listTable := NewTable("Frequency")
	listTable.AddTableColumnUnique("Radio Module", -1)

	for _, gp := range pg.Child {
		if re.FindString(gp.Name) != "" {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(re.ReplaceAllString(obj.Name, ""))
				listTable.AddTableColumn(name, -1)
			}
		}
	}
	viewItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		var index int64 = -1
		if subs := re.FindStringSubmatch(gp.Name); len(subs) == 2 {
			index, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if index < 1 {
			continue
		}
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		rowIndex := index - 1
		listTable.MustSetTableRowData(rowIndex, "Radio Module", NewLabel("Radio Module", fmt.Sprintf("%v", index)))
		for _, obj := range gp.Child {
			name := strings.TrimSpace(re.ReplaceAllString(obj.Name, ""))
			elem := MustMakeElementFromParam(obj, info, groupPath)
			elem.SetName(name)
			for _, item := range elem.Items {
				item.SetName(strings.TrimSpace(re.ReplaceAllString(item.Name, "")))
			}

			listTable.MustSetTableRowData(rowIndex, name, elem)
			key := fmt.Sprintf("%v", index)
			viewItems, ok := viewItemsMap[key]
			if !ok {
				viewItems = []*Element{}
			}
			viewItems = append(viewItems, elem)
			viewItemsMap[key] = viewItems
		}
	}
	for k, items := range viewItemsMap {
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v", index), items...)
		listTable.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		listTable.MustSetTableRowData(rowIndex, "Radio Module", NewLabel("Radio Module", fmt.Sprintf("%v", index)))
	}

	pages := []*Element{}
	for _, gp := range pg.Child {
		if re.FindString(gp.Name) == "" {
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)

			items := []*Element{}
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				items = append(items, elem)
			}
			pages = append(pages, NewPageWithSetParameterValuesFormItems(gp.Name, items...))
		}
	}
	if len(listTable.Items) > 1 {
		pages = append(pages, NewPageWithLayouts("Frequency",
			NewSingleColLayoutWithItems(listTable),
		))
	}
	if len(pages) == 1 {
		return pages[0], nil
	}
	tabsLayout := NewTabsLayout("Frequency", "left")
	tabsLayout.Items = append(tabsLayout.Items, pages...)
	return NewPageWithLayouts("Frequency", tabsLayout), nil
}

func makePageForSmallSignalGain(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModuleGroup := regexp.MustCompile(`Radio\s*Module\s*(\d+)`)
	reModule := regexp.MustCompile(`Radio Module\s*(\d+)`)

	listTable := NewTable("Gain")
	listTable.AddTableColumnUnique("Radio Module", -1)

	for _, gp := range pg.Child {
		if gp.Name == "UpLink DSA Attention" || gp.Name == "DownLink VOP Attention" {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				listTable.AddTableColumn(name, -1)
			}
		} else {
			if reModuleGroup.FindString(gp.Name) != "" {
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
					listTable.AddTableColumn(name, -1)
				}
			}
		}
	}

	viewItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		if gp.Name == "UpLink DSA Attention" || gp.Name == "DownLink VOP Attention" {
			for _, obj := range gp.Child {
				var index int64 = -1
				if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 2 {
					index, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if index < 1 {
					continue
				}
				rowIndex := index - 1
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				listTable.MustSetTableRowData(rowIndex, name, elem)
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

				listTable.MustSetTableRowData(rowIndex, name, elem)
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
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v Gain", index), items...)
		listTable.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		listTable.MustSetTableRowData(rowIndex, "Radio Module", NewLabel("Radio Module", fmt.Sprintf("%v", index)))
	}
	return NewPageWithLayouts("Gain",
		NewSingleColLayoutWithItems(listTable),
	), nil
}

func makePageForOtherParameters(pg *model.ParameterTreeNode, info DeviceInfo, path []string) ([]*Element, error) {
	reModuleFeedbackGroup := regexp.MustCompile(`Radio\s*Module\s*(\d+)\s*Feedback`)
	reModuleTemperatureCompensationGroup := regexp.MustCompile(`Radio\s*Module\s*(\d+)\s*Temperature Compensation`)

	reModule := regexp.MustCompile(`Radio\s*Module\s*(\d+)`)

	pages := []*Element{}

	listTable := NewTable("Feedback")
	listTable.AddTableColumnUnique("Radio Module", -1)

	listTable2 := NewTable("Temperature Compensation")
	listTable2.AddTableColumnUnique("Radio Module", -1)

	for _, gp := range pg.Child {
		if reModuleFeedbackGroup.FindString(gp.Name) != "" {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				listTable.AddTableColumn(name, -1)
			}
		} else if reModuleTemperatureCompensationGroup.FindString(gp.Name) != "" {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				listTable2.AddTableColumn(name, -1)
			}
		} else {
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)
			items := []*Element{}
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				items = append(items, elem)
			}
			if len(items) > 0 {
				pages = append(pages, NewPageWithSetParameterValuesFormItems(gp.Name, items...))
			}
		}
	}

	viewItemsMap := map[string][]*Element{}
	viewItemsMap2 := map[string][]*Element{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		if reModuleFeedbackGroup.FindString(gp.Name) != "" {
			var index int64 = -1
			if subs := reModuleFeedbackGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
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

				listTable.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap[key] = viewItems
			}
		} else if reModuleTemperatureCompensationGroup.FindString(gp.Name) != "" {
			var index int64 = -1
			if subs := reModule.FindStringSubmatch(gp.Name); len(subs) == 2 {
				index, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if index < 1 {
				continue
			}
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)
			rowIndex := index - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModule.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				listTable2.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v", index)
				viewItems, ok := viewItemsMap2[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				viewItemsMap2[key] = viewItems
			}
		}
	}
	for k, items := range viewItemsMap {
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v Feedback", index), items...)
		listTable.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		listTable.MustSetTableRowData(rowIndex, "Radio Module", NewLabel("Radio Module", fmt.Sprintf("%v", index)))
	}

	for k, items := range viewItemsMap2 {
		index := int64(0)
		index, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Radio Module %v Temperature Compensation", index), items...)
		listTable2.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		listTable2.MustSetTableRowData(rowIndex, "Radio Module", NewLabel("Radio Module", fmt.Sprintf("%v", index)))
	}

	if len(listTable.Items) > 1 {
		pages = append(pages, NewPageWithLayouts("Feedback", NewSingleColLayoutWithItems(listTable)))
	}
	if len(listTable2.Items) > 1 {
		pages = append(pages, NewPageWithLayouts("Temperature Compensation", NewSingleColLayoutWithItems(listTable2)))
	}
	return pages, nil
}
