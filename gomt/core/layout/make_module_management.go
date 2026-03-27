package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func makePageForServiceSwitch(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reWeekdayGroup := regexp.MustCompile(`(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday)`)
	reWeekday := regexp.MustCompile(`(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday)`)

	rows := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	listTable := NewTable("Service Schedule")
	listTable.AddTableColumnUnique("Day", -1)
	for _, gp := range pg.Child {
		if reWeekdayGroup.MatchString(gp.Name) {
			for _, obj := range gp.Child {
				name := strings.TrimSpace(reWeekday.ReplaceAllString(obj.Name, ""))
				listTable.AddTableColumn(name, -1)
			}
		}
	}

	viewItemsMap := map[string][]*Element{}
	for _, gp := range pg.Child {
		var key string
		if subs := reWeekdayGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
			key = subs[1]
		}
		rowIndex := 0
		for i, v := range rows {
			if v == key {
				rowIndex = i
				break
			}
		}
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		for _, obj := range gp.Child {
			name := strings.TrimSpace(reWeekday.ReplaceAllString(obj.Name, ""))
			elem := MustMakeElementFromParam(obj, info, groupPath)
			elem.SetName(name)
			for _, item := range elem.Items {
				item.SetName(strings.TrimSpace(reWeekday.ReplaceAllString(item.Name, "")))
			}

			listTable.MustSetTableRowData(int64(rowIndex), name, elem)
			viewItems, ok := viewItemsMap[key]
			if !ok {
				viewItems = []*Element{}
			}
			if name == "Working Hours (24h)" {
				viewItems = append(viewItems, elem.Items[0])
				viewItems = append(viewItems, elem.Items[1])
			} else {
				viewItems = append(viewItems, elem)
			}
			viewItemsMap[key] = viewItems
		}
	}
	for k, items := range viewItemsMap {
		rowIndex := 0
		for i, v := range rows {
			if v == k {
				rowIndex = i
				break
			}
		}
		viewPage := NewPageWithSetParameterValuesFormItems(k, items...)
		listTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
		listTable.MustSetTableRowData(int64(rowIndex), "Day", NewLabel("Day", k))
	}
	if listTable.Data != nil && len(listTable.Data) > 0 {
		return NewPageWithLayouts("Service Schedule",
			NewSingleColLayoutWithItems(listTable),
		), nil
	}
	return nil, nil
}

