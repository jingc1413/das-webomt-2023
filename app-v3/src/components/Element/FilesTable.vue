<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
    <el-col :span="12" v-if="viewMode !== provideKeys.viewModePrintValue"
      style="display: flex;justify-content: flex-end;">
      <div class="toolbar">
        <el-upload v-if="supportUpload" ref="upload" :action="uploadUrl" :data="uploadFormData" :headers="uploadHeaders"
          :show-file-list="false" :accept="uploadAccept" :before-upload="handleBeforeUpload"
          :on-success="handleUploadSuccess" :on-error="handleUploadError">
          <template #trigger>
            <el-button :loading="uploadLoading" type="primary" plain>Upload</el-button>
          </template>
        </el-upload>
        <el-dropdown trigger="click" v-if="supportConfiguration">
          <el-button type="primary" plain style="margin-left: 12px">Export
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <!-- <el-dropdown-item @click="handleLoadDeviceFile('moke.json')">Load</el-dropdown-item> -->
              <el-dropdown-item @click="handleExportConfigurationView()">
                Export Configuration
              </el-dropdown-item>
              <el-dropdown-item v-if="dasDevices.isCurrentDevice(0) && dev.cfg.supportCarrierConfig()"
                @click="handleExportCarrierConfiguration()">
                Export Carrier Configuration
              </el-dropdown-item>
              <el-dropdown-item v-if="dasDevices.isCurrentDevice(0) && dev.cfg.supportCustomCarrierConfig()"
                @click="handleExportCustomCarrierConfiguration()">
                Export Custom Carrier Configuration
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button v-if="supportUpgrade" @click="handleGetDeviceCurrentPacketInfo()" style="margin-left: 12px">
          Current Version
        </el-button>
        <el-button v-if="supportRestoreDefaultSettings && isAdmin" @click="handleRestoreDefaultSettings()"
          style="margin-left: 12px">
          Restore Default Settings
        </el-button>
        <el-button v-if="supportDeleteAll" @click="handleDeleteAll()" style="margin-left: 12px">
          Delete All
        </el-button>
        <el-button @click="getDeviceFileList()" style="margin-left: 12px">
          Refresh
        </el-button>
      </div>
    </el-col>
  </el-row>
  <el-table v-loading="loading" :data="tableData" :border="false" style="width: 100%" stripe table-layout="auto"
    :class="{ 'my-table-height': viewMode != provideKeys.viewModePrintValue }"
    :default-sort="{ prop: 'ModTime', order: 'descending' }">
    <el-table-column prop="FileName" label="File Name" sortable>
      <template #default="scope">
        <el-link v-if="supportViewFile" type="primary" target="_blank"
          @click="handleViewDeviceFile(scope.row['FileName'])">
          {{ scope.row["FileName"] }}
        </el-link>
        <el-link v-else-if="supportUpgrade" type="primary" target="_blank"
          @click="handleGetDeviceUpgradeFilePacketInfo(scope.row['FileName'])">
          {{ scope.row["FileName"] }}
        </el-link>
        <el-text v-else>{{ scope.row["FileName"] }}</el-text>
      </template>
    </el-table-column>
    <el-table-column v-if="fileType === 'UpgradeFile'" prop="ProductType" label="Product Type"
      :filters="productTypeFilters" :filter-method="handlerFilter">
      <template #default="scope">
        <span v-if="scope.row['ProductType']"> {{ scope.row["ProductType"] }}</span>
      </template>
    </el-table-column>
    <el-table-column v-if="fileType === 'UpgradeFile'" prop="ProductModel" label="Product Model"
      :filters="productModelFilters" :filter-method="handlerFilter">
      <template #default="scope">
        <span v-if="scope.row['ProductModel']"> {{ scope.row["ProductModel"] }}</span>
      </template>
    </el-table-column>
    <el-table-column v-if="fileType === 'UpgradeFile'" prop="Version" label="Version">
      <template #default="scope">
        <span v-if="scope.row['Version']"> {{ scope.row["Version"] }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="FileSize" label="File Size">
      <template #default="scope">
        <span v-if="scope.row['FileSize']"> {{ bytesToSize(scope.row["FileSize"]) }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="ModTime" label="Date" sortable>
      <template #default="scope">
        <span v-if="scope.row['ModTime']">{{
          deviceTimestampToDayJs(scope.row["ModTime"]).format("YYYY-MM-DD HH:mm")
        }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="Action" label="Operations" v-if="viewMode !== provideKeys.viewModePrintValue">
      <template #default="scope">
        <el-button link v-if="scope.row['FileName']" type="primary" size="small"
          @click="dev.files.downloadFile(fileType, scope.row['FileName'])">Download</el-button>
        <el-button link v-if="supportUpgrade" type="primary" size="small"
          @click="handleUpgradeDeviceFile(scope.row['FileName'], false)">Upgrade</el-button>
        <el-button link v-if="supportUpgrade" type="primary" size="small"
          @click="handleUpgradeDeviceFile(scope.row['FileName'], true)">Force Upgrade</el-button>
        <el-button link v-if="supportConfiguration" type="primary" size="small"
          @click="handleLoadDeviceFile(scope.row['FileName'])">Load</el-button>
        <el-button link type="primary" size="small"
          @click="handleDeleteDeviceFile(scope.row['FileName'])">Delete</el-button>
      </template>
    </el-table-column>
  </el-table>

  <export-config-page-view v-if="exportConfigurationState.visible"
    v-model:configDialogVisible="exportConfigurationState.visible" />

  <import-config-page-view v-if="loadConfigurationState.visible"
    v-model:configDialogVisible="loadConfigurationState.visible" :fileName="loadConfigurationState.filename" />

  <batch-operation-view v-if="batchActionsState.visible" v-model:batchDialogVisible="batchActionsState.visible"
    :batch-action-call-back="batchActionsState.actionCallBack"
    :batch-action-arguments="batchActionsState.actionArguments" :batch-action-names="batchActionsState.actionNames"
    :title="batchActionsState.title" />

  <el-dialog v-model="restoreDefaultSettingsState.visible" title="Restore Default Settings" width="40%"
    :show-close="false" :close-on-click-modal="false" :close-on-press-escape="false">
    <el-row justify="center">
      <el-col :span="16">
        <p>
          <el-icon style="margin-right: 32px;">
            <SuccessFilled v-if="restoreDefaultSettingsState.resetCpriLossCounter === true" color="#67C23A" />
            <CircleCloseFilled v-else-if="restoreDefaultSettingsState.resetCpriLossCounter === false" color="#f56c6c" />
            <InfoFilled v-else-if="restoreDefaultSettingsState.resetCpriLossCounter === undefined" />
          </el-icon>
          {{ 'Reset CPRI sync loss counter' }}
        </p>
        <p>
          <el-icon style="margin-right: 32px;">
            <SuccessFilled v-if="restoreDefaultSettingsState.resetDeviceAlarmState === true" color="#67C23A" />
            <CircleCloseFilled v-else-if="restoreDefaultSettingsState.resetDeviceAlarmState === false"
              color="#f56c6c" />
            <InfoFilled v-else-if="restoreDefaultSettingsState.resetDeviceAlarmState === undefined" />
          </el-icon>
          {{ 'Reset device alarm State' }}
        </p>
        <p v-if="dasDevices.currentProductType === 'AU'">
          <el-icon style="margin-right: 32px;">
            <SuccessFilled v-if="restoreDefaultSettingsState.deleteTopoRoot === true" color="#67C23A" />
            <CircleCloseFilled v-else-if="restoreDefaultSettingsState.deleteTopoRoot === false" color="#f56c6c" />
            <InfoFilled v-else-if="restoreDefaultSettingsState.deleteTopoRoot === undefined" />
          </el-icon>
          {{ 'Delete device topo root' }}
        </p>
        <p>
          <el-icon style="margin-right: 32px;">
            <SuccessFilled v-if="restoreDefaultSettingsState.deleteAccountsAndLogs === true" color="#67C23A" />
            <CircleCloseFilled v-else-if="restoreDefaultSettingsState.deleteAccountsAndLogs === false"
              color="#f56c6c" />
            <InfoFilled v-else-if="restoreDefaultSettingsState.deleteAccountsAndLogs === undefined" />
          </el-icon>
          {{ 'Delete Accounts And Logs' }}
        </p>
      </el-col>
    </el-row>

    <template #footer>
      <el-row justify="center">
        <el-button type="primary" plain v-if="restoreDefaultSettingsState.result !== undefined"
          @click="restoreDefaultSettingsState.visible = false">{{ $t("button.close") }}</el-button>
      </el-row>
    </template>
  </el-dialog>
</template>

<script>
import { useAuthStore } from "@/stores/auth"
import { useAppStore } from "@/stores/app";
import { useDasDevices } from "@/stores/das-devices";
import { useDasTopo } from "@/stores/topo";
import provideKeys from "@/utils/provideKeys.js";
import { ElMessage } from "element-plus";
import ExportConfigPageView from "@/components/ExportDeviceConfigPageView/ExportConfigPageView.vue";
import ImportConfigPageView from "@/components/ImportDeviceConfigPageView/ImportConfigPageView.vue";
import BatchOperationView from '@/components/BatchOperationProgressBar/BatchOperationView.vue';
import { bytesToSize, sleep } from "@/utils"
import apix from "@/api";
import model from "@/stores/model"

export default {
  name: "MyFilesTable",
  components: { ExportConfigPageView, ImportConfigPageView, BatchOperationView },
  inject: ["viewMode"],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const auth = useAuthStore();
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dasTopo = useDasTopo();
    const dev = dasDevices.currentDevice;

    return {
      auth,
      appStore,
      dasDevices,
      dasTopo,
      dev,
      deviceTimestampToDayJs: model.deviceTimestampToDayJs,
      provideKeys,
      bytesToSize,
    };
  },
  data() {
    const fileType = this.data.Name;
    const isAdmin = this.auth.loginUserName === "admin";
    let supportUpload = false;
    let uploadAccept = '';
    let uploadCheckExtension = undefined;
    let supportUpgrade = false;
    let supportConfiguration = false;
    let supportDeleteAll = false;
    let supportRestoreDefaultSettings = false;
    let supportViewFile = false;

    if (fileType === "UpgradeFile") {
      supportUpload = true;
      uploadAccept = '.zip'
      supportUpgrade = true;
      uploadCheckExtension = function (file) {
        const allowedExtensions = ['zip'];
        const fileExtension = file.name.split('.').pop();
        if (!allowedExtensions.includes(fileExtension)) {
          this.$message.error(`File type not allowed: ${fileExtension}`);
          return false;
        }
        if (!file.name.match(/^(iDAS|DDAS)_.*$/)) {
          this.$message.error(`Invalid upgrade file`);
          return false;
        }
        return true;
      }


    }
    if (fileType === "ConfigFile") {
      supportUpload = true;
      uploadAccept = '.csv,.json'
      uploadCheckExtension = function (file) {
        const allowedExtensions = ['csv', 'json'];
        const fileExtension = file.name.split('.').pop();
        if (!allowedExtensions.includes(fileExtension)) {
          this.$message.error(`File type not allowed: ${fileExtension}`);
          return false;
        }
        return true;
      }
      supportConfiguration = true;
      supportRestoreDefaultSettings = true;
    }
    if (fileType.endsWith("Log")) {
      supportDeleteAll = this.auth.superTestDisabled === false;
      supportViewFile = true;
    }

    return {
      loading: false,
      uploadLoading: false,

      fileType,
      isAdmin,
      supportUpload,
      supportUpgrade,
      supportConfiguration,
      supportDeleteAll,
      supportRestoreDefaultSettings,
      supportViewFile,

      uploadAccept,
      uploadCheckExtension,
      uploadFormData: {}, // Form data to be sent to the API
      uploadHeaders: {}, // Headers to be sent with the request
      uploadUrl: apix.DasApiBase + "/devices/" + this.dasDevices.currentDeviceSub + "/files/" + fileType,

      exportConfigurationState: {
        visible: false,
      },
      loadConfigurationState: {
        visible: false,
        filename: undefined,
      },
      restoreDefaultSettingsState: {
        visible: false,
      },
      batchActionsState: {
        visible: false,
        actionCallBack: null,
        actionArguments: [],
        actionNames: [],
        title: 'Batch Action',
      },
    };
  },
  computed: {
    tableData() {
      return this.dev.files.fileList(this.fileType);
    },
    productTypeFilters() {
      if (this.fileType === "UpgradeFile" && this.tableData) {
        const values = [];
        this.tableData.forEach(row => {
          if (row.ProductType && !values.includes(row.ProductType)) {
            values.push(row.ProductType);
          }
        });
        values.sort((a, b) => a.localeCompare(b))
        const filters = values.map(v => {
          return { text: v, value: v };
        })
        return filters;
      }
      return undefined;
    },
    productModelFilters() {
      if (this.fileType === "UpgradeFile" && this.tableData) {
        const values = [];
        this.tableData.forEach(row => {
          if (row.ProductModel && !values.includes(row.ProductModel)) {
            values.push(row.ProductModel);
          }
        });
        values.sort((a, b) => a.localeCompare(b))
        const filters = values.map(v => {
          return { text: v, value: v };
        })
        return filters;
      }
      return undefined;
    },
  },
  mounted() {
    this.getDeviceFileList();
  },
  methods: {
    handleGetDeviceCurrentPacketInfo: async function () {
      this.appStore.openViewFileDialog({
        title: "Current Version",
        handleLoadingData:()=>{
          return this.dev.files.getCurrentPacketInfo();
        },
      })
    },
    handleGetDeviceUpgradeFilePacketInfo: async function (filename = undefined) {
      this.appStore.openViewFileDialog({
        title: filename,
        handleLoadingData:()=>{
          return this.dev.files.getUpgradeFilePacketInfo(filename);
        },
      })
    },
    getDeviceFileList: async function () {
      this.loading = true;
      await this.dev.files.getFileList(this.fileType);
      this.loading = false;
    },
    handlerFilter(value, row, column) {
      const property = column['property'];
      return row[property] === value;
    },
    handleUpgradeDeviceFile: async function (filename, force) {
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: "Confirm to upgrade " + filename,
        callback: ok => {
          if (ok) {
            self.deviceUpgrade.startDeviceUpgrade(filename, force)
          }
        },
      })
    },
    handleLoadDeviceFile: async function (filename) {
      const self = this;
      // iDAS_A402_Flatness_Coefficient_v1.0_4c25_20230224.csv
      if (filename.match(/^.*Flatness_Coefficient.*\.csv$/)) {
        this.appStore.openConfirmDialog({
          title: 'Confirm',
          content: "Load Flatness Coefficient Configuration.",
          callback: ok => {
            if (ok) {
              self.dev.cfg.loadFlatnessCoefficientConfigFile(filename);
            }
          }
        });
      } else if (filename.match(/^.*PAInitConfig.*\.csv$/)) {
        this.appStore.openConfirmDialog({
          title: 'Confirm',
          content: "Load PA Initial Configuration.",
          callback: ok => {
            if (ok) {
              self.dev.cfg.loadPAInitConfigFile(filename);
            }
          }
        });
      } else if (filename.match(/^.*CustomCarrier.*\.json$/)) {
        this.appStore.openConfirmDialog({
          title: 'Confirm',
          content: "Before loading custom carrier configuration file, make sure that the Carrier Config Control is disabled. Otherwise loading file will fail.",
          callback: ok => {
            if (ok) {
              self.dev.cfg.loadCustomCarrierConfigFile(filename);
            }
          }
        });
      } else if (filename.match(/^.*Carrier.*\.csv$/)) {
        this.appStore.openConfirmDialog({
          title: 'Confirm',
          content: "Before loading carrier configuration file, make sure that the Carrier Config Control is disabled. Otherwise loading file will fail.",
          callback: ok => {
            if (ok) {
              self.dev.cfg.loadCarrierConfigFile(filename);
            }
          }
        });
      } else if (filename.match(/^.*\.json$/)) {
        this.loadConfigurationState.filename = filename;
        this.loadConfigurationState.visible = true;
      }
      return
    },
    handleViewDeviceFile: async function (filename) {
      this.appStore.openViewFileDialog({
        title: filename,
        handleLoadingData:()=>{
          return this.dev.files.getFile(this.fileType, filename)
        },
        // supportSave: true,
        // handleSave: ()=>{
        //   this.dev.files.downloadFile(this.fileType, filename)
        // }
      })
      return
    },
    handleExportConfigurationView() {
      this.exportConfigurationState.visible = true;
    },
    handleExportCarrierConfiguration() {
      this.dev.cfg.getCarrierConfigFile();
    },
    handleExportCustomCarrierConfiguration() {
      this.dev.cfg.getCustomCarrierConfigFile();
    },
    handleDeleteDeviceFile: async function (filename) {
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: 'Confirm to delete the device file',
        callback: ok => {
          if (ok) {
            self.doDeleteDeviceFile(filename);
          }
        }
      });
    },
    doDeleteDeviceFile: async function (filename) {
      if (filename) {
        this.loading = true;
        await this.dev.files.deleteFile(this.fileType, filename, true);
        this.loading = false;
      }
    },
    handleBeforeUpload(file) {
      const ok = this.uploadCheckExtension ? this.uploadCheckExtension(file) : true;
      if (!ok) {
        return false;
      }
      this.uploadLoading = true;
      return true;
    },
    handleUploadSuccess(response) {
      ElMessage.success("File uploaded successfully");
      this.uploadLoading = false;
      this.getDeviceFileList();
    },
    handleUploadError(error) {
      try {
        const msg = JSON.parse(error.message).message;
        ElMessage.error("File uploaded failed: " + msg);
      } catch (e) {
        ElMessage.error("File uploaded failed");
      }
      this.uploadLoading = false;
      this.getDeviceFileList();
    },

    handleRestoreDefaultSettings() {
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: 'Confirm to restore default settings',
        callback: ok => {
          if (ok) {
            setTimeout(() => {
              self.doRestoreDefaultSettings();
            }, 500);
          }
        }
      });
    },
    doRestoreDefaultSettings: async function () {
      const state = this.restoreDefaultSettingsState;

      try {
        state.resetCpriLossCounter = undefined;
        state.resetDeviceAlarmState = undefined;
        state.deleteTopoRoot = undefined;
        state.deleteAccountsAndLogs = undefined;
        state.result = undefined;
        state.visible = true;

        await sleep(500);

        await this.dev.funcs.enterFactoryMode(true)
        const isFacotryMode = await this.dev.funcs.isFacotryMode();
        if (!isFacotryMode) {
          state.result = false;
          return;
        }
        state.resetCpriLossCounter = await this.dev.funcs.resetCpriSyncLossCounter(true);
        if (!state.resetCpriLossCounter) {
          state.result = false;
          return;
        }
        state.resetDeviceAlarmState = await this.dev.funcs.resetDeviceAlarmState(true);
        if (!state.resetDeviceAlarmState) {
          state.result = false;
          return;
        }
        if (this.dasDevices.currentProductType === 'AU') {
          state.deleteTopoRoot = await this.dasTopo.deleteTopoRootNode();
          if (!state.deleteTopoRoot) {
            state.result = false;
          }
        }

        await this.dev.funcs.quitFactoryMode();
        state.deleteAccountsAndLogs = await this.dev.funcs.deleteDeviceKeyAndLogs();
        if (!state.deleteAccountsAndLogs) {
          state.result = false;
          return;
        }
        state.result = true;
        ElMessage.success("Restore default settings successfully")
      } catch (e) {
        console.error(e);
        ElMessage.error("Restore default settings failed")
      } finally {
        state.visible = true;
        this.getDeviceFileList();
      }
    },
    handleDeleteAll() {
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: 'Confirm to delete all logs',
        callback: ok => {
          if (ok) {
            setTimeout(() => {
              self.batchActionsState.actionNames = self.tableData.map(item => {
                return {
                  ...item,
                  name: item.FileName
                }
              })
              self.batchActionsState.actionCallBack = self.dev.files.deleteFile;
              self.batchActionsState.actionArguments = self.batchActionsState.actionNames.map(item => {
                return [
                  self.fileType,
                  item.FileName
                ]
              })
              self.batchActionsState.visible = true;
            }, 500);
          }
        }
      });

    }
  },
};
</script>
<style lang="scss" scoped></style>
