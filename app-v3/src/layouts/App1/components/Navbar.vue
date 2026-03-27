<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <div class="navbar">
    <!-- <hamburger :is-active="appStore.sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" /> -->

    <span class="app_title">{{ appTitle }}</span>

    <el-popover placement="bottom" trigger="click" width="320" ref="deviceTreeRef">
      <template #reference>
        <el-button text class="device_tree_button">
          <el-icon>
            <CaretBottom />
          </el-icon>
        </el-button>
      </template>
      <template #default>
        <my-device-tree-selected @selectDevice="handleRouteToDevice" />
      </template>
    </el-popover>

    <breadcrumb class="breadcrumb-container" />

    <div class="right-menu">
      <el-text v-if="runTime" class="mx-1">
        SiteName
        <span style="font-weight: 700; font-size: 13px;">
          {{ siteName }}
        </span>
      </el-text>
      <el-text v-if="runTime" class="mx-1" style="margin-left:8px;">
        UpTime
        <span style="font-weight: 700; font-size: 13px;">
          {{ runTime.fromNow(true) }}
        </span>
      </el-text>
      <el-text v-if="serveTime" class="mx-1" style="margin-left:8px; margin-right:30px">
        DeviceTime
        <span style="font-weight: 700; font-size: 13px;">
          {{ serveTime.format('YYYY-MM-DD HH:mm:ss') }}
        </span>
      </el-text>
      <el-dropdown class="avatar-container">
        <div class="avatar-wrapper">
          <el-button link>
            <el-icon size="20px">
              <Menu />
            </el-icon>
          </el-button>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="snapshotDialogVisible = !snapshotDialogVisible">Snapshot</el-dropdown-item>
          </el-dropdown-menu>

          <el-dropdown-menu>
            <el-dropdown-item @click="reportDialogVisible = !reportDialogVisible">Report</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-dropdown class="avatar-container">
        <div class="avatar-wrapper">
          <el-button link>
            <span class="navbar_button_size">{{ auth.loginUserName || "----" }}</span>
          </el-button>
          <!-- <span></span> -->
          <!--          <el-icon class="el-icon&#45;&#45;right"><i-ep-arrow-down /></el-icon>-->
        </div>
        <template #dropdown>
          <el-dropdown-menu class="user-dropdown">
            <!-- <router-link to="/">
              <el-dropdown-item>
                Home
              </el-dropdown-item>
            </router-link>
            <a target="_blank" href="https://github.com/PanJiaChen/vue-admin-template/">
              <el-dropdown-item>Github</el-dropdown-item>
            </a>
            <a target="_blank" href="https://panjiachen.github.io/vue-element-admin-site/#/">
              <el-dropdown-item>Docs</el-dropdown-item>
            </a> -->
            <el-dropdown-item @click="handleChangePassword" v-if="$permissions.hasPermission(['api.current.change-password'])">
              <span style="display:block;">Change Password</span>
            </el-dropdown-item>
            <el-dropdown-item @click="handleLogout">
              <span style="display:block;">Log Out</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <snapshot-page-view v-if="snapshotDialogVisible" v-model:snapshotDialogVisible="snapshotDialogVisible" />
    <report-page-view v-if="reportDialogVisible" v-model:reportDialogVisible="reportDialogVisible" />
    <my-change-password v-if="passwordDialogVisible" :isOpen="passwordDialogVisible" @dialogHide="closeChangePassword" />
    <my-config-wizard v-if="wizardDialogVisible" :isOpen="wizardDialogVisible" @dialogHide="closeWizardDialog" />
  </div>
</template>

<script>
import Breadcrumb from '@/components/Breadcrumb/Breadcrumb.vue'

import SnapshotPageView from "@/components/SnapshotPageView/SnapshotPageView.vue";
import ReportPageView from '@/components/ParameterReportPageView/ReportPageView.vue';
import { useAppStore } from '@/stores/app'
import { useAuthStore } from "@/stores/auth";
import { useIntervalFn } from '@vueuse/core'
import { dayjs, ElMessage } from 'element-plus';
import { useDasDevices } from '@/stores/das-devices';
import model from '@/stores/model';



