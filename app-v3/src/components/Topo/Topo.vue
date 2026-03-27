<template>
  <div class="" style="">
    <el-row style="padding-inline:16px; border-bottom:solid 1px var(--el-menu-border-color);">
      <el-col :span="viewMode == provideKeys.viewModeDefaultValue ? 16 : 24">

        <el-row v-if="viewMode == provideKeys.viewModeDefaultValue">
          <div class="loading_device_progress" :style="{ width: loadingProgress }" v-show="loadingState?.loading">
          </div>
          <el-col :span="8">
            <h3 v-if="page.Name">
              {{ page.Name }}
              <span v-show="loadingState?.loading" class="sub-title">
                {{ "Query devices... " }}
              </span>
              <span v-show="loadingState?.loading && loadingState.total > 0" class="sub-title">
                {{ loadingState.index }} / {{ loadingState.total }}
              </span>
            </h3>
          </el-col>
          <el-col :span="16" style="display: flex;justify-content: flex-end;">
            <div class="toolbar">
              <el-button size="small" v-if="dasDevices.isCurrentDevice(0)" circle @click="handleDeleteAll()">
                <el-icon>
                  <Delete />
                </el-icon>
              </el-button>
              <el-button :disabled="loadingState?.loading" size="small" style="margin-right: 12px;" circle
                @click="handleRefresh(true)">
                <el-icon>
                  <Refresh />
                </el-icon>
              </el-button>
              <el-select v-model="secondContent" placeholder="Select" size="small" style="margin-right: 12px;width: 180px">
                <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-radio-group v-model="topoTypeModel" size="small" style="margin-right: 12px;">
                <el-radio-button label="Graph" value="Graph" />
                <el-radio-button label="Tree" value="Tree" />
              </el-radio-group>
            </div>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24" v-show="isGraphTopo" style="width: 100%;">
            <span class="account-for-topo">
              <div>
                <img src="@/assets/green.png" width="20" height="20" style="margin-top:20px;">
                <span style="margin-left: 16px;">Normal</span>
              </div>
              <div>
                <img src="@/assets/gray.png" width="20" height="20" style="margin-top:20px;">
                <span style="margin-left: 16px;">Offline</span>
              </div>
              <div>
                <img src="@/assets/red.png" width="20" height="20" style="margin-top:20px;">
                <span style="margin-left: 16px;">Alarm</span>
              </div>
            </span>
            <my-graph-topo :secondContent="secondContent" @select-device="handleSelectDevice"
              @delete-device="handleDeleteDevice" :isGraphTopo="isGraphTopo"></my-graph-topo>
          </el-col>
          <el-col :span="24" v-show="!isGraphTopo">
            <el-scrollbar
              :style="{ height: viewMode == provideKeys.viewModeDefaultValue ? 'calc(100vh - var(--header-height) - 90px)' : 'auto' }">
              <my-tree-topo :secondContent="secondContent" @select-device="handleSelectDevice"
                @delete-device="handleDeleteDevice"></my-tree-topo>
            </el-scrollbar>
          </el-col>
        </el-row>
      </el-col>
      <el-col :span="8" style="padding-inline:16px; border-left:solid 1px var(--el-menu-border-color);"
        v-if="viewMode == provideKeys.viewModeDefaultValue && selectSubID != undefined">
        <el-row>
          <h4>Device Info</h4>
        </el-row>
        <my-topo-info :subID="selectSubID" />
        <el-row v-if="true">
          <el-col :span="12">
            <el-button type="primary" size="small" :disable="selectDeviceInfo.ConnectState >= 6"
              @click="handleOpenDevicePage(selectSubID)">Jump</el-button>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import { useDasDevices } from '@/stores/das-devices'
