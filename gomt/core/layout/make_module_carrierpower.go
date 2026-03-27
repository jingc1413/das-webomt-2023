package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakeModuleForCarrierPower(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reModulePage := regexp.MustCompile(`Module\s*(\d+)`)
	reCarrierGroup := regexp.MustCompile(`Carrier\s*(\d+)`)
	reModuleCarrier := regexp.MustCompile(`Radio Module\s*(\d+) Carrier\s*(\d+)`)

	module := NewModule(m.Name)

	carrierListTable := NewTable("Carrier List")
	carrierListTable.AddTableColumnUnique("Carrier", -1)

	var carrierNum int64 = 0
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			var carrierIndex int64 = -1
			if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				if carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32); carrierIndex > carrierNum {
					carrierNum = carrierIndex
				}
			}
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
				carrierListTable.AddTableColumn(name, -1)
			}
		}
	}

	carrierViewItemsMap := map[string][]*Element{}

	for _, pg := range m.Child {
		var moduleIndex int64 = -1
		if subs := reModulePage.FindStringSubmatch(pg.Name); len(subs) == 2 {
			moduleIndex, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if moduleIndex < 1 {
			continue
		}
		for _, gp := range pg.Child {
			var carrierIndex int64 = -1
			if subs := reCarrierGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				carrierIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if carrierIndex < 1 {
				continue
			}
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)

			rowIndex := (moduleIndex-1)*carrierNum + carrierIndex - 1
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reModuleCarrier.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)

				carrierListTable.MustSetTableRowData(rowIndex, name, elem)
				key := fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)
				viewItems, ok := carrierViewItemsMap[key]
				if !ok {
					viewItems = []*Element{}
				}
				viewItems = append(viewItems, elem)
				carrierViewItemsMap[key] = viewItems
			}
		}
	}
	for k, items := range carrierViewItemsMap {
		parts := strings.Split(k, "-")
		moduleIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		carrierIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (moduleIndex-1)*carrierNum + carrierIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Module %v Carrier %v", moduleIndex, carrierIndex),
			items...)
		carrierListTable.MustSetTableRowData(int64(rowIndex), "Carrier",
			NewLabel("Carrier", fmt.Sprintf("%v-%v", moduleIndex, carrierIndex)))
		carrierListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	moduleConfigurationPage := NewPageWithLayouts("Power Configuration",
		NewSingleColLayoutWithItems(carrierListTable),
	)
	module.Items = append(module.Items, moduleConfigurationPage)
	return module, nil
}
