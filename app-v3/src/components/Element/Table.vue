<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <my-files-table v-if="data.Type === 'Table:Files'" :owner="owner" :data="data" />
  <my-users-table v-else-if="data.Type === 'Table:Users'" :owner="owner" :data="data" />
  <my-alarm-logs-table v-else-if="data.Type === 'Table:AlarmLogs'" :owner="owner" :data="data" />
  <my-inventory-table v-else-if="data.Type === 'Table:Inventory'" :owner="owner" :data="data" />
  <my-firmwares-table v-else-if="data.Type === 'Table:Firmwares'" :owner="owner" :data="data" />
  <template v-else>
    <el-row style="height: 48px;">
      <el-col :span="12">
        <h4 v-if="data.Name">
          <span>{{ data.Name }}</span>
        </h4>
      </el-col>
      <el-col :span="12" style="display: flex;justify-content: flex-end;">
        <div v-if="toolbarItems[0]?.length > 0" class="toolbar">
          <div v-for="item in toolbarItems[0]" :key="item.Key" style="margin-right: 12px;">
            <my-toolbar-element :owner="owner" :data="item"
              v-if="(isPrintMode && item.Access != 'wo') || !isPrintMode" />
          </div>
        </div>
        <div v-if="toolbarItems[1]?.length > 0 && !isPrintMode" class="toolbar">
          <div v-for="item in toolbarItems[1]" :key="item.Key" style="margin-right: 12px;">
            <el-button v-if="item.Key === 'add'" @click="addTableItem" type="primary" plain>{{ item.Name }}</el-button>
          </div>
        </div>
        <div v-if="toolbarVisible" class="toolbar">
          <el-button v-if="supportTableImport" type="primary" plain @click="doExport()">Export</el-button>
          <el-popover v-if="supportTableImport" :visible="importDialogVisible" placement="bottom-end" :width="800"
            trigger="click">
            <template #reference>
              <el-button v-if="supportTableImport" type="primary" plain @click="openImportDialog()">Import</el-button>
            </template>
            <el-row>
              <el-col :span="12">
                <h4 style="margin-left:8px;">Import Parameters</h4>
              </el-col>
              <el-col :span="12" style="text-align: end;">
                <div class="toolbar">
                  <el-button link @click="closeImportDialog()">
                    <el-icon :size="16">
                      <Close />
                    </el-icon>
                  </el-button>
                </div>
              </el-col>
            </el-row>
            <el-row style="margin-bottom: 16px;">
              <el-upload ref="importFilesRef" v-model:file-list="importFiles" accept=".xlsx,.xls" :limit="1"
                :http-request="handleCustomUpload" :on-remove="removeImportFile" :on-exceed="handleFileExceed"
                style="display: inline-flex; margin-left: 8px;">
                <el-button>Select</el-button>
              </el-upload>
              <template v-if="importTableVisible">
                <el-table :data="importTableData" :border="false" @selection-change="handleImportSectionChange"
                  table-layout="auto" style="padding-bottom: 16px;" :cell-class-name="importTableCellClassName"
                  height="calc(40vh - 80px)">
                  <el-table-column type="selection" width="55" />
                  <template v-for="column in importTableColumns">
                    <el-table-column v-if="column.Style?.fixed || column.Access !== 'ro'" :key="column.Key"
                      :prop="column.Key" :label="column.Name" :fixed="column.Style?.fixed"
                      :width="column?.Style?.width">
                      <template #default="scope">
                        <my-element v-if="scope.row[column.Key]" :owner="owner" :data="scope.row[column.Key]" />
                      </template>
                    </el-table-column>
                  </template>
                </el-table>
                <el-col :span="24">
                  <el-button type="primary" @click="setImportParameterValues"
                    :disabled="importTableSelected === undefined || importTableSelected.length < 1">Submit</el-button>
                </el-col>
              </template>
            </el-row>
          </el-popover>
          <el-button v-if="supportRefresh" @click="getParameterValues(true)">Refresh</el-button>
        </div>
      </el-col>
    </el-row>
    <el-table ref="tableRef" v-loading="loading" :data="tableData" :border="false" table-layout="auto"
      style="width: 100%;" :style="{ minHeight: data?.Style?.height }"
      :show-header="tableData.length == 0 || !isPrintMode" :highlight-current-row="supportTableRowClick"
      @row-click="handleRowChange" :row-class-name="tableRowClassName">
      <el-table-column v-for="column in tableColumns" :key="column.Key" :prop="column.Key" :label="column.Name"
        :fixed="column?.Style?.fixed" :width="column?.Style?.width">
        <template #default="scope">
          <span v-if="scope.$index == 0 && isPrintMode">{{ column.Name }}</span>
          <my-element v-else-if="scope.row[column.Key]" :owner="owner" :data="scope.row[column.Key]" />
        </template>
      </el-table-column>
    </el-table>
  </template>
</template>

<script>
import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import { read, utils, writeFile } from "xlsx";
import provideKeys from '@/utils/provideKeys.js';
import { useAppStore } from "@/stores/app";

