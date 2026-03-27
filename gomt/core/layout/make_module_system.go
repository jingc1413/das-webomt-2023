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

func MakePageForSNMPConfiguration(pg *model.ParameterTreeNode, pg2 *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	var (
		applyItem  *Element = nil
		resetItem  *Element = nil
		deleteItem *Element = nil
	)
	trapFormItems := []*Element{}
	trapResendFormItems := []*Element{}

	reTrap := regexp.MustCompile(`Trap\s*(\d+)`)
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		for _, obj := range gp.Child {
			if gp.Name == "Trap Settings" && reTrap.FindString(obj.Name) != "" {
				continue
			}
			elem := MustMakeElementFromParam(obj, info, groupPath)
			if strings.HasPrefix(gp.Name, "SNMPV3 USM") {
				if strings.HasSuffix(elem.Name, "User Apply") {
					applyItem = elem
				} else if strings.HasSuffix(elem.Name, "Reset USM") {
					resetItem = elem
				} else if strings.HasSuffix(elem.Name, "User Delete") {
					deleteItem = elem
				}
			} else if gp.Name == "Trap Settings" {
				trapFormItems = append(trapFormItems, elem)
			} else if gp.Name == "Trap Resend" {
				if obj.Name == "Community" {
					trapFormItems = append(trapFormItems, elem)
				} else {
					trapResendFormItems = append(trapResendFormItems, elem)
				}
			}
		}
	}
	trapFormItems = append(trapFormItems, trapResendFormItems...)

	tabsLayout := NewTabsLayout(pg.Name, "left")
	if tabPage := NewPageWithLayouts("General", NewSingleColLayoutWithItems(
		NewSetParameterValuesForm("Trap settings", trapFormItems...),
		// NewSetParameterValuesForm("Trap Resend", trapResendFormItems...),
	)); tabPage != nil {
		tabsLayout.Items = append(tabsLayout.Items, tabPage)
	}

	trapAddressTable, err := makeTableForTrapAddress(pg, info, path)
	if err != nil {
		return nil, errors.Wrap(err, "make table for snmp trap address")
	}
	if tabPage := NewPageWithLayouts(trapAddressTable.Name, NewSingleColLayoutWithItems(trapAddressTable)); tabPage != nil {
		tabsLayout.Items = append(tabsLayout.Items, tabPage)
	}

	if pg2 != nil {
		userTable, err := makeTableForSNMPUserInfo(pg2, deleteItem, info, path)
		if err != nil {
			return nil, errors.Wrap(err, "make table for snmp user info")
		}
		toolbar := userTable.GetAction("toolbar")
		if toolbar == nil {
			toolbar = NewToolbar("")
		}
		if applyItem != nil {
			btn := *applyItem
			btn.SetName("Apply USM")
			if clickAction := btn.Actions["click"]; clickAction != nil {
				btn.SetAction("click", NewMultipleActionsAction(clickAction, NewAction("SetDeviceUnavailable")))
			}

			btn.SetStyle("confirmTitle", "WARNING")
			btn.SetStyle("confirmMessage", "This command will reboot device, do you want to continue?")
			btn.SetStyle("confirmType", "warning")
			btn.SetStyle("willReboot", true)
			toolbar.Items = append(toolbar.Items, &btn)
		}
		if resetItem != nil {
			btn := *resetItem
			if clickAction := btn.Actions["click"]; clickAction != nil {
				btn.SetAction("click", NewMultipleActionsAction(clickAction, NewAction("SetDeviceUnavailable")))
			}
			btn.SetStyle("confirmTitle", "WARNING")
			btn.SetStyle("confirmMessage", "This command will reboot device, do you want to continue?")
			btn.SetStyle("confirmType", "warning")
			btn.SetStyle("willReboot", true)
			toolbar.Items = append(toolbar.Items, &btn)
		}
		if len(toolbar.Items) > 0 {
			userTable.SetAction("toolbar", NewToolbarWithItems("", toolbar.Items...))
		}
		if tabPage := NewPageWithLayouts(userTable.Name, NewSingleColLayoutWithItems(userTable)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	page := NewPageWithLayouts("SNMP Configuration", tabsLayout)
	return page, nil
}

func makeTableForSNMPUserInfo(pg *model.ParameterTreeNode, deleteItem *Element, info DeviceInfo, path []string) (*Element, error) {
	reUserGroup := regexp.MustCompile(`User\s*(\d+)`)
	reUser := regexp.MustCompile(`User\s*(\d+)`)

	var deleteDef *model.ParameterDefine = nil
	if deleteItem != nil {
		deleteDef = info.Params.GetParameterDefine(model.PrivObjectId(deleteItem.OID))
	}

	table := NewTable("SNMP User List")
	table.AddTableColumnUnique("User", -1)
	for _, obj := range pg.Child[0].Child {
		name := strings.TrimSpace(reUser.ReplaceAllString(obj.Name, ""))
		if name == "Edit User Confirm" {
			continue
		}
		table.AddTableColumn(name, -1)
	}
	table.AddTableColumn("Action", -1)

	for _, gp := range pg.Child {
		var index int64 = -1
		if subs := reUserGroup.FindStringSubmatch(gp.Name); len(subs) == 2 {
			index, _ = strconv.ParseInt(subs[1], 10, 32)
		}
		if index < 1 {
			continue
		}
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		rowIndex := index - 1
		viewFormItems := []*Element{}
		var (
			privacyProtocol        *Element
			privacyPassword        *Element
			authenticationProtocol *Element
			authenticationPassword *Element
		)
		for _, obj := range gp.Child {
			name := strings.TrimSpace(reUser.ReplaceAllString(obj.Name, ""))
			elem := MustMakeElementFromParam(obj, info, groupPath)
			elem.SetName(name)

			if name == "Privacy Protocol" {
				privacyProtocol = elem
			} else if name == "Privacy Password" {
				privacyPassword = elem
			} else if name == "Authentication Protocol" {
				authenticationProtocol = elem
			} else if name == "Authentication Password" {
				authenticationPassword = elem
			}

			if name == "Edit User Confirm" {
				elem.SetValue("01")
				elem.SetStyle("hidden", true)
				elem.RemoveStyle("input")
				viewFormItems = append(viewFormItems, elem)
			} else {
				table.MustSetTableRowData(rowIndex, name, elem)
				viewFormItems = append(viewFormItems, elem)
			}
		}
		if privacyProtocol != nil && privacyPassword != nil {
			privacyPassword.SetStyle("disableParam", privacyProtocol.OID)
			privacyPassword.SetStyle("disableInputValue", "00")
		}
		if authenticationProtocol != nil && authenticationPassword != nil {
			authenticationPassword.SetStyle("disableParam", authenticationProtocol.OID)
			authenticationPassword.SetStyle("disableInputValue", "00")
		}
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("User %v", index), viewFormItems...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "User", NewLabel("User", fmt.Sprintf("%v", index)))

		if deleteDef != nil {
			value := ""
			for k, v := range deleteDef.Options {
				if v == fmt.Sprintf("User%v", index) {
					value = k
					break
				}
			}
			if value != "" {
				param := *deleteItem
				param.SetValue(value)
				action := NewSetParameterValuesAction(&param)
				btn := NewButtonWithAction("Delete", "Delete", "text", action)
				if clickAction := btn.Actions["click"]; clickAction != nil {
					btn.SetAction("click", NewMultipleActionsAction(clickAction, NewAction("SetDeviceUnavailable")))
				}
				btn.SetStyle("confirmTitle", "WARNING")
				btn.SetStyle("confirmMessage", "This command will reboot device, do you want to continue?")
				btn.SetStyle("confirmType", "warning")
				btn.SetStyle("willReboot", true)
				table.MustSetTableRowData(rowIndex, "Action", btn)
			}
		}
	}
	// for rowIndex, row := range table.Data {
	// 	deleteItems := []*Element{}
	// 	for k, v := range row {
	// 		item := *v
	// 		if k == "security_username" || k == "privacy_password" || k == "authentication_password" {
	// 			item.SetValue("")
	// 			deleteItems = append(deleteItems, &item)
	// 		} else if k == "privacy_protocol" || k == "authentication_protocol" {
	// 			item.SetValue("00")
	// 			deleteItems = append(deleteItems, &item)
	// 		}
	// 	}
	// 	action := NewSetParameterValuesAction(deleteItems...)
	// 	table.MustSetTableRowData(int64(rowIndex), "Action", NewButtonWithAction("Delete", "text", action))
	// }

	addButton := NewButtonWithAction("Add", "Add", "primary", NewAction("AddItem"))
	table.SetAction("toolbar", NewToolbarWithItems("", addButton))
	table.SetStyle("invalidKey", "security_username")
	table.SetStyle("invalidValue", "")
	table.SetStyle("viewKey", "user")
	return table, nil
}

