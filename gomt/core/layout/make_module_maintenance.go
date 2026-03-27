package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func MakeModuleForMaintenance(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	module := NewModule("Maintenance")

	var syncLossCounterReset *Element
	for _, pg := range m.Child {
		if pg.Name == "Engineering" {
			for _, gp := range pg.Child {
				if gp.Name == "OP Info" {
					groupPath := append([]string{}, path...)
					groupPath = append(groupPath, pg.Name, gp.Name)
					for _, obj := range gp.Child {
						if obj.Name == "Sync Loss Counter Reset" {
							syncLossCounterReset = MustMakeElementFromParam(obj, info, groupPath)
						}
					}
				}
			}
		}
	}

	for _, pg := range m.Child {
		if pg.Name == "Optical Info" {
			page, err := makePageForMaintenanceOpticalInfo(pg, syncLossCounterReset, info, path)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
		} else if pg.Name == "Engineering" {
			page, err := MakePageWithExcludes(m, pg, info, []string{"OP Info"})
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
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

func makePageForMaintenanceOpticalInfo(pg *model.ParameterTreeNode, syncLossCounterReset *Element, info DeviceInfo, path []string) (*Element, error) {
	reOP := regexp.MustCompile(`OP\s*(\d+|P|S|M)\s*(Transceiver|)`)

	table := NewTable("Optical Module Information List")
	table.AddTableColumnUnique("Optical Module", -1)
	for _, gp := range pg.Child {
		if gp.Name == "Optical Module Info" {
			for _, obj := range gp.Child {
				p := obj.Params[0]
				if p.Child != nil && len(p.Child) > 0 {
					for _, v := range p.Child {
						table.AddTableColumn(v.Name, -1)
					}
				} else if reOP.FindString(obj.Name) != "" {
					name := strings.TrimSpace(reOP.ReplaceAllString(obj.Name, ""))
					table.AddTableColumn(name, -1)
				}
			}
		}
	}

	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		if gp.Name == "Optical Module Info" {
			for _, obj := range gp.Child {
				p := obj.Params[0]
				p.SetPath(groupPath)
				indexText := ""
				if subs := reOP.FindStringSubmatch(p.Name); len(subs) == 3 {
					indexText = subs[1]
				}
				index := getOpticalIndex(indexText)
				if index <= 0 {
					continue
				}

				rowIndex := index - 1
				table.MustSetTableRowData(rowIndex, "Optical Module", NewLabel("Optical Module", fmt.Sprintf("OP %v", indexText)))

				if p.Child != nil && len(p.Child) > 0 {
					for _, v := range p.Child {
						elem := NewParam(v.Name, fmt.Sprintf("%v", v.PrivOid), string(v.Access))
						elem.SetStyle("nullValue", "Invalid")
						table.MustSetTableRowData(rowIndex, v.Name, elem)
					}
				} else if reOP.FindString(obj.Name) != "" {
					name := strings.TrimSpace(string(reOP.ReplaceAllString(obj.Name, "")))
					elem := MustMakeElementFromParam(obj, info, groupPath)
					elem.SetName(name)
					table.MustSetTableRowData(rowIndex, name, elem)
				}
			}
		}
	}
	if syncLossCounterReset != nil {
		table.SetAction("toolbar", NewToolbarWithItems("", syncLossCounterReset))
	}
	page := NewPageWithLayouts(pg.Name, NewSingleColLayoutWithItems(table))
	return page, nil
}
