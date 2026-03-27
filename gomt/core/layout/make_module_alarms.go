package layout

import (
	"fmt"
	"gomt/core/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type _AlarmElements struct {
	Name       string
	Index      *Element
	Status     *Element
	Severity   *Element
	Threshold  *Element
	Label      *Element
	Indication *Element
	Mode       *Element
}

type _AlarmElementsList []*_AlarmElements

func (m _AlarmElementsList) Get(name string) *_AlarmElements {
	for _, v := range m {
		if v.Name == name {
			return v
		}
	}
	return nil
}
func (m *_AlarmElementsList) Set(v *_AlarmElements) {
	index := -1
	for i, v2 := range *m {
		if v2.Name == v.Name {
			index = i
			break
		}
	}
	if index < 0 {
		*m = append(*m, v)
	} else {
		(*m)[index] = v
	}
}

func (m *_AlarmElementsList) SetElem(name string, typ string, p *Element) *_AlarmElements {
	v := m.Get(name)
	if v == nil {
		v = &_AlarmElements{
			Name: name,
		}
	}
	switch typ {
	case "Index":
		v.Index = p
	case "Status":
		v.Status = p
	case "Severity":
		v.Severity = p
	case "Threshold":
		v.Threshold = p
	case "Label":
		v.Label = p
	case "Indication":
		v.Indication = p
	case "Mode":
		v.Mode = p
	}
	m.Set(v)
	return v
}

