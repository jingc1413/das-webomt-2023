package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func MakeModuleForPA(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	rePA := regexp.MustCompile(`^PA\s*(CH|)\s*(\d+)$`)

	module := NewModule(m.Name)
	if page, err := makePageForPATabsPage(m, info, path); err != nil {
		return nil, errors.Wrap(err, "make page for pa list")
	} else if page != nil {
		addPageToModule(page, module)
	}

	for _, pg := range m.Child {
		if rePA.MatchString(pg.Name) {
			continue
		}
		page, err := MakePageWithExcludes(m, pg, info, []string{"RX Switch", "TX Switch", "PA InitFile"})
		if err != nil {
			return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
		}
		addPageToModule(page, module)
	}

	return module, nil
}

func makePageForPATabsPage(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	rePA := regexp.MustCompile(`PA\s*(CH|)\s*(\d+)`)
	reRx := regexp.MustCompile(`RX\s*Switch\s*(\d+)`)
	reTx := regexp.MustCompile(`TX\s*Switch\s*(\d+)`)

	viewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		if pg.Name == "PA SWITCH" {
			for _, gp := range pg.Child {
				var re *regexp.Regexp
				if gp.Name == "RX Switch" {
					re = reRx
				} else if gp.Name == "TX Switch" {
					re = reTx
				}
				if re == nil {
					continue
				}
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				for _, obj := range gp.Child {
					index := int64(0)
					if subs := re.FindStringSubmatch(obj.Name); len(subs) == 2 {
						index, _ = strconv.ParseInt(subs[1], 10, 32)
					}
					if index < 1 {
						continue
					}

					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(gp.Name)
					key := fmt.Sprintf("%v", index)
					viewItems, ok := viewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					viewItemsMap[key] = viewItems
				}
			}
		} else {
			var paIndex int64 = -1

			if subs := rePA.FindStringSubmatch(pg.Name); len(subs) == 3 {
				paIndex, _ = strconv.ParseInt(subs[2], 10, 32)
			}
			if paIndex < 1 {
				continue
			}
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				for _, obj := range gp.Child {
					name := strings.TrimSpace(rePA.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)
					for _, item := range elem.Items {
						item.SetName(strings.TrimSpace(rePA.ReplaceAllString(item.Name, "")))
					}
					key := fmt.Sprintf("%v", paIndex)
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
	tabs := NewTabsLayout("", "left")
	for k, items := range viewItemsMap {
		paIndex := int64(0)
		paIndex, _ = strconv.ParseInt(k, 10, 32)
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("PA %v", paIndex), items...)
		tabs.Items = append(tabs.Items, viewPage)
	}
	if len(tabs.Items) > 0 {
		return NewPageWithLayouts("PA Configuration", tabs), nil
	}
	return nil, nil
}

func makePageForPAList(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	rePA := regexp.MustCompile(`PA\s*(CH|)\s*(\d+)`)
	reRx := regexp.MustCompile(`RX\s*Switch\s*(\d+)`)
	reTx := regexp.MustCompile(`TX\s*Switch\s*(\d+)`)

	listTable := NewTable("PA List")
	listTable.AddTableColumnUnique("PA", -1)

	for _, pg := range m.Child {
		if rePA.FindString(pg.Name) != "" {
			for _, gp := range pg.Child {
				for _, obj := range gp.Child {
					name := strings.TrimSpace(rePA.ReplaceAllString(obj.Name, ""))
					listTable.AddTableColumn(name, -1)
				}
			}
		} else if pg.Name == "PA SWITCH" {
			for _, gp := range pg.Child {
				if gp.Name == "RX Switch" {
					listTable.AddTableColumn("RX Switch", -1)
				}
				if gp.Name == "TX Switch" {
					listTable.AddTableColumn("TX Switch", -1)
				}
			}
		}
	}
	viewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		if pg.Name == "PA SWITCH" {
			for _, gp := range pg.Child {
				var re *regexp.Regexp
				if gp.Name == "RX Switch" {
					re = reRx
				} else if gp.Name == "TX Switch" {
					re = reTx
				}
				if re == nil {
					continue
				}
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				for _, obj := range gp.Child {
					index := int64(0)
					if subs := re.FindStringSubmatch(obj.Name); len(subs) == 2 {
						index, _ = strconv.ParseInt(subs[1], 10, 32)
					}
					if index < 1 {
						continue
					}
					rowIndex := index - 1

					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(gp.Name)
					listTable.MustSetTableRowData(rowIndex, gp.Name, elem)
					key := fmt.Sprintf("%v", index)
					viewItems, ok := viewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					viewItemsMap[key] = viewItems
				}
			}
		} else {
			var paIndex int64 = -1

			if subs := rePA.FindStringSubmatch(pg.Name); len(subs) == 3 {
				paIndex, _ = strconv.ParseInt(subs[2], 10, 32)
			}
			if paIndex < 1 {
				continue
			}
			rowIndex := paIndex - 1
			for _, gp := range pg.Child {
				groupPath := append([]string{}, path...)
				groupPath = append(groupPath, pg.Name, gp.Name)
				for _, obj := range gp.Child {
					name := strings.TrimSpace(rePA.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)
					for _, item := range elem.Items {
						item.SetName(strings.TrimSpace(rePA.ReplaceAllString(item.Name, "")))
					}

					listTable.MustSetTableRowData(rowIndex, name, elem)
					key := fmt.Sprintf("%v", paIndex)
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
		paIndex := int64(0)
		paIndex, _ = strconv.ParseInt(k, 10, 32)
		rowIndex := paIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("PA %v", paIndex), items...)
		listTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
		listTable.MustSetTableRowData(int64(rowIndex), "PA", NewLabel("PA", fmt.Sprintf("%v", paIndex)))
	}
	if listTable.Data != nil && len(listTable.Data) > 0 {
		return NewPageWithLayouts("PA Configuration",
			NewSingleColLayoutWithItems(listTable),
		), nil
	}
	return nil, nil
}
