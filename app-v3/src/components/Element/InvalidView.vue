<template>
  <div class="invalid_view" v-show="showInvalidView">
    <el-row justify="center" style="margin-top: calc( 50vh - 200px );">
      <el-result icon="warning" title="The device is not available."
        sub-title="If the device is rebooting, please wait for the device to finish starting up.">
        <template #extra>
          <el-row v-loading="loading" element-loading-text="Checking..." style="height: 100px; width: 100px;">
          </el-row>
        </template>
      </el-result>
    </el-row>
  </div>
</template>

<script>

import { useRoute } from "vue-router";
import { useDasDevices } from "@/stores/das-devices";
import settings from "@/settings";
import { useAppStore } from "@/stores/app";

export default {
  name: 'MyInvalidView',
  setup() {
    const dasDevices = useDasDevices();
    const appStore = useAppStore();
    return {
      appStore,
      dasDevices,
    }
  },
  computed: {
    sub() {
      return useRoute().params.sub || "local";
    },
    deviceInfo() {
      return this.dasDevices.getDeviceInfo(this.sub);
    },
    showInvalidView() {
      if (!this.appStore.isAppSetupDone) {
        return false;
      }
      if (settings.nodeTest) {
        return false;
      }
      if (this.deviceInfo?.state?.available != undefined) {
        return this.deviceInfo?.state?.available !== true;
      }
      return false;
    },
    loading() {
      return this.deviceInfo?.state?.availableChecking === true;
    },
  },
}

</script>

<style scoped>
.invalid_view {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: #fff;
  z-index: 10002;
}
</style>