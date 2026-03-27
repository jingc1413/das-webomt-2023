package layout

import (
	"gomt/core/model"

	"github.com/pkg/errors"
)

type DeviceInfo struct {
	Device *model.ProductDefine
	Params model.ParameterDefines
	Alarms model.AlarmDefines
}

func MakeApplication(name string, info DeviceInfo) (*Element, error) {
	app := NewApplication(name)
	root := model.GetParameterTree(info.Params)
	for _, m := range root.Child {
		module, err := MakeApplicationModule(m, info)
		if err != nil {
			return nil, errors.Wrapf(err, "make module for %v", m.Name)
		}
		app.Items = append(app.Items, module)
	}
	app2, err := TransferApplication(app, info)
	if err != nil {
		return nil, errors.Wrap(err, "transfer")
	}

	return app2, nil
}

func MakeApplicationModule(m *model.ParameterTreeNode, info DeviceInfo) (*Element, error) {
	path := []string{m.Name}

	if m.Name == "Carrier Config" {
		if info.Device.ProductTypeName == "RU" {
			return MakeModuleForRuCarrierConfig(m, info, path)
		} else if info.Device.ProductTypeName == "AU" || info.Device.ProductTypeName == "SAU" {
			return MakeModuleForAuCarrierConfig(m, info, path)
		}
	} else if m.Name == "Carrier Power" {
		return MakeModuleForCarrierPower(m, info, path)
	} else if m.Name == "Channel Config" {
		return MakeModuleForChannelConfig(m, info, path)
	} else if m.Name == "Channel Power" {
		return MakeModuleForChannelPower(m, info, path)
	} else if m.Name == "Management" {
		return MakeModuleForCapacityAllocationManagement(m, info, path)
	} else if m.Name == "Combiners" {
		return MakeModuleForCombiners(m, info, path)
	} else if m.Name == "PA" {
		return MakeModuleForPA(m, info, path)
	} else if m.Name == "Alarms" {
		return MakeModuleForAlarms(m, info)
	} else if m.Name == "Maintenance" {
		return MakeModuleForMaintenance(m, info, path)
	} else if m.Name == "Small-Signal" {
		return MakeModuleForSmallSignal(m, info, path)
	}
	module := NewModule(m.Name)
	for _, pg := range m.Child {
		page, err := MakePage(m, pg, info)
		if err != nil {
			return nil, errors.Wrapf(err, "make page layout for %v", pg.Name)
		}
		if page != nil {
			addPageToModule(page, module)
		}
	}
	return module, nil
}

func MakePage(
	m *model.ParameterTreeNode,
	pg *model.ParameterTreeNode,
	info DeviceInfo,
) (*Element, error) {
	return MakePageWithExcludes(m, pg, info, []string{})
}

func MakePageWithExcludes(
	m *model.ParameterTreeNode,
	pg *model.ParameterTreeNode,
	info DeviceInfo,
	excludeGroups []string,
) (*Element, error) {
	path := []string{m.Name}
	if pg.Name == "Band Configuration" {
		if info.Device.ProductTypeName == "RU" {
			pg2 := m.GetChild("Radio Signal Information")
			return MakePageForRuBandConfiguration(pg, pg2, info, path)
		}
	} else if pg.Name == "Radio Signal Information" || pg.Name == "Input Signal Information" {
		if info.Device.ProductTypeName == "AU" || info.Device.ProductTypeName == "SAU" {
			return MakePageForAuRadioSignalInformation(pg, info, path)
		} else if info.Device.ProductTypeName == "RU" {
			return MakePageForRuRadioSignalInformation(pg, info, path)
		}
	} else if pg.Name == "Radio Interface Modules" || pg.Name == "Input Module Information" {
		return MakePageForAuRadioInterfaceModule(pg, info, path)
	} else if pg.Name == "Optical Module Information" {
		return MakePageForOpticalModuleInformation(pg, info, path)
	} else if pg.Name == "LAN Connectivity" || pg.Name == "LAN Configuration" {
		return MakePageForLanConnectivity(pg, info, path)
	} else if pg.Name == "SNMP Configuration" {
		pg2 := m.GetChild("SNMP User Info")
		return MakePageForSNMPConfiguration(pg, pg2, info, path)
	} else if pg.Name == "SNMP User Info" {
		return nil, nil
	} else if pg.Name == "Factory Command" {
		return MakePageForFactoryCommand(pg, info, path)
	} else if pg.Name == "Address Interface" {
		return MakePageForAddressInterface(pg, info, path)
	}

	tabsLayout := NewTabsLayout(pg.Name, "left")
	for _, gp := range pg.Child {
		excluded := false
		for _, excludeGroup := range excludeGroups {
			if gp.Name == excludeGroup {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		page, err := MakePageForPageGroup(pg, gp, info, path)
		if err != nil {
			return nil, errors.Wrapf(err, "make page for %v", gp.Name)
		}
		if page != nil {
			tabsLayout.Items = append(tabsLayout.Items, page)
		}
	}
	page := NewPageWithLayouts(pg.Name, tabsLayout)
	return page, nil
}

func MakePageForPageGroup(pg *model.ParameterTreeNode, gp *model.ParameterTreeNode, info DeviceInfo, path []string) (*Element, error) {
	items := []*Element{}
	groupPath := append([]string{}, path...)
	groupPath = append(groupPath, pg.Name, gp.Name)
	for _, obj := range gp.Child {
		elem := MustMakeElementFromParam(obj, info, groupPath)
		if elem == nil {
			continue
		}
		items = append(items, elem)
	}
	name := gp.Name
	if gp.Name == "" {
		name = pg.Name
	}

	form := NewSetParameterValuesForm(name, items...)
	if form.Key == "dual_sfp_mode" {
		form.SetStyle("confirmTitle", "WARNING")
		form.SetStyle("confirmMessage",
			info.Device.DeviceTypeName+" will reboot after switch transmission mode and all carrier configurations will be reset. It is recommended to export the current carrier configurations for backup before switching.",
		)
		form.SetStyle("confirmType", "warning")
	} else if form.Key == "ru_radio_module_signal_test_active" {
		form.SetStyle("confirmTitle", "WARNING")
		form.SetStyle("confirmMessage",
			"Make sure that all RUs have been connected to antenna or terminated before enabling CW Test Signal function. Otherwise the PA of RU will be damaged if it is not terminated.",
		)
		form.SetStyle("confirmType", "warning")
	} else if form.Key == "all_amplifier_module_signal_test_active" {
		form.SetStyle("confirmTitle", "WARNING")
		form.SetStyle("confirmMessage",
			"Please ensure that all active RU DL ports are connected to their antennas, or properly terminated, before enabling the CW Test Signal feature. If not power amplifier damage may occur.",
		)
		form.SetStyle("confirmType", "warning")
	}
	if form.Key == "system_delay" || form.Key == "figer_delay" {
		form.FilterItemsByKeys([]string{"flatnesscsv"})
	}

	page := NewPageWithLayouts(name, NewSingleColLayoutWithItems(form))
	return page, nil
}
