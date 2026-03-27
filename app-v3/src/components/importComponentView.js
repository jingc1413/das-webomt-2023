import MyButton from "./Element/Button.vue"
import MyLabel from "./Element/Label.vue"
import MyAlert from "./Element/Alert.vue"

import MyParam from "./Element/Param.vue"
import MyParamGroup from "./Element/ParamGroup.vue"
import MyTable from "./Element/Table.vue"
import MyFilesTable from "./Element/FilesTable.vue"
import MyUsersTable from "./Element/UsersTable.vue"
import MyFirmwaresTable from "./Element/FirmwaresTable.vue"
import MyAlarmLogsTable from "./Element/AlarmLogsTable.vue"
import MyInventoryTable from "./Element/InventoryTable.vue"
import MyAddressInterfaceForm from "./Element/AddressInterace.vue"

import MyParamInput from "./Element/ParamInput.vue"
import MyParamGroupInput from "./Element/ParamGroupInput.vue"
import MyStatisticGroup from "./Element/StatisticGroup.vue"

import MyToolbarElement from "./Element/ToolbarElement.vue"
import MyFormElement from "./Element/FormElement.vue"
import MyFormLayout from "./Element/FormLayout.vue"
import MyForm from "./Element/Form.vue"
import MyElement from "./Element/Element.vue"
import MyLayout from "./Element/Layout.vue"
import MyDialog from "./Element/Dialog.vue"
import MyConfirmDialog from "./Element/ConfirmDialog.vue"
import MyViewFileDialog from "./Element/ViewFileDialog.vue"

import MyMainPage from "./Element/MainPage.vue"
import MyViewPage from "./Element/ViewPage.vue"
import MyEditableInput from "./Element/EditableInput.vue"
import FormEditableInput from "@/components/Element/FormEditableInput.vue";
import MyTopo from "@/components/Topo/Topo.vue"
import MyGraphTopo from "@/components/Topo/GraphTopo.vue"
import MyTreeTopo from "@/components/Topo/TreeTopo.vue"
import MyTopoInfo from "@/components/Topo/TopoInfo.vue"
import MyInvalidView from "./Element/InvalidView.vue"
import MyStatisticsPage from "@/components/SystemStats/Stats.vue"
import MyPingDiagPage from "@/components/PingView/PingView.vue"

import TimePiece from "./time-piece.vue";

import SidebarItem from '@/layouts/App1/components/Sidebar/SidebarItem.vue'
import ChangePassword from '@/components/ChangePassword/ChangePassword.vue';

import ConfigWizard from '@/components/ConfigurationWizard/ConfigWizardMainPage.vue';
import DeviceTreeSelected from '@/components/Element/DeviceTreeSelected.vue';

import BandwidthSelectorView from './BandwidthSelector/BandwidthSelectorView.vue';
import TableSearchBar from '@/components/Element/TableSearchBar.vue';
import TableDataExportButton from '@/components/Element/TableDataExportButton.vue';

export default function importComponentView(appRoot) {

  appRoot.component("my-button", MyButton)
  appRoot.component("my-label", MyLabel)
  appRoot.component("my-alert", MyAlert)

  appRoot.component("my-param", MyParam)
  appRoot.component("my-param-group", MyParamGroup)
  appRoot.component("my-statistic-group", MyStatisticGroup)

  appRoot.component("my-table", MyTable)
  appRoot.component("my-files-table", MyFilesTable)
  appRoot.component("my-users-table", MyUsersTable)
  appRoot.component("my-firmwares-table", MyFirmwaresTable)
  appRoot.component("my-alarm-logs-table", MyAlarmLogsTable)
  appRoot.component("my-inventory-table", MyInventoryTable)

  appRoot.component("my-address-interface-form", MyAddressInterfaceForm)

  appRoot.component("my-param-input", MyParamInput)
  appRoot.component("my-param-group-input", MyParamGroupInput)
  appRoot.component("my-form-element", MyFormElement)
  appRoot.component("my-form-layout", MyFormLayout)
  appRoot.component("my-form", MyForm)
  appRoot.component("my-toolbar-element", MyToolbarElement)

  appRoot.component("my-element", MyElement)
  appRoot.component("my-layout", MyLayout)
  appRoot.component("my-dialog", MyDialog)
  appRoot.component("my-confirm-dialog", MyConfirmDialog);
  appRoot.component("my-view-file-dialog", MyViewFileDialog);

  appRoot.component("my-main-page", MyMainPage)
  appRoot.component("my-view-page", MyViewPage)
  appRoot.component("my-editable-input", MyEditableInput)
  appRoot.component("my-form-editable-input", FormEditableInput)
  appRoot.component("my-topo", MyTopo)
  appRoot.component("my-graph-topo", MyGraphTopo)
  appRoot.component("my-tree-topo", MyTreeTopo)
  appRoot.component("my-topo-info", MyTopoInfo)
  appRoot.component("my-statistics-page", MyStatisticsPage)
  appRoot.component("my-ping-diag-page", MyPingDiagPage)

  appRoot.component("my-invalid-view", MyInvalidView)

  appRoot.component("my-time-piece", TimePiece)

  appRoot.component('SidebarItem', SidebarItem)

  appRoot.component('my-change-password', ChangePassword)

  appRoot.component('my-config-wizard', ConfigWizard)
  appRoot.component('my-device-tree-selected', DeviceTreeSelected)

  appRoot.component('my-bandwidth-selector-view', BandwidthSelectorView)
  appRoot.component('my-table-search-bar', TableSearchBar)
  appRoot.component('my-table-data-export-button', TableDataExportButton)

}