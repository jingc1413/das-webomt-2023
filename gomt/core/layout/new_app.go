package layout

import (
	"fmt"
	"sort"
	"strings"
)

func TransferApplication(source *Element, info DeviceInfo) (*Element, error) {
	pageMap := GetPageMap(source, "Page")
	keys := []string{}
	for _, k := range pageMap.Keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	usedKeys := []string{}
	app := NewApplication(source.Name)

	if module := NewModule("Overview"); module != nil {
		module.SetStyle("icon", "Monitor")
		module.Items = append(module.Items, NewPage("DAS Topo"))
		if page := NewPage("Element Information"); page != nil {
			if page2 := pageMap.GetPage("settings.overview.element_identification"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.element_identification")
				addPageToPage(page2, page)
			}
			if page2 := pageMap.GetPage("settings.overview.date_and_time"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.date_and_time")
				items := []*Element{}
				if form := page2.GetItemByType("Form"); form != nil {
					for _, item := range form.Items {
						tmp := *item
						if tmp.Key == "system_date_and_time" {
							tmp.Style = nil
							tmp.SetStyle("input", "default")
							tmp.SetStyle("readonly", true)
							tmp.Access = "ro"
						}
						items = append(items, &tmp)
					}
				}
				if len(items) > 0 {
					_page := NewPageWithSetParameterValuesFormItems(page2.Name, items...)
					addPageToPage(_page, page)
				}
			}
			if page2 := pageMap.GetPage("settings.overview.temperature"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.temperature")
				addPageToPage(page2, page)
			}
			if page2 := pageMap.GetPage("settings.overview.general"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.general")
				addPageToPage(page2, page)
			}
			if page2 := pageMap.GetPage("settings.overview.reset"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.reset")
				addPageToPage(page2, page)
			}
			addPageToModule(page, module)
		}

		if page := pageMap.GetPage("settings.optical_module_information"); page != nil {
			usedKeys = append(usedKeys, "settings.optical_module_information")
			page.SetName("Optical Module")
			addPageToModule(page, module)
		}
		addModuleToApp(module, app)
	}

	// if module := NewModule("Radio Signal Configuration"); module != nil {
	// 	module.SetStyle("icon", "Operation")

	// 	if page := pageMap.GetPage("settings.band_configuration"); page != nil {
	// 		usedKeys = append(usedKeys, "settings.band_configuration")
	// 		page.SetName("Band Configuration")
	// 		if table := page.GetItemByType("Table"); table != nil {
	// 			if page2 := pageMap.GetPage("settings.radio_signal_information"); page2 != nil {
	// 				usedKeys = append(usedKeys, "settings.radio_signal_information")
	// 				// if table2 := page2.GetItemByType("Table"); table2 != nil {
	// 				// 	mergeTables(table, table2, "Radio Module")
	// 				// }
	// 			}
	// 		}
	// 		addPageToModule(page, module)
	// 	} else {
	// 		if page := pageMap.GetPage("settings.radio_signal_information"); page != nil {
	// 			usedKeys = append(usedKeys, "settings.radio_signal_information")
	// 			page.SetName("Radio Signal")
	// 			addPageToModule(page, module)
	// 		}
	// 		if page := pageMap.GetPage("settings.input_signal_information"); page != nil {
	// 			usedKeys = append(usedKeys, "settings.input_signal_information")
	// 			page.SetName("Input Signal")
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := NewPage("Radio Configuration"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		if page2 := pageMap.GetPage("settings.radio_interface_modules.general"); page2 != nil {
	// 			usedKeys = append(usedKeys, "settings.radio_interface_modules.general")
	// 			page2.SetName("Radio Module")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, true) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := pageMap.GetPage("settings.radio_interface_modules.radio_interface"); page != nil {
	// 		usedKeys = append(usedKeys, "settings.radio_interface_modules.radio_interface")
	// 		page.SetName("Radio Interface")
	// 		addPageToModule(page, module)
	// 	}

	// 	if page := NewPage("Input Configuration"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		if page2 := pageMap.GetPage("settings.input_module_information.general"); page2 != nil {
	// 			usedKeys = append(usedKeys, "settings.input_module_information.general")
	// 			page2.SetName("Input Module")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, true) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := pageMap.GetPage("settings.input_module_information.input_interface"); page != nil {
	// 		usedKeys = append(usedKeys, "settings.input_module_information.input_interface")
	// 		page.SetName("Input Interface")
	// 		addPageToModule(page, module)
	// 	}
	// 	addModuleToApp(module, app)
	// }

	// if module := NewModule("Digital Signal Configuration"); module != nil {
	// 	module.SetStyle("icon", "Cpu")

	// 	hasCarrierConfiguration := false
	// 	if page := NewPage("Carrier Configuration"); page != nil {
	// 		tabs := NewTabsLayout("", "top")
	// 		if page2 := pageMap.GetPage("carrier_config.system_signal.carrier_list"); page2 != nil {
	// 			usedKeys = append(usedKeys, "carrier_config.system_signal.carrier_list")
	// 			page2.SetName("Carrier List")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("carrier_config.system_signal.channel_list"); page2 != nil {
	// 			usedKeys = append(usedKeys, "carrier_config.system_signal.channel_list")
	// 			page2.SetName("Carrier List")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("carrier_power.power_configuration"); page2 != nil {
	// 			usedKeys = append(usedKeys, "carrier_power.power_configuration")
	// 			page2.SetName("Carrier Power")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}

	// 		if page2 := pageMap.GetPage("channel_config.system_signal.channel_list"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.channel_list")
	// 			page2.SetName("Channel List")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("channel_power.power_configuration"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_power.power_configuration")
	// 			page2.SetName("Carrier Power")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			hasCarrierConfiguration = true
	// 			addPageToModule(page, module)
	// 		}
	// 	}

	// 	carrierConfigurationName := "Carrier Configuration"
	// 	if hasCarrierConfiguration {
	// 		carrierConfigurationName = "Carrier Configuration2"
	// 	}
	// 	if page := NewPage(carrierConfigurationName); page != nil {
	// 		tabs := NewTabsLayout("", "top")

	// 		if page2 := pageMap.GetPage("channel_config.system_signal.carrier_configuration"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.carrier_configuration")
	// 			page2.SetName("Carrier Configuration")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("channel_config.system_signal.carrier_status"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.carrier_status")
	// 			page2.SetName("Carrier Status")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("channel_config.system_signal.power_sharing"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.power_sharing")
	// 			page2.SetName("Power Configuration")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("channel_config.system_signal.input_module_list"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.input_module_list")
	// 			page2.SetName("Input Module")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := NewPage("Signal Configuration"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		if page2 := pageMap.GetPage("carrier_config.system_signal.dual_sfp_mode"); page2 != nil {
	// 			usedKeys = append(usedKeys, "carrier_config.system_signal.dual_sfp_mode")
	// 			page2.SetName("Dual SFP Mode")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("carrier_config.system_signal.system_delay"); page2 != nil {
	// 			usedKeys = append(usedKeys, "carrier_config.system_signal.system_delay")
	// 			page2.SetName("System Delay")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}

	// 		if page2 := pageMap.GetPage("channel_config.system_signal.dual_sfp_mode"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.dual_sfp_mode")
	// 			page2.SetName("Dual SFP Mode")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("channel_config.system_signal.fiber_delay"); page2 != nil {
	// 			usedKeys = append(usedKeys, "channel_config.system_signal.fiber_delay")
	// 			page2.SetName("System Delay")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := NewPage("TDD Configuration"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		if page2 := pageMap.GetPage("settings.tdd_configuration.5g_tdd"); page2 != nil {
	// 			usedKeys = append(usedKeys, "settings.tdd_configuration.5g_tdd")
	// 			page2.SetName("5G TDD")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("settings.tdd_configuration.4g_tdd"); page2 != nil {
	// 			usedKeys = append(usedKeys, "settings.tdd_configuration.4g_tdd")
	// 			page2.SetName("4G TDD")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := NewPage("Capacity Allocation"); page != nil {
	// 		tabs := NewTabsLayout("", "top")
	// 		if pageMap2 := pageMap.GetPageMapWithPrefix("capacity_allocation.service_configuration."); len(pageMap2.Keys) > 0 {
	// 			for _, k := range pageMap2.Keys {
	// 				v := pageMap2.Pages[k]
	// 				usedKeys = append(usedKeys, k)
	// 				tabs.Items = append(tabs.Items, v)
	// 			}
	// 		}
	// 		if page2 := pageMap.GetPage("capacity_allocation.rf_module_mapping"); page2 != nil {
	// 			usedKeys = append(usedKeys, "capacity_allocation.rf_module_mapping")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}
	// 		if page2 := pageMap.GetPage("capacity_allocation.service_schedule"); page2 != nil {
	// 			usedKeys = append(usedKeys, "capacity_allocation.service_schedule")
	// 			tabs.Items = append(tabs.Items, page2)
	// 		}

	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}

	// 	addModuleToApp(module, app)
	// }

	if module := NewModule("Signal Configuration"); module != nil {
		module.SetStyle("icon", "Operation")

		if page := pageMap.GetPage("settings.band_configuration"); page != nil {
			usedKeys = append(usedKeys, "settings.band_configuration")
			page.SetName("Band Configuration")

			if table := page.GetItemByType("Table"); table != nil {
				if page2 := pageMap.GetPage("settings.radio_signal_information"); page2 != nil {
					usedKeys = append(usedKeys, "settings.radio_signal_information")
					// if table2 := page2.GetItemByType("Table"); table2 != nil {
					// 	mergeTables(table, table2, "Radio Module")
					// }
				}
			}
			addPageToModule(page, module)
		} else {
			if page := NewPage("Radio Module"); page != nil {
				tabs := NewTabsLayout("", "top")
				if page2 := pageMap.GetPage("settings.radio_signal_information"); page2 != nil {
					usedKeys = append(usedKeys, "settings.radio_signal_information")
					page2.SetName("Radio Signal")
					tabs.Items = append(tabs.Items, page2)
				}
				if page2 := pageMap.GetPage("settings.radio_interface_modules.radio_interface"); page2 != nil {
					usedKeys = append(usedKeys, "settings.radio_interface_modules.radio_interface")
					page2.SetName("Radio Interface")
					tabs.Items = append(tabs.Items, page2)
				}
				if page2 := pageMap.GetPage("settings.radio_interface_modules.general"); page2 != nil {
					usedKeys = append(usedKeys, "settings.radio_interface_modules.general")
					page2.SetName("General Setting")
					tabs.Items = append(tabs.Items, page2)
				}
				if addTabsToPage(tabs, page, true) {
					addPageToModule(page, module)
				}
			}
			if page := NewPage("Input Module"); page != nil {
				tabs := NewTabsLayout("", "top")

				if page2 := pageMap.GetPage("settings.input_signal_information"); page2 != nil {
					usedKeys = append(usedKeys, "settings.input_signal_information")
					page2.SetName("Input Signal")
					tabs.Items = append(tabs.Items, page2)
				}
				if page2 := pageMap.GetPage("settings.input_module_information.input_interface"); page2 != nil {
					usedKeys = append(usedKeys, "settings.input_module_information.input_interface")
					page2.SetName("Input Interface")
					tabs.Items = append(tabs.Items, page2)
				}
				if page2 := pageMap.GetPage("settings.input_module_information.general"); page2 != nil {
					usedKeys = append(usedKeys, "settings.input_module_information.general")
					page2.SetName("General Setting")
					tabs.Items = append(tabs.Items, page2)
				}

				if addTabsToPage(tabs, page, true) {
					addPageToModule(page, module)
				}
			}
		}

		hasCarrierConfiguration := false
		if page := NewPage("Carrier Configuration"); page != nil {
			tabs := NewTabsLayout("", "top")
			if page2 := pageMap.GetPage("carrier_config.system_signal.carrier_list"); page2 != nil {
				usedKeys = append(usedKeys, "carrier_config.system_signal.carrier_list")
				page2.SetName("Carrier List")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("carrier_config.system_signal.channel_list"); page2 != nil {
				usedKeys = append(usedKeys, "carrier_config.system_signal.channel_list")
				page2.SetName("Carrier List")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("carrier_power.power_configuration"); page2 != nil {
				usedKeys = append(usedKeys, "carrier_power.power_configuration")
				page2.SetName("Carrier Power")
				tabs.Items = append(tabs.Items, page2)
			}

			if page2 := pageMap.GetPage("channel_config.system_signal.channel_list"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.channel_list")
				page2.SetName("Channel List")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("channel_power.power_configuration"); page2 != nil {
				usedKeys = append(usedKeys, "channel_power.power_configuration")
				page2.SetName("Carrier Power")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				hasCarrierConfiguration = true
				addPageToModule(page, module)
			}
		}

		carrierConfigurationName := "Carrier Configuration"
		if hasCarrierConfiguration {
			carrierConfigurationName = "Carrier Configuration2"
		}
		if page := NewPage(carrierConfigurationName); page != nil {
			tabs := NewTabsLayout("", "top")

			if page2 := pageMap.GetPage("channel_config.system_signal.carrier_configuration"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.carrier_configuration")
				page2.SetName("Carrier Configuration")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("channel_config.system_signal.carrier_status"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.carrier_status")
				page2.SetName("Carrier Status")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("channel_config.system_signal.power_sharing"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.power_sharing")
				page2.SetName("Power Configuration")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("channel_config.system_signal.input_module_list"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.input_module_list")
				page2.SetName("Input Module")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Signal Configuration"); page != nil {
			tabs := NewTabsLayout("", "left")
			if page2 := pageMap.GetPage("carrier_config.system_signal.dual_sfp_mode"); page2 != nil {
				usedKeys = append(usedKeys, "carrier_config.system_signal.dual_sfp_mode")
				page2.SetName("Dual SFP Mode")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("carrier_config.system_signal.system_delay"); page2 != nil {
				usedKeys = append(usedKeys, "carrier_config.system_signal.system_delay")
				page2.SetName("System Delay")
				tabs.Items = append(tabs.Items, page2)
			}

			if page2 := pageMap.GetPage("channel_config.system_signal.dual_sfp_mode"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.dual_sfp_mode")
				page2.SetName("Dual SFP Mode")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("channel_config.system_signal.fiber_delay"); page2 != nil {
				usedKeys = append(usedKeys, "channel_config.system_signal.fiber_delay")
				page2.SetName("System Delay")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.tdd_configuration.5g_tdd"); page2 != nil {
				usedKeys = append(usedKeys, "settings.tdd_configuration.5g_tdd")
				page2.SetName("5G TDD")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.tdd_configuration.4g_tdd"); page2 != nil {
				usedKeys = append(usedKeys, "settings.tdd_configuration.4g_tdd")
				page2.SetName("4G TDD")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Capacity Allocation"); page != nil {
			tabs := NewTabsLayout("", "top")
			if pageMap2 := pageMap.GetPageMapWithPrefix("capacity_allocation.service_configuration."); len(pageMap2.Keys) > 0 {
				for _, k := range pageMap2.Keys {
					v := pageMap2.Pages[k]
					usedKeys = append(usedKeys, k)
					tabs.Items = append(tabs.Items, v)
				}
			}
			if page2 := pageMap.GetPage("capacity_allocation.rf_module_mapping"); page2 != nil {
				usedKeys = append(usedKeys, "capacity_allocation.rf_module_mapping")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("capacity_allocation.service_schedule"); page2 != nil {
				usedKeys = append(usedKeys, "capacity_allocation.service_schedule")
				tabs.Items = append(tabs.Items, page2)
			}

			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}

		addModuleToApp(module, app)
	}

	if module := NewModule("Diagnostics"); module != nil {
		module.SetStyle("icon", "Histogram")
		module.Items = append(module.Items, NewPage("Statistics"))
		if page := NewPage("Network Diagnostic"); page != nil {
			tabs := NewTabsLayout("", "top")
			if page2 := NewPage("Ping Diagnostic"); page2 != nil {
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, true) {
				addPageToModule(page, module)
			}
		}
		module.Items = append(module.Items, NewPage("Spectrum Diagnostic"))
		addModuleToApp(module, app)
	}

	if module := NewModule("System Settings"); module != nil {
		module.SetStyle("icon", "Tools")

		if page := NewPage("Connectivity"); page != nil {
			tabs := NewTabsLayout("", "left")
			if page2 := pageMap.GetPage("settings.network_configuration.interfaces"); page2 != nil {
				usedKeys = append(usedKeys, "settings.network_configuration.interfaces")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.network_configuration.management"); page2 != nil {
				usedKeys = append(usedKeys, "settings.network_configuration.management")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.network_configuration.remote_upgrade"); page2 != nil {
				usedKeys = append(usedKeys, "settings.network_configuration.remote_upgrade")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		// if page := NewPage("Connectivity"); page != nil {
		// 	tabs := NewTabsLayout("", "left")
		// 	if page2 := pageMap.GetPage("settings.lan_connectivity.console_ip_settings"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_connectivity.console_ip_settings")
		// 		page2.SetName("Network")
		// 		tabs.Items = append(tabs.Items, page2)
		// 		if form := page2.GetItemByType("Form"); form != nil {
		// 			items := []*Element{}
		// 			items2 := []*Element{}
		// 			for _, v := range form.Items {
		// 				if strings.Contains(v.Name, "NMS ") {
		// 					items2 = append(items2, v)
		// 				} else {
		// 					items = append(items, v)
		// 				}
		// 			}
		// 			form.Items = items
		// 			if len(items2) > 0 {
		// 				tabs.Items = append(tabs.Items,
		// 					NewPageWithSetParameterValuesFormItems("NMS", items2...),
		// 				)
		// 			}
		// 		}
		// 	}
		// 	if page2 := pageMap.GetPage("settings.lan_connectivity.sftp_settings"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_connectivity.sftp_settings")
		// 		page2.SetName("SFTP Settings")
		// 		tabs.Items = append(tabs.Items, page2)
		// 	}
		// 	if page2 := pageMap.GetPage("settings.lan_connectivity.local_console_port_control"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_connectivity.local_console_port_control")
		// 		page2.SetName("Console Port")
		// 		tabs.Items = append(tabs.Items, page2)
		// 	}

		// 	if page2 := pageMap.GetPage("settings.lan_configuration.console_ip_settings"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_configuration.console_ip_settings")
		// 		page2.SetName("Network")
		// 		tabs.Items = append(tabs.Items, page2)
		// 		if form := page2.GetItemByType("Form"); form != nil {
		// 			items := []*Element{}
		// 			items2 := []*Element{}
		// 			for _, v := range form.Items {
		// 				if strings.Contains(v.Name, "NMS ") {
		// 					items2 = append(items2, v)
		// 				} else {
		// 					items = append(items, v)
		// 				}
		// 			}
		// 			form.Items = items
		// 			if len(items2) > 0 {
		// 				tabs.Items = append(tabs.Items,
		// 					NewPageWithSetParameterValuesFormItems("NMS", items2...),
		// 				)
		// 			}
		// 		}
		// 	}
		// 	if page2 := pageMap.GetPage("settings.lan_configuration.sftp_settings"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_configuration.sftp_settings")
		// 		page2.SetName("SFTP Settings")
		// 		tabs.Items = append(tabs.Items, page2)
		// 	}
		// 	if page2 := pageMap.GetPage("settings.lan_configuration.local_console_port_control"); page2 != nil {
		// 		usedKeys = append(usedKeys, "settings.lan_configuration.local_console_port_control")
		// 		page2.SetName("Console Port")
		// 		tabs.Items = append(tabs.Items, page2)
		// 	}
		// 	if addTabsToPage(tabs, page, false) {
		// 		addPageToModule(page, module)
		// 	}
		// }
		if page := NewPage("SNMP"); page != nil {
			tabs := NewTabsLayout("", "left")
			if page2 := pageMap.GetPage("settings.snmp_configuration.general"); page2 != nil {
				usedKeys = append(usedKeys, "settings.snmp_configuration.general")
				page2.SetName("General")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.snmp_configuration.trap_settings_list"); page2 != nil {
				usedKeys = append(usedKeys, "settings.snmp_configuration.trap_settings_list")
				page2.SetName("Trap Address List")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.snmp_configuration.snmp_user_list"); page2 != nil {
				usedKeys = append(usedKeys, "settings.snmp_configuration.snmp_user_list")
				page2.SetName("USM List")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Alarm"); page != nil {
			tabs := NewTabsLayout("", "top")
			tabs.Items = append(tabs.Items, NewPageWithLayouts("Alarm Logs",
				NewSingleColLayoutWithItems(NewAlarmLogsTable("Alarm Logs"))))

			if page2 := pageMap.GetPage("alarms.element_alarms.device"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.device")
				page2.SetName("Device")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("alarms.element_alarms.optical_module"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.optical_module")
				page2.SetName("Radio Module")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("alarms.element_alarms.radio_module"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.radio_module")
				page2.SetName("Radio Module")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("alarms.element_alarms.radio_interface"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.radio_interface")
				page2.SetName("Radio Interface")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("alarms.element_alarms.external_input"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.external_input")
				page2.SetName("External Input")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("alarms.element_alarms.alarm_indication"); page2 != nil {
				usedKeys = append(usedKeys, "alarms.element_alarms.alarm_indication")
				page2.SetName("Alarm Indication")
				tabs.Items = append(tabs.Items, page2)
			}

			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Logs"); page != nil {
			tabs := NewTabsLayout("", "top")
			tabs.Items = append(tabs.Items, NewPageWithLayouts("Guard Logs",
				NewSingleColLayoutWithItems(NewFilesTable("GuardLog"))))
			tabs.Items = append(tabs.Items, NewPageWithLayouts("Web Logs",
				NewSingleColLayoutWithItems(NewFilesTable("WebLog"))))
			tabs.Items = append(tabs.Items, NewPageWithLayouts("Device Logs",
				NewSingleColLayoutWithItems(NewFilesTable("DeviceLog"))))
			tabs.Items = append(tabs.Items, NewPageWithLayouts("Login Logs",
				NewSingleColLayoutWithItems(NewFilesTable("LoginLog"))))
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPageWithLayouts("Upgrade", NewSingleColLayoutWithItems(NewFilesTable("UpgradeFile"))); page != nil {
			addPageToModule(page, module)
		}
		if page := NewPageWithLayouts("Configuration", NewSingleColLayoutWithItems(NewFilesTable("ConfigFile"))); page != nil {
			addPageToModule(page, module)
		}
		if page := NewPage("Time and Date"); page != nil {
			tabs := NewTabsLayout("", "left")
			if page2 := pageMap.GetPage("settings.overview.date_and_time"); page2 != nil {
				usedKeys = append(usedKeys, "settings.overview.date_and_time")
				page2.SetName("Time and Date")
				tabs.Items = append(tabs.Items, page2)
			}
			if page2 := pageMap.GetPage("settings.network_configuration.ntp"); page2 != nil {
				usedKeys = append(usedKeys, "settings.network_configuration.ntp")
				page2.SetName("NTP")
				tabs.Items = append(tabs.Items, page2)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}

		if info.Device.ProductTypeName == "AU" {
			if page := NewPageWithLayouts("Inventory",
				NewSingleColLayoutWithItems(NewInventoryTable("Inventory"))); page != nil {
				addPageToModule(page, module)
			}
			if page := NewPage("Account"); page != nil {
				tabs := NewTabsLayout("", "left")
				tabs.Items = append(tabs.Items, NewPageWithLayouts("Users",
					NewSingleColLayoutWithItems(NewUsersTable("User"))))
				if page2 := pageMap.GetPage("settings.network_configuration.login"); page2 != nil {
					usedKeys = append(usedKeys, "settings.network_configuration.login")
					page2.SetName("Login Session")
					tabs.Items = append(tabs.Items, page2)
				}

				if addTabsToPage(tabs, page, false) {
					addPageToModule(page, module)
				}
			}
		}

		if page := NewPage("Contact Info"); page != nil {
			tabs := NewTabsLayout("", "left")
			pageMap2 := pageMap.GetPageMapWithPrefix("maintenance.contact_info.")
			for _, k := range pageMap2.Keys {
				v := pageMap2.Pages[k]
				usedKeys = append(usedKeys, k)
				tabs.Items = append(tabs.Items, v)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		addModuleToApp(module, app)
	}

	if module := NewModule("Maintenance"); module != nil {
		module.SetStyle("icon", "OfficeBuilding")
		if page := NewPage("Engineering"); page != nil {
			tabs := NewTabsLayout("", "left")
			pageMap2 := pageMap.GetPageMapWithPrefix("maintenance.engineering.")
			for _, k := range pageMap2.Keys {
				v := pageMap2.Pages[k]
				usedKeys = append(usedKeys, k)
				tabs.Items = append(tabs.Items, v)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := pageMap.GetPage("maintenance.firmware_status.package_info"); page != nil {
			usedKeys = append(usedKeys, "maintenance.firmware_status.package_info")
			// page.SetName("Firmware Information")
			// addPageToModule(page, module)

			page2 := NewPageWithLayouts("Firmware Information", NewSingleColLayoutWithItems(NewFirmwaresTable("Firmware")))
			addPageToModule(page2, module)
		}
		if page := pageMap.GetPage("maintenance.optical_info"); page != nil {
			usedKeys = append(usedKeys, "maintenance.optical_info")
			page.SetName("Optical Module Information")
			addPageToModule(page, module)
		}
		if page := NewPage("Factory"); page != nil {
			tabs := NewTabsLayout("", "left")
			pageMap2 := pageMap.GetPageMapWithPrefix("maintenance.factory_command.")
			for _, k := range pageMap2.Keys {
				v := pageMap2.Pages[k]
				usedKeys = append(usedKeys, k)
				tabs.Items = append(tabs.Items, v)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		addModuleToApp(module, app)
	}

	if module := NewModule("Factory Maintenance"); module != nil {
		module.SetStyle("icon", "Menu")

		if page := pageMap.GetPage("digital_module.address_interface"); page != nil {
			usedKeys = append(usedKeys, "digital_module.address_interface")
			page.SetName("Address Interface")
			addPageToModule(page, module)
		}
		if page := NewPage("Debug Parameter"); page != nil {
			tabs := NewTabsLayout("", "left")
			pageMap2 := pageMap.GetPageMapWithPrefix("digital_module.debug_parameters.")
			for _, k := range pageMap2.Keys {
				v := pageMap2.Pages[k]
				usedKeys = append(usedKeys, k)
				tabs.Items = append(tabs.Items, v)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("MCU Parameter"); page != nil {
			tabs := NewTabsLayout("", "left")
			pageMap2 := pageMap.GetPageMapWithPrefix("digital_module.mcu_parameter.")
			for _, k := range pageMap2.Keys {
				v := pageMap2.Pages[k]
				usedKeys = append(usedKeys, k)
				tabs.Items = append(tabs.Items, v)
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Small Signal"); page != nil {
			tabs := NewTabsLayout("", "left")
			if page2 := NewPage("General"); page2 != nil {
				for _, k := range pageMap.Keys {
					if strings.HasPrefix(k, "small_signal.general.") && len(strings.Split(k, ".")) == 3 {
						v := pageMap.Pages[k]
						usedKeys = append(usedKeys, k)
						page2.Items = append(tabs.Items, v.Items...)
					}
				}
				if len(page2.Items) > 0 {
					tabs.Items = append(tabs.Items, page2)
				}
			}
			for _, k := range pageMap.Keys {
				if strings.HasPrefix(k, "small_signal.") && len(strings.Split(k, ".")) == 2 {
					v := pageMap.Pages[k]
					usedKeys = append(usedKeys, k)
					tabs.Items = append(tabs.Items, v)
				}
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}
		if page := NewPage("Frequency Gain Test"); page != nil {
			for _, k := range pageMap.Keys {
				if strings.HasPrefix(k, "digital_module.freqgaincomp_test.") {
					v := pageMap.Pages[k]
					usedKeys = append(usedKeys, k)
					addPageToPage(v, page)
				}
			}
			addPageToModule(page, module)
		}
		if page := NewPage("Combiner"); page != nil {
			tabs := NewTabsLayout("", "top")
			for _, k := range pageMap.Keys {
				if strings.HasPrefix(k, "combiners.") {
					v := pageMap.Pages[k]
					usedKeys = append(usedKeys, k)
					tabs.Items = append(tabs.Items, v)
				}
			}
			if addTabsToPage(tabs, page, false) {
				addPageToModule(page, module)
			}
		}

		//PA
		{
			// for _, k := range pageMap.Keys {
			// 	if strings.HasPrefix(k, "pa.") && len(strings.Split(k, ".")) == 2 {
			// 		v := pageMap.Pages[k]
			// 		usedKeys = append(usedKeys, k)
			// 		addPageToModule(v, module)
			// 	}
			// }
			if page := NewPage("PA Configuration"); page != nil {
				tabs := NewTabsLayout("", "left")
				for i := 1; i <= 32; i++ {
					k := fmt.Sprintf("pa.pa_configuration.pa_%d", i)
					if page2 := pageMap.GetPage(k); page2 != nil {
						usedKeys = append(usedKeys, k)
						tabs.Items = append(tabs.Items, page2)
					}
				}
				if addTabsToPage(tabs, page, false) {
					addPageToModule(page, module)
				}
			}
			for i := 1; i <= 32; i++ {
				for j := 1; j <= 32; j++ {
					prefix := fmt.Sprintf("pa.pa%d_ch%d.", i, j)
					if page := NewPage(fmt.Sprintf("PA%d CH%d", i, j)); page != nil {
						tabs := NewTabsLayout("", "left")
						for _, k := range pageMap.Keys {
							if strings.HasPrefix(k, prefix) && len(strings.Split(k, ".")) == 3 {
								v := pageMap.Pages[k]
								usedKeys = append(usedKeys, k)
								tabs.Items = append(tabs.Items, v)
							}
						}
						if addTabsToPage(tabs, page, false) {
							addPageToModule(page, module)
						}
					}
				}
			}
			if page := NewPage("PA Common"); page != nil {
				tabs := NewTabsLayout("", "left")
				for _, k := range pageMap.Keys {
					if strings.HasPrefix(k, "pa.pa_common.") && len(strings.Split(k, ".")) == 3 {
						v := pageMap.Pages[k]
						usedKeys = append(usedKeys, k)
						tabs.Items = append(tabs.Items, v)
					}
				}
				if addTabsToPage(tabs, page, false) {
					addPageToModule(page, module)
				}
			}
			if page := NewPage("PA Switch"); page != nil {
				tabs := NewTabsLayout("", "left")
				for _, k := range pageMap.Keys {
					if strings.HasPrefix(k, "pa.pa_switch.") && len(strings.Split(k, ".")) == 3 {
						v := pageMap.Pages[k]
						usedKeys = append(usedKeys, k)
						tabs.Items = append(tabs.Items, v)
					}
				}
				if addTabsToPage(tabs, page, false) {
					addPageToModule(page, module)
				}
			}
		}
		addModuleToApp(module, app)
	}

	// if module := NewModule("Combiner"); module != nil {
	// 	module.SetStyle("icon", "Document")

	// 	if page := NewPage("Combiner"); page != nil {
	// 		tabs := NewTabsLayout("", "top")
	// 		for _, k := range pageMap.Keys {
	// 			if strings.HasPrefix(k, "combiners.") {
	// 				v := pageMap.Pages[k]
	// 				usedKeys = append(usedKeys, k)
	// 				tabs.Items = append(tabs.Items, v)
	// 			}
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	addModuleToApp(module, app)
	// }

	// if module := NewModule("PA"); module != nil {
	// 	module.SetStyle("icon", "Document")

	// 	for _, k := range pageMap.Keys {
	// 		if strings.HasPrefix(k, "pa.") && len(strings.Split(k, ".")) == 2 {
	// 			v := pageMap.Pages[k]
	// 			usedKeys = append(usedKeys, k)
	// 			addPageToModule(v, module)
	// 		}
	// 	}
	// 	for i := 1; i <= 32; i++ {
	// 		prefix := fmt.Sprintf("pa.pa_ch%d.", i)
	// 		if page := NewPage(fmt.Sprintf("PA CH%d", i)); page != nil {
	// 			tabs := NewTabsLayout("", "left")
	// 			for _, k := range pageMap.Keys {
	// 				if strings.HasPrefix(k, prefix) && len(strings.Split(k, ".")) == 3 {
	// 					v := pageMap.Pages[k]
	// 					usedKeys = append(usedKeys, k)
	// 					tabs.Items = append(tabs.Items, v)
	// 				}
	// 			}
	// 			if addTabsToPage(tabs, page, false) {
	// 				addPageToModule(page, module)
	// 			}
	// 		}
	// 	}
	// 	for i := 1; i <= 32; i++ {
	// 		for j := 1; j <= 32; j++ {
	// 			prefix := fmt.Sprintf("pa.pa%d_ch%d.", i, j)
	// 			if page := NewPage(fmt.Sprintf("PA%d CH%d", i, j)); page != nil {
	// 				tabs := NewTabsLayout("", "left")
	// 				for _, k := range pageMap.Keys {
	// 					if strings.HasPrefix(k, prefix) && len(strings.Split(k, ".")) == 3 {
	// 						v := pageMap.Pages[k]
	// 						usedKeys = append(usedKeys, k)
	// 						tabs.Items = append(tabs.Items, v)
	// 					}
	// 				}
	// 				if addTabsToPage(tabs, page, false) {
	// 					addPageToModule(page, module)
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if page := NewPage("PA Common"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		for _, k := range pageMap.Keys {
	// 			if strings.HasPrefix(k, "pa.pa_common.") && len(strings.Split(k, ".")) == 3 {
	// 				v := pageMap.Pages[k]
	// 				usedKeys = append(usedKeys, k)
	// 				tabs.Items = append(tabs.Items, v)
	// 			}
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	if page := NewPage("PA Switch"); page != nil {
	// 		tabs := NewTabsLayout("", "left")
	// 		for _, k := range pageMap.Keys {
	// 			if strings.HasPrefix(k, "pa.pa_switch.") && len(strings.Split(k, ".")) == 3 {
	// 				v := pageMap.Pages[k]
	// 				usedKeys = append(usedKeys, k)
	// 				tabs.Items = append(tabs.Items, v)
	// 			}
	// 		}
	// 		if addTabsToPage(tabs, page, false) {
	// 			addPageToModule(page, module)
	// 		}
	// 	}
	// 	addModuleToApp(module, app)
	// }

	for _, v := range keys {
		used := false
		for _, v2 := range usedKeys {
			if v == v2 {
				used = true
				break
			}
		}
		if strings.HasPrefix(v, "settings.engineering_cmd") {

		} else if !used {
			fmt.Printf("!!! %v\n", v)
		} else {
			// fmt.Printf("%v\n", v)
		}
	}
	return app, nil
}

func addTabsToPage(tabs *Element, page *Element, keepTabs bool) bool {
	num := len(tabs.Items)
	size := len(page.Items)
	if num < 1 {
		return false
	}

	if size == 0 {
		if num == 1 && !keepTabs {
			page.Items = tabs.Items[0].Items
		} else {
			page.Items = append(page.Items, tabs)
		}
	} else {
		if page.Items[0].IsTabsLayout() && !keepTabs {
			page.Items[0].Items = append(page.Items[0].Items, tabs.Items...)
		} else {
			page.Items = append(page.Items, tabs)
		}
	}

	return true
}

func addPageToPage(page1 *Element, page2 *Element) bool {
	if size := len(page1.Items); size > 0 {
		if size == 1 && page1.Items[0].IsTabsLayout() {
			return addTabsToPage(page1.Items[0], page2, false)
		} else {
			forms1 := page1.GetItemsByType("Form")
			forms2 := page2.GetItemsByType("Form")
			if forms1 != nil && forms2 != nil && len(forms2) == 1 && len(forms1) == 1 {
				params1 := forms1[0].GetParams(nil, nil)
				params2 := forms2[0].GetParams(nil, nil)
				params2 = append(params2, params1...)
				_form := NewSetParameterValuesForm(forms2[0].Name, params2...)
				page2.ReplaceItem(forms2[0], _form)
			} else {
				page2.Items = append(page2.Items, page1.Items...)
			}
		}
		return true
	}
	return false
}
func addPageToModule(page *Element, module *Element) bool {
	if page == nil || module == nil {
		return false
	}
	if size := len(page.Items); size > 0 {
		module.Items = append(module.Items, page)
		return true
	}
	return false
}
func addModuleToApp(module *Element, app *Element) bool {
	if size := len(module.Items); size > 0 {
		app.Items = append(app.Items, module)
		return true
	}
	return false
}