func MakeModuleForAlarms(m *model.ParameterTreeNode, info DeviceInfo) (*Element, error) {
	reModule := regexp.MustCompile(`(Radio|Input) Module\s*(\d+)`)
	rePort := regexp.MustCompile(`Port\s*(\d+)`)
	reOptical := regexp.MustCompile(`OP\s*(\d+|S|M)`)
	reInput := regexp.MustCompile(`External Input \s*(\d+)`)
	reOutput := regexp.MustCompile(`External Output Alarm\s*(\d+)`)

	alarmsList := &_AlarmElementsList{}
	externalInputAlarmsList := &_AlarmElementsList{}
	alarmIndicationList := &_AlarmElementsList{}

	var vswrThresholdElem *Element

	portNum := int64(0)
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			for _, obj := range gp.Child {
				var portIndex int64 = -1
				if subs := rePort.FindStringSubmatch(obj.Name); len(subs) == 2 {
					portIndex, _ = strconv.ParseInt(subs[1], 10, 32)
					if portIndex > int64(portNum) {
						portNum = portIndex
					}
				}
			}
		}
	}
	for _, pg := range m.Child {
		for _, gp := range pg.Child {
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, []string{"Alarms"})
				var (
					moduleIndex    int64  = -1
					portIndex      int64  = -1
					opIndex        string = ""
					extInputIndex  int64  = -1
					extOutputIndex int64  = -1
				)
				modulePrefix := ""
				if subs := reModule.FindStringSubmatch(gp.Name); len(subs) == 3 {
					moduleIndex, _ = strconv.ParseInt(subs[2], 10, 32)
					modulePrefix = subs[1]
				}
				if moduleIndex < 0 {
					if subs := reModule.FindStringSubmatch(obj.Name); len(subs) == 3 {
						moduleIndex, _ = strconv.ParseInt(subs[2], 10, 32)
						modulePrefix = subs[1]
					}
				}
				if subs := rePort.FindStringSubmatch(obj.Name); len(subs) == 2 {
					portIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if subs := reOptical.FindStringSubmatch(obj.Name); len(subs) == 2 {
					opIndex = subs[1]
				}
				if subs := reInput.FindStringSubmatch(obj.Name); len(subs) == 2 {
					extInputIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if subs := reOutput.FindStringSubmatch(obj.Name); len(subs) == 2 {
					extOutputIndex, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				var list *_AlarmElementsList = nil
				name := formatAlarmName(obj.Name)
				key := name
				if strings.HasPrefix(name, "Band ") {
					list = alarmsList
				} else if strings.HasSuffix(name, "Input Over Power Threshold") || strings.HasSuffix(name, "Input Under Power Threshold") {
					list = alarmsList
				} else if extInputIndex > 0 {
					key = fmt.Sprintf("External Input %v Alarm", extInputIndex)
					list = externalInputAlarmsList
				} else if extOutputIndex > 0 {
					key = fmt.Sprintf("Alarm Indication %v", extOutputIndex)
					list = alarmIndicationList
				} else if opIndex != "" {
					list = alarmsList
				} else if moduleIndex > 0 && portIndex > 0 {
					key = fmt.Sprintf("%v Module %v ", modulePrefix, moduleIndex) + key
					list = alarmsList
				} else if moduleIndex > 0 {
					list = alarmsList
				} else {
					list = alarmsList
				}

				if strings.HasSuffix(name, " Alarm") {
					if obj.Params[0].InputType == "status:alarm" {
						statusParam := elem
						if !strings.HasSuffix(statusParam.Name, " Status") {
							statusParam.Name = statusParam.Name + " Status"
						}
						list.SetElem(key, "Status", statusParam)
					} else {
						severityParam := elem
						if !strings.HasSuffix(severityParam.Name, " Severity") {
							severityParam.Name = severityParam.Name + " Severity"
						}
						list.SetElem(key, "Severity", severityParam)
					}
				} else if strings.HasSuffix(name, " Level") {
					list.SetElem(key, "Severity", elem)
				} else if strings.HasSuffix(name, " Indication") {
					list.SetElem(key, "Indication", elem)
				} else if strings.HasSuffix(name, " Mode Select") || strings.HasSuffix(name, " Mode") {
					list.SetElem(key, "Mode", elem)
				} else if strings.HasSuffix(name, " Label") {
					list.SetElem(key, "Label", elem)
				} else if strings.HasSuffix(name, " Threshold") {
					if name == "VSWR Threshold" {
						vswrThresholdElem = elem
					} else if strings.HasSuffix(name, "Input Over Power Threshold") || strings.HasSuffix(name, "Input Under Power Threshold") {
						re := regexp.MustCompile(`(Radio|Input) Module\s*(\d+) (.*) Threshold`)
						for i := int64(0); i < portNum; i++ {
							key2 := re.ReplaceAllString(key, fmt.Sprintf("$1 Module $2 Port %v $3 Alarm", i+1))
							list.SetElem(key2, "Threshold", elem)
						}
					} else {
						key2 := strings.ReplaceAll(key, "Threshold", "Alarm")
						list.SetElem(key2, "Threshold", elem)
					}
				} else if strings.HasSuffix(name, " Mode Select") {
					list.SetElem(key, "Mode", elem)
				} else {
					logrus.Error(errors.Errorf("unknow alarm parameter %v", obj.Name))
				}
			}
		}
	}

	if vswrThresholdElem != nil {
		for _, v := range *alarmsList {
			re := regexp.MustCompile(`(Radio|Amplifier) Module\s*(\d+)\s*VSWR\s*Alarm`)
			if re.MatchString(v.Name) {
				alarmsList.SetElem(v.Name, "Threshold", vswrThresholdElem)
			}
		}
	}

	alarmsTable := NewTable("Device Alarms")
	alarmsTable.AddTableColumnUnique("Alarm Name", 320)
	alarmsTable.AddTableColumn("Severity Level", -1)
	alarmsTable.AddTableColumn("Status", -1)
	alarmsTable.AddTableColumn("Threshold", -1)

	for i, v := range *alarmsList {
		rowIndex := i
		alarmName := v.Name
		def := info.Alarms.GetAlarmDefine(v.Name)
		if def != nil {
			alarmName = def.Name
		}
		elems := map[string]*Element{}
		if v.Severity != nil {
			elems["Severity Level"] = v.Severity
		} else {
			elems["Severity Level"] = NewLabel("Severity Level", "")
		}
		if v.Status != nil {
			elems["Status"] = v.Status
		} else {
			elems["Status"] = NewLabel("Status", "")
		}
		if v.Threshold != nil {
			elems["Threshold"] = v.Threshold
		} else {
			elems["Threshold"] = NewLabel("Threshold", "-")
		}
		viewItems := []*Element{}
		orderKeys := []string{"Severity Level", "Status", "Threshold"}
		for _, k := range orderKeys {
			elem := elems[k]
			alarmsTable.MustSetTableRowData(int64(rowIndex), k, elem)
			if elem.Key == "threshold" && elem.Type == "Label" {
				//ingore
			} else {
				viewItems = append(viewItems, elem)
			}
		}

		viewPage := NewPageWithSetParameterValuesFormItems(alarmName, viewItems...)
		alarmsTable.MustSetTableRowData(int64(rowIndex), "Alarm Name", NewLabel("Alarm Name", alarmName))
		alarmsTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	externalInputAlarmsTable := NewTable("External Input Alarms")
	externalInputAlarmsTable.AddTableColumnUnique("Alarm Name", 320)
	externalInputAlarmsTable.AddTableColumn("Severity Level", -1)
	externalInputAlarmsTable.AddTableColumn("Status", -1)
	externalInputAlarmsTable.AddTableColumn("Label", -1)
	externalInputAlarmsTable.AddTableColumn("Mode", -1)

	for i, v := range *externalInputAlarmsList {
		rowIndex := i
		alarmName := v.Name
		def := info.Alarms.GetAlarmDefine(v.Name)
		if def != nil {
			alarmName = def.Name
		}
		elems := map[string]*Element{}
		if v.Severity != nil {
			elems["Severity Level"] = v.Severity
		} else {
			elems["Severity Level"] = NewLabel("Severity Level", "")
		}
		if v.Status != nil {
			elems["Status"] = v.Status
		} else {
			elems["Status"] = NewLabel("Status", "")
		}
		if v.Label != nil {
			elems["Label"] = v.Label
		} else {
			elems["Label"] = NewLabel("Label", "")
		}
		if v.Mode != nil {
			elems["Mode"] = v.Mode
		} else {
			elems["Mode"] = NewLabel("Mode", "")
		}

		viewItems := []*Element{}
		orderKeys := []string{"Severity Level", "Status", "Label", "Mode"}
		for _, k := range orderKeys {
			elem := elems[k]
			externalInputAlarmsTable.MustSetTableRowData(int64(rowIndex), k, elem)
			viewItems = append(viewItems, elem)
		}
		viewPage := NewPageWithSetParameterValuesFormItems(alarmName, viewItems...)
		externalInputAlarmsTable.MustSetTableRowData(int64(rowIndex), "Alarm Name", NewLabel("Alarm Name", alarmName))
		externalInputAlarmsTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	alarmIndicationTable := NewTable("Alarm Indications")
	alarmIndicationTable.AddTableColumnUnique("ID", 320)
	alarmIndicationTable.AddTableColumn("Mode", -1)
	alarmIndicationTable.AddTableColumn("Severity Level", -1)
	alarmIndicationTable.AddTableColumn("Indication", -1)
	for i, v := range *alarmIndicationList {
		rowIndex := i
		alarmName := v.Name

		elems := map[string]*Element{}
		if v.Mode != nil {
			elems["Mode"] = v.Mode
		} else {
			elems["Mode"] = NewLabel("Mode", "")
		}
		if v.Severity != nil {
			elems["Severity Level"] = v.Severity
		} else {
			elems["Severity Level"] = NewLabel("Severity Level", "")
		}
		if v.Indication != nil {
			elems["Indication"] = v.Indication
		} else {
			elems["Indication"] = NewLabel("Indication", "")
		}

		viewItems := []*Element{}
		orderKeys := []string{"Mode", "Severity Level", "Indication"}
		for _, k := range orderKeys {
			elem := elems[k]
			alarmIndicationTable.MustSetTableRowData(int64(rowIndex), k, elem)
			viewItems = append(viewItems, elem)
		}
		viewPage := NewPageWithSetParameterValuesFormItems(alarmName, viewItems...)
		alarmIndicationTable.MustSetTableRowData(int64(rowIndex), "ID", NewLabel("ID", v.Name))
		alarmIndicationTable.SetTableRowClickAction(int64(rowIndex), NewViewPageAction(viewPage))
	}

	tabsLayout := NewTabsLayout("Element Alarms", "top")

	if alarmsTable.Data != nil && len(alarmsTable.Data) > 0 {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Device",
			NewSingleColLayoutWithItems(alarmsTable),
		))
	}
	if externalInputAlarmsTable.Data != nil && len(externalInputAlarmsTable.Data) > 0 {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("External Input",
			NewSingleColLayoutWithItems(externalInputAlarmsTable),
		))
	}
	if alarmIndicationTable.Data != nil && len(alarmIndicationTable.Data) > 0 {
		tabsLayout.Items = append(tabsLayout.Items, NewPageWithLayouts("Alarm Indication",
			NewSingleColLayoutWithItems(alarmIndicationTable),
		))
	}

	module := NewModule(m.Name)
	module.Items = append(module.Items, NewPageWithLayouts("Element Alarms", tabsLayout))

	return module, nil
}

func formatAlarmName(name string) string {
	if strings.HasPrefix(name, "Fan") && strings.HasSuffix(name, "Loss") {
		name = name + " Alarm"
	}
	return name
}