func makePageForServiceConfiguration(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	auModuleNumber := 4
	auModulePortNumber := 4
	groupNumber := 0
	ruModuleNumber := 0

	reInputModule := regexp.MustCompile(`(.*)(Radio|Input)\s*Module\s*(\d+)\s*(UL|DL)`)
	reAmplifierModule := regexp.MustCompile(`(Radio|Amplifier)\s*Module\s*(\d+)\s*(\(UL\/DL\)|Mapping)`)
	reGroup := regexp.MustCompile(`Group\s*(\d+)\s*`)

	for _, gp := range pg.Child {
		if strings.HasSuffix(gp.Name, "Operator/Service Configuration") {
			for _, obj := range gp.Child {
				if subs := reInputModule.FindStringSubmatch(obj.Name); len(subs) == 5 {
					moduleIndex := cast.ToInt(subs[3])
					if auModuleNumber < moduleIndex {
						auModuleNumber = moduleIndex
					}
				}
			}
		}
	}

	var (
		switchParam      *Element
		groupSelectParam *Element
		groupUpdateParam *Element

		auModuleListTable *Element
		groupTable        *Element
		groupViewItems    []*Element
	)
	groupTable = NewTable("RF Module Mapping")
	groupTable.AddTableColumnUnique("Group", -1)

	for _, gp := range pg.Child {
		if strings.HasSuffix(gp.Name, "RF Module Mapping") {
			for _, obj := range gp.Child {
				if reGroup.MatchString(obj.Name) {
					if n := len(obj.Params[0].Child); n > ruModuleNumber {
						ruModuleNumber = n
					}
				}
			}
		}
	}
	for i := 0; i < ruModuleNumber; i++ {
		moduleName := fmt.Sprintf("RF Module %v", i+1)
		groupTable.AddTableColumn(moduleName, -1)
	}

	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		if gp.Name == "Capacity Allocation Switch" || gp.Name == "Automatic RF Module Mapping" {
			for _, obj := range gp.Child {
				if obj.Name == "Capacity Allocation Switch" || obj.Name == "Automatic RF Module Mapping" {
					switchParam = MustMakeElementFromParam(obj, info, groupPath)
				}
			}
		} else if strings.HasSuffix(gp.Name, "Operator/Service Configuration") {
			auModuleListTable = NewTable(gp.Name)
			auModuleListTable.SetStyle("invalidKey", "uplink_frequency")
			auModuleListTable.SetStyle("invalidValue", "0-0")

			auModuleListTable.AddTableColumnUnique("Interface", -1)
			auModuleListTable.AddTableColumn("Uplink Frequency", -1)
			auModuleListTable.AddTableColumn("Downlink Frequency", -1)
			for _, obj := range gp.Child {
				if subs := reInputModule.FindStringSubmatch(obj.Name); len(subs) == 5 {
					deviceName := strings.TrimSpace(subs[1])
					moduleName := strings.TrimSpace(subs[2])
					moduleIndex := cast.ToInt(subs[3])
					deviceIndex := 0

					if strings.HasPrefix(deviceName, "1st") {
						deviceIndex = 1
					} else if strings.HasPrefix(deviceName, "2nd") {
						deviceIndex = 2
					}
					elem := MustMakeElementFromParam(obj, info, groupPath)
					name := ""
					if strings.HasSuffix(elem.Name, " UL") {
						name = "Uplink"
					} else if strings.HasSuffix(elem.Name, " DL") {
						name = "Downlink"
					} else {
						return nil, errors.Errorf("invalid element for %v", obj.Name)
					}
					elem.SetName(name)
					for i := 0; i < auModulePortNumber; i++ {
						portIndex := i + 1
						index := deviceIndex*auModuleNumber*auModulePortNumber + (moduleIndex-1)*auModulePortNumber + i
						p := NewParam(
							fmt.Sprintf("%v Module %v Port %v %v Frequency", moduleName, moduleIndex, portIndex, name),
							fmt.Sprintf("%v[%d]", elem.OID, i),
							elem.Access)
						auModuleListTable.MustSetTableRowData(int64(index), "Interface",
							NewLabel("Interface", fmt.Sprintf("%v %v-%v", deviceName, moduleIndex, portIndex)))
						auModuleListTable.MustSetTableRowData(int64(index), name+" Frequency", p)
					}
				}
			}
		} else if strings.HasSuffix(gp.Name, "RF Module Mapping") {
			for _, obj := range gp.Child {
				if subs := reGroup.FindStringSubmatch(obj.Name); len(subs) == 2 {
					groupIndex := cast.ToInt(subs[1])
					if groupNumber < groupIndex {
						groupNumber = groupIndex
					}
					if groupIndex < 1 {
						continue
					}
					index := groupIndex - 1
					elem := MustMakeElementFromParam(obj, info, groupPath)
					for i := 0; i < ruModuleNumber; i++ {
						moduleName := fmt.Sprintf("RF Module %v", i+1)
						param := NewParam(fmt.Sprintf("Group %v %v", groupIndex, moduleName), fmt.Sprintf("%v[%v]", elem.OID, i), elem.Access)
						groupTable.MustSetTableRowData(int64(index), moduleName, param)
					}
				}
			}
		} else if gp.Name == "RF Module Mapping Configuration" {
			for _, obj := range gp.Child {
				if obj.Name == "Capacity Group Update" {
					groupUpdateParam = MustMakeElementFromParam(obj, info, groupPath)
				} else if obj.Name == "Capacity Group" || obj.Name == "Group" {
					groupSelectParam = MustMakeElementFromParam(obj, info, groupPath)
				} else if subs := reAmplifierModule.FindStringSubmatch(obj.Name); len(subs) == 4 {
					moduleType := subs[1]
					moduleIndex := cast.ToInt(subs[2])
					if groupViewItems == nil {
						groupViewItems = []*Element{}
					}

					elem := MustMakeElementFromParam(obj, info, groupPath)
					if strings.HasSuffix(elem.Name, "(UL/DL)") {
						label := NewLabel("RF Module", fmt.Sprintf("%v Module %v", moduleType, moduleIndex))
						groupViewItems = append(groupViewItems, label)
						p1 := NewParam("Uplink Frequency", fmt.Sprintf("%v[0]", elem.OID), elem.Access)
						p1.SetStyle("readonly", true)
						p2 := NewParam("Downlink Frequency", fmt.Sprintf("%v[1]", elem.OID), elem.Access)
						p2.SetStyle("readonly", true)
						groupViewItems = append(groupViewItems, p1)
						groupViewItems = append(groupViewItems, p2)
					} else {
						elem.SetName("Mapping")
						groupViewItems = append(groupViewItems, elem)
					}

				}
			}
		}
	}

	for i := 0; i < groupNumber; i++ {
		getItems := []*Element{}
		setItems := []*Element{}
		if groupSelectParam != nil {
			selectParam := *groupSelectParam
			selectParam.SetStyle("readonly", true)
			selectParam.SetValue(fmt.Sprintf("%02X", i))
			setItems = append(setItems, &selectParam)
			getItems = append(getItems, &selectParam)
		}

		for _, v := range groupViewItems {
			if v.Type == "Param" {
				tmp := v
				setItems = append(setItems, tmp)
				getItems = append(getItems, tmp)
			}
		}

		if groupUpdateParam != nil {
			updateParam := *groupUpdateParam
			updateParam.SetValue("01")
			updateParam.SetStyle("hidden", true)
			updateParam.RemoveStyle("input")
			setItems = append(setItems, &updateParam)
		}

		name := fmt.Sprintf("Group %v", i+1)
		viewForm := NewSetParameterValuesForm(name, groupViewItems...)
		viewPageItems := []*Element{viewForm}
		if switchParam != nil {
			viewForm.SetStyle("disableParam", switchParam.OID)
			viewForm.SetStyle("disableValue", "00")
		}
		viewPage := NewPageWithLayouts(name, NewSingleColLayoutWithItems(viewPageItems...))
		viewAction := NewViewPageAction(viewPage)
		if switchParam != nil {
			viewAction.SetStyle("disableParam", switchParam.OID)
			viewAction.SetStyle("disableValue", "00")
		}

		groupTable.MustSetTableRowData(int64(i), "Group",
			NewLabel("Group", fmt.Sprintf("%v", i+1)))
		groupTable.SetTableRowClickActions(int64(i), NewGetParameterValuesAction(getItems...), viewAction)
	}

	if switchParam != nil {
		action := NewAction("StateAction")
		action.SetAction("00", NewSetParameterValuesAction(CopyParamWithValue(switchParam, "01")))
		action.SetAction("01", NewSetParameterValuesAction(CopyParamWithValue(switchParam, "00")))

		switchParam.SetAction("click", action)
		switchParam.SetStyle("input", "button")
		switchParam.SetStyle("inactiveValue", "00")
		// switchParam.SetStyle("inactiveText", "Capacity Allocation Disabled")
		switchParam.SetStyle("inactiveText", "Capacity Allocation Automatic")
		switchParam.SetStyle("inactiveType", "info")
		switchParam.SetStyle("activeValue", "01")
		switchParam.SetStyle("activeText", "Capacity Allocation Manual   ")
		// switchParam.SetStyle("activeText", "Capacity Allocation Enabled ")
		switchParam.SetStyle("activeType", "success")
		groupTable.SetAction("toolbar", NewToolbarWithItems("", switchParam))
	}

	pageItems := []*Element{}
	if groupTable != nil {
		// groupTable.SetStyle("height", "140px")
		pageItems = append(pageItems, groupTable)
	}
	if switchParam != nil {
		disableAlert := NewAlert("Disable Alert", "Capacity allocation is disabled, please enable before configuration.", "warning")
		disableAlert.SetStyle("visibleParam", switchParam.OID)
		disableAlert.SetStyle("visibleValue", "00")
		pageItems = append(pageItems, disableAlert)
	}
	pageItems = append(pageItems, auModuleListTable)
	page := NewPageWithLayouts("RF Module Mapping",
		NewColLayoutWithItems(0, pageItems...),
	)
	return page, nil
}

func MakeModuleForCapacityAllocationManagement(m *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	module := NewModule("Capacity Allocation")

	for _, pg := range m.Child {
		if pg.Name == "Service Configuration" {
			page, err := makePageForServiceConfiguration(pg, info, path)
			if err != nil {
				return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
			}
			addPageToModule(page, module)
		} else if pg.Name == "Service Switch" {
			page, err := makePageForServiceSwitch(pg, info, path)
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