func makeTableForTrapAddress(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	reTrap := regexp.MustCompile(`Trap\s*(\d+)`)

	table := NewTable("Trap Settings List")
	table.AddTableColumnUnique("Trap", -1)

	for _, gp := range pg.Child {
		if gp.Name == "Trap Settings" {
			for _, obj := range gp.Child {
				if reTrap.FindString(obj.Name) != "" {
					name := strings.TrimSpace(reTrap.ReplaceAllString(obj.Name, ""))
					table.AddTableColumn(name, -1)
				}
			}
		}
	}
	table.AddTableColumn("Action", -1)

	viewFormItemsMap := map[int64][]*Element{}
	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)
		if gp.Name == "Trap Settings" {
			for _, obj := range gp.Child {
				var index int64 = -1
				if subs := reTrap.FindStringSubmatch(obj.Name); len(subs) == 2 {
					index, _ = strconv.ParseInt(subs[1], 10, 32)
				}
				if index < 1 {
					continue
				}
				rowIndex := index - 1
				name := strings.TrimSpace(reTrap.ReplaceAllString(obj.Name, ""))
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elem.SetName(name)
				table.MustSetTableRowData(rowIndex, name, elem)

				viewFormItems, ok := viewFormItemsMap[index]
				if !ok {
					viewFormItems = []*Element{}
				}
				viewFormItems = append(viewFormItems, elem)
				viewFormItemsMap[index] = viewFormItems
			}
		}
	}
	for index, items := range viewFormItemsMap {
		rowIndex := index - 1
		viewPage := NewPageWithSetParameterValuesFormItems(fmt.Sprintf("Trap %v", index), items...)
		table.SetTableRowClickAction(rowIndex, NewViewPageAction(viewPage))
		table.MustSetTableRowData(rowIndex, "Trap", NewLabel("Trap", fmt.Sprintf("%v", index)))
	}

	for rowIndex, row := range table.Data {
		deleteItems := []*Element{}
		for k, v := range row {
			item := *v
			if k == "ip_address" {
				item.SetValue("0.0.0.0")
				deleteItems = append(deleteItems, &item)
			}
		}
		for k, v := range row {
			item := *v
			if k == "security_engine_id" {
				item.SetValue("80000523010A0703")
				deleteItems = append(deleteItems, &item)
			}
		}
		action := NewSetParameterValuesAction(deleteItems...)
		table.MustSetTableRowData(int64(rowIndex), "Action", NewButtonWithAction("Delete", "Delete", "text", action))
	}

	addButton := NewButtonWithAction("Add", "Add", "primary", NewAction("AddItem"))
	table.SetAction("toolbar", NewToolbarWithItems("", addButton))
	table.SetStyle("invalidKey", "ip_address")
	table.SetStyle("invalidValue", "0.0.0.0")
	table.SetStyle("viewKey", "trap")
	return table, nil
}
func MakePageForLanConnectivity(pg *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	managementSetingsItems := []*Element{}
	nmsSettingsItems := []*Element{}
	trapSettingsItems := []*Element{}

	sftpSettingsItems := []*Element{}
	ntpSettingsItems := []*Element{}
	loginSettingsItems := []*Element{}
	consoleInterfaceSettingsItems := []*Element{}
	consoleConfigSettingsItems := []*Element{}
	wanInterfaceSettingsItems := []*Element{}

	elementIdentificationItems := []*Element{}
	dateAndTimeItems := []*Element{}

	hasOpticalModuleTable := false
	table := NewTable("Optical Module List")
	table.AddTableColumnUnique("Optical Module", -1)
	re := regexp.MustCompile(`OP\s*(\d+|P|S|M)`)
	viewFormItemsMap := map[string][]any{}

	for _, gp := range pg.Child {
		groupPath := append([]string{}, path...)
		groupPath = append(groupPath, pg.Name, gp.Name)

		switch gp.Name {
		case "CONSOLE IP Settings", "CONSOLE IP Setting":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				switch elem.Name {
				case "Protocol":
					elem.SetName("Management Protocol")
					managementSetingsItems = append(managementSetingsItems, elem)
				case "Device Recv Port(UDP)":
					elem.SetName("UDP Recv Port")
					managementSetingsItems = append(managementSetingsItems, elem)
				case "Heartbeat Clock":
					elem.SetName("Heartbeat Interval")
					managementSetingsItems = append(managementSetingsItems, elem)
				case "Primary NMS IP Address", "Secondary NMS IP Address":
					nmsSettingsItems = append(nmsSettingsItems, elem)
				case "Primary NMS Port Number", "Secondary NMS Port Number":
					nmsSettingsItems = append(nmsSettingsItems, elem)
				case "Trap IP Address1", "Trap IP Address2", "Trap Port":
					trapSettingsItems = append(trapSettingsItems, elem)
				case "Device IP Addr1(CONSOLE)", "Device IP Address 1(CONSOLE)":
					elem.SetName("IP Address")
					consoleInterfaceSettingsItems = append(consoleInterfaceSettingsItems, elem)
				case "Subnet Mask1(CONSOLE)":
					elem.SetName("Subnet Mask")
					consoleInterfaceSettingsItems = append(consoleInterfaceSettingsItems, elem)
				case "Default Gateway1(CONSOLE)":
					elem.SetName("Default Gateway")
					consoleInterfaceSettingsItems = append(consoleInterfaceSettingsItems, elem)
				case "Device IP Addr2(NMS)", "Device IP Address 2(NMS)":
					elem.SetName("IP Address")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				case "Subnet Mask2(NMS)":
					elem.SetName("Subnet Mask")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				case "Default Gateway2(NMS)":
					elem.SetName("Default Gateway")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				case "Device IP Addr(IPv6)", "Device IP Address(IPv6)":
					elem.SetName("IPv6 Address")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				case "Device IPv6 prefixlen(NMS)":
					elem.SetName("IPv6 Prefix Length")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				case "Device IPv6 Default Gateway(NMS)":
					elem.SetName("IPv6 Default Gateway")
					wanInterfaceSettingsItems = append(wanInterfaceSettingsItems, elem)
				default:
					logrus.Errorf("unknow param for %v, %v", gp.Name, obj.Name)
				}
			}
		case "SFTP Settings":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				switch elem.Name {
				case "Server IP Address (SFTP)":
					elem.SetName("Server IP Address")
					sftpSettingsItems = append(sftpSettingsItems, elem)
				case "Server Port Number (SFTP)":
					elem.SetName("Server Port Number")
					sftpSettingsItems = append(sftpSettingsItems, elem)
				case "SFTP Account Username":
					elem.SetName("Account Username")
					sftpSettingsItems = append(sftpSettingsItems, elem)
				case "SFTP Account Password":
					elem.SetName("Account Password")
					sftpSettingsItems = append(sftpSettingsItems, elem)
				case "Firmware Upgrade Filepath", "Firmware Upgrade Filename":
					sftpSettingsItems = append(sftpSettingsItems, elem)
				case "SFTP File Transfer Control":
					sftpSettingsItems = append(sftpSettingsItems, elem)
				default:
					logrus.Errorf("unknow param for %v, %v", gp.Name, obj.Name)
				}
			}
		case "NTP Settings", "NTP":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				ntpSettingsItems = append(ntpSettingsItems, elem)
			}
		case "OMT Logout":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				loginSettingsItems = append(loginSettingsItems, elem)
			}
		case "Local Console Port Control", "Local Debug Port Control":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				consoleConfigSettingsItems = append(consoleConfigSettingsItems, elem)
			}
		case "Date And Time":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				dateAndTimeItems = append(dateAndTimeItems, elem)
			}
		case "Element Identification":
			for _, obj := range gp.Child {
				elem := MustMakeElementFromParam(obj, info, groupPath)
				elementIdentificationItems = append(elementIdentificationItems, elem)
			}
		case "Optical Module Serial Number And Vendor Name", "Optical Module Tx Power And Rx Power", "OP Transceiver Serial Number And Vendor Name":
			hasOpticalModuleTable = true
			if gp.Name == "Optical Module Serial Number And Vendor Name" || gp.Name == "OP Transceiver Serial Number And Vendor Name" {
				table.AddTableColumn("Serial Number", -1)
				table.AddTableColumn("Vendor Name", -1)
			} else if gp.Name == "Optical Module Tx Power And Rx Power" {
				table.AddTableColumn("Tx Power", -1)
				table.AddTableColumn("Rx Power", -1)
			}
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
				if gp.Name == "Optical Module Serial Number And Vendor Name" || gp.Name == "OP Transceiver Serial Number And Vendor Name" {
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
		default:
			logrus.Errorf("unknow group for lan connectivity, %v", gp.Name)
		}
	}

	formItems := []*Element{}
	if len(wanInterfaceSettingsItems) > 0 {
		formItems = append(formItems, NewLabel("WAN Interface", "WAN Interface"))
		formItems = append(formItems, wanInterfaceSettingsItems...)
	}
	if len(consoleInterfaceSettingsItems) > 0 {
		formItems = append(formItems, NewLabel("Console Interface", "Console Interface"))
		formItems = append(formItems, consoleInterfaceSettingsItems...)
	}
	if len(consoleConfigSettingsItems) > 0 {
		formItems = append(formItems, NewLabel("Console Settings", "Console Settings"))
		formItems = append(formItems, consoleConfigSettingsItems...)
	}
	tabsLayout := NewTabsLayout(pg.Name, "left")
	if tabPage := NewPageWithLayouts("Interfaces", NewSingleColLayoutWithItems(
		NewSetParameterValuesForm("Interface Settings", formItems...),
	)); tabPage != nil {
		tabsLayout.Items = append(tabsLayout.Items, tabPage)
	}

	formItems2 := []*Element{}
	if len(managementSetingsItems) > 0 {
		formItems2 = append(formItems2, NewLabel("Management Settings", "Management Settings"))
		formItems2 = append(formItems2, managementSetingsItems...)
	}
	if len(trapSettingsItems) > 0 {
		formItems2 = append(formItems2, NewLabel("Trap Settings", "Trap Settings"))
		formItems2 = append(formItems2, trapSettingsItems...)
	}
	if len(nmsSettingsItems) > 0 {
		formItems2 = append(formItems2, NewLabel("NMS Settings", "NMS Settings"))
		formItems2 = append(formItems2, nmsSettingsItems...)
	}
	if tabPage := NewPageWithLayouts("Management", NewSingleColLayoutWithItems(
		NewSetParameterValuesForm("Management Settings", formItems2...),
	)); tabPage != nil {
		tabsLayout.Items = append(tabsLayout.Items, tabPage)
	}
	if len(sftpSettingsItems) > 0 {
		if tabPage := NewPageWithLayouts("Remote Upgrade", NewSingleColLayoutWithItems(
			NewSetParameterValuesForm("SFTP Settings", sftpSettingsItems...),
		)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	if len(elementIdentificationItems) > 0 {
		if tabPage := NewPageWithLayouts("Element Indentification", NewSingleColLayoutWithItems(
			NewSetParameterValuesForm("Element Indentification", elementIdentificationItems...),
		)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	if hasOpticalModuleTable {
		for indexText, _ := range viewFormItemsMap {
			index := getOpticalIndex(indexText)
			rowIndex := index - 1
			table.MustSetTableRowData(rowIndex, "Optical Module", NewLabel("Optical Module", fmt.Sprintf("OP %v", indexText)))
		}

		if tabPage := NewPageWithLayouts("Optical Module", NewSingleColLayoutWithItems(table)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	if len(dateAndTimeItems) > 0 {
		if tabPage := NewPageWithLayouts("Time", NewSingleColLayoutWithItems(
			NewSetParameterValuesForm("Date and Time", dateAndTimeItems...),
		)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	if len(ntpSettingsItems) > 0 {
		if tabPage := NewPageWithLayouts("NTP", NewSingleColLayoutWithItems(
			NewSetParameterValuesForm("NTP Settings", ntpSettingsItems...),
		)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	if len(loginSettingsItems) > 0 {
		if tabPage := NewPageWithLayouts("Login", NewSingleColLayoutWithItems(
			NewSetParameterValuesForm("Login Settings", loginSettingsItems...),
		)); tabPage != nil {
			tabsLayout.Items = append(tabsLayout.Items, tabPage)
		}
	}

	page := NewPageWithLayouts("Network Configuration", tabsLayout)
	return page, nil
}
