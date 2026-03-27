package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"
)

func MakeModuleForCombiners(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reCombinerPage := regexp.MustCompile(`Combiner\s*(\d+)`)
	rePortGroup := regexp.MustCompile(`Port\s*(\d+)`)

	reCombiner := regexp.MustCompile(`Combiner\s*(\d+)`)
	reCombinerPort := regexp.MustCompile(`Combiner\s*(\d+) Port\s*(\d+)`)

	module := NewModule(m.Name)

	moduleListTable := NewTable("Combiner List")
	moduleListTable.AddTableColumnUnique("Combiner", -1)

	portListTable := NewTable("Combiner Port List")
	portListTable.AddTableColumnUnique("Combiner Port", -1)

	var portNum int64 = 0
	// pg := m.Child[0]
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			var portIndex int64 = -1
			if subs := rePortGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				if portIndex, _ = strconv.ParseInt(subs[1], 10, 32); portIndex > portNum {
					portNum = portIndex
				}
				if portIndex > 0 {
					for _, obj := range gp.Child {
						name := strings.TrimSpace(reCombinerPort.ReplaceAllString(obj.Name, ""))
						portListTable.AddTableColumn(name, -1)
					}
				}
			} else {
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reCombiner.ReplaceAllString(obj.Name, ""))
					moduleListTable.AddTableColumn(name, -1)
				}
			}
		}
	}

	portViewItemsMap := map[string][]*Element{}
	moduleViewItemsMap := map[string][]*Element{}
	for _, pg := range m.Child {
		var combinerIndex int64 = -1
		if subs := reCombinerPage.FindStringSubmatch(pg.Name); len(subs) == 2 {
			combinerIndex, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if combinerIndex < 1 {
			continue
		}
		for _, gp := range pg.Child {
			groupPath := append([]string{}, path...)
			groupPath = append(groupPath, pg.Name, gp.Name)

			var portIndex int64 = -1
			if subs := rePortGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
				portIndex, _ = strconv.ParseInt(subs[1], 10, 32)
			}
			if portIndex > 0 {
				rowIndex := (combinerIndex-1)*portNum + portIndex - 1
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reCombinerPort.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)

					portListTable.MustSetTableRowData(rowIndex, name, elem)
					key := fmt.Sprintf("%v-%v", combinerIndex, portIndex)
					viewItems, ok := portViewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					portViewItemsMap[key] = viewItems
				}
			} else {
				rowIndex := combinerIndex - 1
				for _, obj := range gp.Child {
					name := strings.TrimSpace(reCombiner.ReplaceAllString(obj.Name, ""))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)
					for _, item := range elem.Items {
						n := strings.TrimSpace(reCombiner.ReplaceAllString(item.Name, ""))
						item.SetName(n)
					}
					moduleListTable.MustSetTableRowData(rowIndex, name, elem)
					key := fmt.Sprintf("%v", combinerIndex)
					viewItems, ok := moduleViewItemsMap[key]
					if !ok {
						viewItems = []*Element{}
					}
					viewItems = append(viewItems, elem)
					moduleViewItemsMap[key] = viewItems

				}
			}

		}
	}

	for k, items := range moduleViewItemsMap {
		combinerIndex, _ := strconv.ParseInt(k, 10, 32)
		rowIndex := combinerIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Combiner %v", combinerIndex),
			items...)
		moduleListTable.MustSetTableRowData(int64(rowIndex), "Combiner",
			NewLabel("Combiner", fmt.Sprintf("%v", combinerIndex)))
		moduleListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	moduleListPage := NewPageWithLayouts("Combiners",
		NewSingleColLayoutWithItems(moduleListTable),
	)

	for k, items := range portViewItemsMap {
		parts := strings.Split(k, "-")
		combinerIndex, _ := strconv.ParseInt(parts[0], 10, 32)
		portIndex, _ := strconv.ParseInt(parts[1], 10, 32)
		rowIndex := (combinerIndex-1)*portNum + portIndex - 1
		viewPage := NewPageWithSetParameterValuesFormItems(
			fmt.Sprintf("Combiner %v Port %v", combinerIndex, portIndex),
			items...)
		portListTable.MustSetTableRowData(int64(rowIndex), "Combiner Port",
			NewLabel("Combiner Port", fmt.Sprintf("%v-%v", combinerIndex, portIndex)))
		portListTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}
	portListPage := NewPageWithLayouts("Combiner Ports",
		NewSingleColLayoutWithItems(portListTable),
	)
	module.Items = append(module.Items, moduleListPage, portListPage)
	return module, nil
}
