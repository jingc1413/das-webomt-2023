package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func getOpticalIndex(indexText string) int64 {
	var index int64
	if indexText == "P" || indexText == "M" {
		index = 1
	} else if indexText == "S" {
		index = 2
	} else {
		index = cast.ToInt64(indexText)
	}
	return index
}

func MakePageForOpticalModuleInformation(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	table := NewTable("Optical Module List")
	table.AddTableColumnUnique("Optical Module", -1)
	for _, gp := range pg.Child {
		if gp.Name == "Optical Module Serial Number And Vendor Name" {
			table.AddTableColumn("Serial Number", -1)
			table.AddTableColumn("Vendor Name", -1)
		} else if gp.Name == "Optical Module Tx Power And Rx Power" {
			table.AddTableColumn("Tx Power", -1)
			table.AddTableColumn("Rx Power", -1)
		}
	}

	re := regexp.MustCompile(`OP\s*(\d+|P|S|M)`)
	viewFormItemsMap := map[string][]any{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		for _, obj := range gp.Child {
			indexText := ""
			if subs := re.FindStringSubmatch(obj.Name); len(subs) == 2 {
				indexText = subs[1]
			}
			index := getOpticalIndex(indexText)
			if index <= 0 {
				continue
			}

			rowIndex := index - 1
			elem := MustMakeElementFromParam(obj, info, groupPath)
			if gp.Name == "Optical Module Serial Number And Vendor Name" {
				if strings.Contains(obj.Name, "Serial Number") {
					table.MustSetTableRowData(rowIndex, "Serial Number", elem)
				} else if strings.Contains(obj.Name, "Vendor Name") {
					table.MustSetTableRowData(rowIndex, "Vendor Name", elem)
				}
			} else if gp.Name == "Optical Module Tx Power And Rx Power" {
				if strings.Contains(obj.Name, "Tx Power") {
					table.MustSetTableRowData(rowIndex, "Tx Power", elem)
				} else if strings.Contains(obj.Name, "Rx Power") {
					table.MustSetTableRowData(rowIndex, "Rx Power", elem)
				}
			}

			viewFormItems, ok := viewFormItemsMap[indexText]
			if !ok {
				viewFormItems = []any{}
			}
			viewFormItems = append(viewFormItems, elem)
			viewFormItemsMap[indexText] = viewFormItems
		}

	}
	for indexText, _ := range viewFormItemsMap {
		index := getOpticalIndex(indexText)
		rowIndex := index - 1
		table.MustSetTableRowData(rowIndex, "Optical Module", NewLabel("Optical Module", fmt.Sprintf("OP %v", indexText)))
	}

	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}

func MakePageForFactoryCommand(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	var factoryMode *Element
	var factoryCommand *Element

	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		if gp.Name == "Factory Parameters" {
			for _, obj := range gp.Child {
				if obj.Name == "Factory Mode" {
					factoryMode = MustMakeElementFromParam(obj, info, groupPath)
				} else if obj.Name == "Factory Mode Password" || obj.Name == "Factory Command" {
					factoryCommand = MustMakeElementFromParam(obj, info, groupPath)
					factoryCommand.SetName("Factory Command")
				}
			}
		}
	}
	pages := []*Element{}
	if factoryMode != nil {
		pages = append(pages, NewPageWithSetParameterValuesFormItems("Factory Mode", factoryMode))
	}
	if factoryCommand != nil {
		pages = append(pages, NewPageWithSetParameterValuesFormItems("Factory Command", factoryCommand))
	}

	tabsLayout := NewTabsLayout(pg.Name, "left")
	tabsLayout.Items = append(tabsLayout.Items, pages...)

	page := NewPageWithLayouts(pg.Name, tabsLayout)
	return page, nil
}

func MakePageForAddressInterface(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	param := info.Params.GetParameterDefine("TB2.P0CCC")
	if param == nil {
		return nil, errors.Errorf("cant find parameter define TB2.P0CCC")
	}
	if param.Child[0].Options == nil || len(param.Child[0].Options) <= 0 {
		return nil, nil
	}
	form := NewSetParameterValuesForm("Address Interface",
		makeParameter(param.Child[0], info),
		makeParameter(param.Child[1], info),
		makeParameter(param.Child[2], info),
		makeParameter(param.Child[3], info),
	)
	form.Actions = nil
	page := NewPageWithLayouts("Address Interface", NewSingleColLayoutWithItems(form))
	return page, nil
}