export default {
  name: "MyTable",
  inject: ['viewMode'],
  props: {
    owner: {
      type: Object,
      default: undefined,
    },
    data: {
      type: Object,
      default: undefined,
    },
  },
  setup() {
    const auth = useAuthStore();
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    return {
      auth,
      appStore,
      dasDevices,
      dev,
      provideKeys,
    };
  },
  data() {
    const writableParams = [];
    const tableInvalidKey = this.data?.Style?.invalidKey;
    const tableInvalidValue = this.data?.Style?.invalidValue;
    const tableViewKey = this.data?.Style?.viewKey;
    let supportTableRowClick = this.data?.supportTableRowClick;

    if (this.viewMode == provideKeys.viewModePrintValue) {
      supportTableRowClick = false;
    }

    if (this.data?.Data) {
      for (const row of this.data.Data) {
        for (const key in row) {
          const param = row[key];
          if (param) {
            if (param.Access === "rw" || param.Access === "wo") {
              writableParams.push(param);
            }
          }
        }
      }
    }
    let supportTableImport = this.auth.superTestDisabled ? false : writableParams?.length > 0;
    if (["RF Module Mapping", "SNMP User List"].includes(this.data?.Name)) {
      supportTableImport = false;
    }
    if (this.data?.accessKeys?.set === false) {
      supportTableImport = false;
    }
    const supportRefresh = this.viewMode != provideKeys.viewModePrintValue;
    const toolbarVisible = this.viewMode != provideKeys.viewModePrintValue;
    const scrollable = this.viewMode != provideKeys.viewModePrintValue && this.data.Key !== 'rf_module_mapping';
    return {
      loading: true,
      toolbarVisible,
      scrollable,
      supportTableImport,
      supportRefresh,
      supportTableRowClick,
      tableInvalidKey,
      tableInvalidValue,
      tableViewKey,
      writableParams,
      importDialogVisible: false,
      importFiles: [],
      importTableVisible: false,
      importTableData: [],
      importTableSelected: [],
    };
  },
  computed: {
    isPrintMode() {
      return this.viewMode == provideKeys.viewModePrintValue
    },
    toolbarItems() {
      const toolbarItems = [];
      const tableToolbarItems = [];
      if (this.data?.Actions?.toolbar?.Items) {
        this.data?.Actions?.toolbar?.Items.forEach(v => {
          if (v?.accessKeys?.set != true) {
            return
          }
          if (v.Type === "Button" && v.Actions?.click?.Name === "AddItem") {
            tableToolbarItems.push(v);
          } else {
            toolbarItems.push(v);
          }
        })
      }
      return [toolbarItems, tableToolbarItems];
    },
    tableColumns() {
      const tableColumns = this.data?.Items?.filter(v => {
        if (v.Style?.hidden) {
          return false;
        }
        if (this.viewMode == provideKeys.viewModePrintValue && v.Key == "action") {
          return false;
        }
        return true;
      }) || [];
      return tableColumns;
    },
    filterTableData() {
      try {
        if (!this.tableInvalidKey) {
          return this.data.Data;
        }

        const filterTableData = this.data.Data.filter(row => {
          const item = row[this.tableInvalidKey];
          if (item) {
            const param = this.dev.params.getParam(item.OID)
            if (param && param.Value) {
              if (param.Value !== undefined && this.tableInvalidValue != undefined && param.Value !== this.tableInvalidValue) {
                return true;
              }
            }
          }
          return false;
        })
        return filterTableData;
      } catch (e) {
        return [];
      }

    },
    tableData() {
      if (this.isPrintMode) {
        return ['', ...this.filterTableData];
      }
      return [...this.filterTableData];
    },
  },
  mounted() {
    this.getParameterValues();
  },
  methods: {
    getParameterValues: async function (showMessage = false) {
      try {
        if (!this.supportRefresh) return;
        if (this.data?.oids?.length > 0) {
          this.loading = true;
          await this.dev.params.getParameterValues({
            oids: this.data.oids,
            showMessage,
            values: this.data.defaultValues,
          });
          this.loading = false;
        }
      } finally {
        this.loading = false;
      }
    },
    addTableItem() {
      const row = this.data.Data.find(v => {
        const item = v[this.tableInvalidKey];
        if (item) {
          const param = this.dev.params.getParam(item.OID)
          if (param) {
            if (!param.Value) {
              return true;
            }
            if (this.tableInvalidValue != undefined && param.Value === this.tableInvalidValue) {
              return true;
            }
          }
        }
        return false;
      });
      if (!row) {
        return;
      }
      if (!this.tableViewKey) {
        return;
      }
      const action = row[`_actions`]?.Actions?.click;
      if (action && action.Name === "ViewPage") {
        const page = action.Items[0];
        if (page) {
          this.dev.layout.openViewPage({ page: this.owner, viewPage: page });
        }
      }
    },
    doExport() {
      const rows = [];
      this.data.Data.forEach((row) => {
        const row2 = {};
        this.tableColumns.forEach((col) => {
          const obj = row[col.Key];
          if (col.Key === this.data?.Unique) {
            row2[col.Key] = obj.Value;
          } else if (obj?.Type === "Param") {
            const param = this.dev.params.getParam(obj.OID);
            row2[col.Key] = param.Value;
          } else if (obj?.Type?.startsWith("Component:ParamGroup")) {
            const values = this.dev.params.getValues(obj.Items);
            row2[col.Key] = values.join(",");
            // } else if (obj?.Type === "Button") {
            //   row2[col.Name] = obj.Name;
          }
        });
        rows.push(row2);
      });
      const worksheet = utils.json_to_sheet(rows);
      const workbook = utils.book_new();
      utils.book_append_sheet(workbook, worksheet, "Parameter Values");
      const filename = this.dasDevices.currentDeviceFileName(
        "Export_" + this.owner.Name.replaceAll(" ", "") + "_" + this.data.Name.replaceAll(" ", ""), false)
        + ".xlsx";
      writeFile(workbook, filename, { compression: true });
    },
    openImportDialog() {
      this.importDialogVisible = true;
    },
    closeImportDialog() {
      this.importDialogVisible = false;
      this.importFiles = [];
      this.importTableVisible = false;
      this.importTableData = [];
      this.importTableSelected = [];
    },
    handleCustomUpload(e) {
      this.openImportFile(e.file)
    },
    openImportFile(file) {
      const self = this;
      const table = this.data;
      const reader = new FileReader();
      reader.onload = (ev) => {
        const filedata = ev.target.result;
        const workbook = read(filedata, { type: "binary" });
        var worksheet = workbook.Sheets[workbook.SheetNames[0]];
        const raw = utils.sheet_to_json(worksheet, { header: 1 });
        const result = self.dev.importer.importTableDataFromRaw(table, raw);
        self.importTableColumns = self.tableColumns.filter(v => result.columneKeys.includes(v.Key));
        self.importTableData = result.data;
        self.importTableVisible = true;
      };
      reader.readAsBinaryString(file);
    },
    handleFileExceed(files) {
      this.$refs['importFilesRef'].clearFiles();
      const file = files[0];
      this.$refs['importFilesRef'].handleStart(file);
      this.openImportFile(file);
    },
    removeImportFile(e) {
      this.importTableData = [];
      this.importTableSelected = [];
      this.importTableVisible = false;
    },
    importTableCellClassName({ row, column, rowIndex, columnIndex }) {
      if (this.importTableData && column) {
        const rowData = this.importTableData[rowIndex];
        if (rowData) {
          const param = rowData[column.property];
          if (param) {
            if (param.validate != undefined) {
              return "cell-invalid";
            }
            if (param.isChanged) {
              return "cell-changed";
            }
          }
        }
      }
      return undefined;
    },
    handleImportSectionChange(v) {
      this.importTableSelected = v;
    },
    setImportParameterValues: async function () {
      const self = this;
      const values = [];
      this.importTableSelected.forEach((row) => {
        Object.keys(row).forEach((key) => {
          const obj = row[key];
          if (obj?.Type === "Param") {
            values.push({
              oid: obj.OID,
              value: obj.Value,
              validate: obj.validate,
            });
          } else if (obj?.Type?.startsWith("Component:ParamGroup")) {
            const objValues = obj.Value.split(",");
            for (const i in obj?.Items) {
              const obj2 = obj.Items[i];
              const obj2Value = objValues[i];
              if (obj2 && obj2Value !== undefined) {
                values.push({
                  oid: obj2.OID,
                  value: obj2Value,
                });
              }
            }
          }
        });
      });
      const invalidValues = values.filter(v => v.validate !== undefined);
      if (invalidValues.length > 0) {
        this.appStore.openConfirmDialog({
          title: 'WARNING',
          content: 'It is forbidden to set parameter values with invalid values.',
          supportCancel: false,
        });
        return
      }
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: 'Confirm to set parameter values',
        callback: ok => {
          if (ok) {
            self.dev.params.setParameterValues({ values, showMessage: true });
          }
        }
      });
    },
    handleRowChange(row) {
      if (this.supportTableRowClick && row._actions?.Actions?.click) {
        const self = this;
        let isDisable = this.getRowActionClickIsDisable(row);
        if (isDisable == -1) {
          this.dev.doAction(row._actions?.Actions?.click, {
            owner: this.owner,
            onCloseViewPage: () => {
              self.$refs['tableRef']?.setCurrentRow();
            }
          })
        } else {
          self.$refs['tableRef']?.setCurrentRow();
        }
        return
      }
    },
    tableRowClassName({ row }) {
      let classString = ''
      if (this.supportTableRowClick && row._actions?.Actions?.click) {
        let isDisable = this.getRowActionClickIsDisable(row);
        if (isDisable == -1) {
          classString += ' support_table_tow_click_cursor'
        }
      }
      return classString
    },
    getRowActionClickIsDisable(row) {
      return row._actions?.Actions?.click?.Items.findIndex(item => item.InputDisabled)
    }
  },
};
</script>
<style lang="scss" scoped>
.my-table-height {
  height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height) - 10px)
}

.support_table_tow_click_cursor {
  cursor: pointer;
}
</style>

<style lang="scss">
.support_table_tow_click_cursor {
  cursor: pointer;
}
</style>