import { useDasTopo } from "@/stores/topo"
import provideKeys from '@/utils/provideKeys.js'
import Cookies from 'js-cookie'
import {getPageVersion} from '@/utils/index.js';
import { useAppStore } from '@/stores/app'
export default {
  name: 'topo',
  inject: ['viewMode'],
  props: {
    page: Object,
  },
  setup() {
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dasTopo = useDasTopo();
    const dev = dasDevices.currentDevice;
    return {
      appStore,
      dev,
      dasDevices,
      dasTopo,
      provideKeys
    }
  },
  data() {
    return {
      loading: false,
      active: false,
      options: [{ value: 'Device Name', label: 'Device Name' }, { value: 'Location', label: 'Location' }],
      selectSubID: undefined,
    }
  },
  computed: {
    topoTypeModel: {
      get() {
        return this.dasTopo.topoTypeModel
      },
      set(value) {
        this.dasTopo.changeTopoTypeModel(value)
      }
    },
    secondContent: {
      get() {
        return this.dasTopo.treeDeviceName
      },
      set(value) {
        console.log({value});
        this.dasTopo.changeTreeDeviceName(value)
      }
    },
    isGraphTopo() {
      return this.topoTypeModel !== "Tree"
    },
    loadingState() {
      return this.dasDevices.deviceInfos?.loadingState;
    },
    loadingProgress() {
      return this.loadingState.progress + '%'
    }
  },
  mounted() {
    // this.isMasterAu()
    this.initTopo();
    //  this.isGraphTopo =Cookies.get('isGraphTopo')
  },
  created() { },
  methods: {
    initTopo() {
    },
    handleSelectDevice: function (subID) {
      if (this.viewMode != provideKeys.viewModeDefaultValue) {
        return;
      }
      this.selectSubID = String(subID)
      this.selectDeviceInfo = this.dasDevices.getDeviceInfo(subID);
      this.dasDevices.updateDeviceInfo(subID)
    },
    handleDeleteDevice: async function (subID) {
      if (this.viewMode != provideKeys.viewModeDefaultValue) {
        return;
      }
      const isFacotryMode = await this.dev.funcs.isFacotryMode();
      if (!isFacotryMode) {
        ElMessage.warning("Not working in factory mode");
        return;
      }
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: "Confirm to delete device node",
        callback: ok => {
          if (ok) {
            self.dasTopo.deleteTopoNode(subID);
          }
        }
      });
    },
    handleOpenDevicePage: function (subID) {
      const info = this.dasDevices.getDeviceInfo(subID);
      if (info == undefined) {
        ElMessage.warning("Invalid sub-device");
        return
      }
      if (info.invalidDeviceTypeName) {
        ElMessage.warning("Invalid device type");
        return
      }
      let pageVersion = getPageVersion();
      window.open(`/#/${pageVersion}/${subID}/overview/das_topo`);
    },
    handleDeleteAll: async function () {
      const isFacotryMode = await this.dev.funcs.isFacotryMode();
      if (!isFacotryMode) {
        ElMessage.warning("Not working in factory mode");
        return;
      }
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: "Confirm to reset all device nodes",
        callback: ok => {
          if (ok) {
            self.dasTopo.deleteTopoRootNode();
          }
        }
      });
    },
    handleRefresh(force = false) {
      this.dasTopo.refreshTopo(force, true)
    },
  },
}
</script>

<style lang="scss" scoped>
.account-for-topo {
  position: absolute;
  top: 10px;
  right: 32px;
  z-index: 2000;
  -webkit-user-select: none;
  -moz-user-select: none;
  -khtml-user-select: none;
  -ms-user-select: none;
}

.item {
  margin-bottom: 12px
}

.logo-text {
  //font-family:Arial;//color:#254c88;color:#606266;font-size:16px;font-style:normal;font-weight:900;line-height:21px;letter-spacing:1px;margin-top:0px}.itemText,.itemText:focus{color:#606266;font-size:14px;font-weight:bold;&:hover {
  color: rgb(32, 160, 255);
}

.itemInfo {
  float: right;
  margin-right: 2%;
  color: #606266;
  font-size: 14px;
  font-weight: bold
}

.iconexp {
  font-size: 14px
}

.el-icon-success {
  color: #67C23A
}

.el-icon-error {
  color: #606266
}

.el-icon-warning {
  color: #F56C6C
}

.el-divider--vertical {
  display: inline-block;
  width: 1px;
  height: 90vh;
  margin: 0 8px;
  vertical-align: middle;
  position: relative
}

.el-row {
  margin-bottom: 20px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}

.sub-title {
  color: var(--el-color-info);
  font-size: 12px;
  margin-left: 12px;
}

.loading_device_progress {
  position: absolute;
  left: 0;
  top: 0;
  height: 2px;
  background-color: var(--el-color-primary);
  text-align: right;
  border-radius: 100px;
  line-height: 1;
  white-space: nowrap;
  transition: width .6s ease;
}
</style>