export default {
  components: {
    Breadcrumb,
  },
  setup() {
    let relativeTime = require('dayjs/plugin/relativeTime');
    dayjs.extend(relativeTime);
    const appStore = useAppStore();
    const auth = useAuthStore();
    const dasDevices = useDasDevices();
    return { appStore, auth, dasDevices, dayjs }
  },
  data() {
    return {
      serveTime: null,
      runTime: null,
      snapshotDialogVisible: false,
      reportDialogVisible: false,
      passwordDialogVisible: false,
      wizardDialogVisible: false
    }
  },
  computed: {
    appTitle() {
      return this.appStore.appTitle;
    },
    systemTimeDiff() {
      const info = this.dasDevices.currentDeviceInfo;
      if (info.SystemTime > 0) {
        const systemTime = model.deviceTimestampToDayJs(info.SystemTime);
        const diff = dayjs().diff(systemTime, 'second');
        return diff;
      }
      return undefined;
    },
    upTime() {
      const info = this.dasDevices.currentDeviceInfo;
      if (info.UpTime > 0) {
        return dayjs().subtract(info.UpTime, 'second');
      }
      return undefined;
    },
    siteName() {
      const info = this.dasDevices.currentDeviceInfo;
      return info.SiteName;
    }
  },
  mounted() {
    this.setupServeTimeAndRunTime()
  },
  methods: {
    toggleSideBar() {
      this.appStore.toggleSideBar()
    },
    setupServeTimeAndRunTime() {
      useIntervalFn(() => {
        if (this.systemTimeDiff) {
          this.serveTime = dayjs().subtract(this.systemTimeDiff, 'second');
        }
        this.runTime = this.upTime;
      }, 1000, { immediate: true })
      useIntervalFn(() => {
        this.dasDevices.updateCurrentDeviceInfo();
      }, 600000, { immediate: true })
    },
    handleLogout() {
      this.auth.logout("Logout");
    },
    handleChangePassword() {
      this.passwordDialogVisible = true;
    },
    closeChangePassword(change) {
      this.passwordDialogVisible = false;
      if (change) {
        this.appStore.openConfirmDialog({
          title: 'Tip',
          content: "Your password has been change,please login again",
          supportCancel: false,
          callback: ok => {
            if (ok) {
              this.authStore.logout()
            }
          }
        })
      }

    },
    handleOpenWizardDialog() {
      this.wizardDialogVisible = true;
    },
    closeWizardDialog(change) {
      this.wizardDialogVisible = false;
    },
    handleRouteToDevice(deviceInfo) {
      let SubID = deviceInfo.SubID;
      let { module, sub, page } = this.$route.params;
      if (SubID == sub) return;
      if (deviceInfo.ConnectState < 6) {
        this.$refs['deviceTreeRef'].hide();
        this.$router.replace({
          name: "app1Page",
          params: {
            sub: SubID,
          }
        })
      } else {
        ElMessage.warning("Device is offline");
      }
    }

  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: var(--header-height);
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, .08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .app_title {
    color: #606266;
    font-size: 16px;
    font-style: normal;
    font-weight: 900;
    letter-spacing: 1px;
    float: left;
    line-height: var(--header-height);
    margin-right: 32px;
    margin-left: 12px;
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: var(--header-height);

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 16px;
      height: 40px;
      line-height: 40px;

      .avatar-wrapper {
        margin-top: 5px;
        position: relative;

        .user-avatar {
          cursor: pointer;
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }

  .runtime_status {
    background-color: var(--el-color-success-light-3);
    border-radius: 8px;
    display: inline-block;
    width: 12px;
    height: 12px;
    margin-right: 8px;
  }
}
.device_tree_button {
  font-size: 16px;
  float: left;
  line-height: var(--header-height);
  height: var(--header-height);
}
.navbar_button_size {
  font-size: 14px;
}

</style>
