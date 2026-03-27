<template>
  <el-row style="padding: 20px 20px 0px 20px;background: #F2F3F5;min-height: 100%;">
    <el-col>
      <el-result :title="upgradeTitle">
        <template #icon>
          <el-icon :size="120">
            <SuccessFilled v-if="upgradeStatus == 'success'" class="icon-success" />
            <WarningFilled v-else-if="upgradeStatus == 'warning'" class="icon-warning" />
            <CircleCloseFilled v-else-if="upgradeStatus == 'error'" class="icon-error" />
            <InfoFilled v-else-if="upgradeStatus == 'info'" class="icon-info" />
          </el-icon>
        </template>
        <template #sub-title>
          <p>
            {{ upgradeSubTitle }}<br>
            {{ upgradeState?.filename }}<br>
            <my-time-piece />
          </p>
        </template>
        <template #extra>
          <template v-if="upgradeFinished">
            <el-row style="width: calc(100vw - 60px);">
              <el-col :span="12" :offset="6" style="margin-top: 12px;">
              </el-col>
            </el-row>
          </template>
          <template v-else>
            <el-row style="width: calc(100vw - 60px);">
              <el-col :span="16" :offset="4" style="margin-bottom: 12px;">
                <el-progress :percentage="upgradeState?.progress" :stroke-width="10" striped striped-flow
                  :duration="5" />
              </el-col>

              <template v-if="upgradeState?.stats !== undefined">
                <el-col :span="4" :offset="6">
                  <el-statistic title="Total devices" :value="upgradeState?.stats.totalDeviceNum || '-'" />
                </el-col>
                <el-col :span="4">
                  <el-statistic title="Failed devices" :value="upgradeState?.stats.failedDeviceNum || '-'" />
                </el-col>
                <el-col :span="4">
                  <el-statistic title="Timeout devices" :value="upgradeState?.stats.timeoutDeviceNum || '-'" />
                </el-col>
              </template>

              <el-col :span="12" :offset="6" style="margin-top: 12px;">
                <el-scrollbar v-if="upgradeState?.logs"
                  style="height:calc(100vh - 20px - 80px - 46px - 29px - 125px - 30px - 50px - 24px - 20px); text-align: left;">
                  <template v-for="(item, index) in upgradeState.logs">
                    <p v-if="item.status === 'success'" :key="index + 'success'">
                      <el-icon style="margin-right: 32px;">
                        <SuccessFilled class="icon-success" />
                      </el-icon>
                      {{ item.text }}
                    </p>
                    <p v-else-if="item.status === 'warning'" :key="index + 'warning'">
                      <el-icon style="margin-right: 32px;">
                        <WarningFilled class="icon-warning" />
                      </el-icon>
                      {{ item.text }}
                    </p>
                    <p v-else-if="item.status === 'error'" :key="index + 'error'">
                      <el-icon style="margin-right: 32px;">
                        <CircleCloseFilled class="icon-error" />
                      </el-icon>
                      {{ item.text }}
                    </p>
                    <p v-else :key="index + 'info'">
                      <el-icon style="margin-right: 32px;">
                        <InfoFilled class="icon-info" />
                      </el-icon>{{ item.text }}
                    </p>
                  </template>
                </el-scrollbar>
              </el-col>
            </el-row>
          </template>


        </template>
      </el-result>
    </el-col>
  </el-row>
</template>

<script>
import { useAppStore } from '@/stores/app';
import { useDasDevices } from '@/stores/das-devices';
import { dayjs } from 'element-plus';

export default {
  name: 'UpgradeView',
  setup() {
    const dev = useDasDevices().currentDevice;
    const upgradeState = dev.upgrade.upgradeState;
    return {
      upgradeState,
      dev,
    }
  },
  data() {
    return {
      upgradeFinished: false,
      upgradeFinishedTime: undefined,
    }
  },
  computed: {
    upgradeStatus() {
      if (this.upgradeFinished) {
        return this.upgradeState?.resultStatus;
      }
      return "info";
    },
    upgradeTitle() {
      if (this.upgradeFinished) {
        switch (this.upgradeState?.resultStatus) {
          case "success":
            return "Upgrade completed"
          case "error":
            return "Upgrade failed"
        }
        return "Upgrade complete"
      }
      return "Upgrading";
    },
    upgradeSubTitle() {
      if (this.upgradeFinished) {
        return `${this.upgradeState?.resultText}`;
      }
      return `Please keep the device powered on and wait for the upgrade to complete.`;
    },
  },
  created() {
    if (this.upgradeState === undefined) {
      useAppStore().gotoDefaultRoute();
    }
  },
  mounted() {
    this.doCheckUpgrading();
  },
  methods: {
    doCheckUpgrading: async function () {
      const self = this;
      if (this.upgradeState === undefined) {
        useAppStore().gotoDefaultRoute();
        return;
      }

      await this.dev.upgrade.checkUpgrading();
      if (this.upgradeState.finished) {
        const now = dayjs();
        if (this.upgradeFinishedTime === undefined) {
          this.upgradeFinishedTime = now;
        }
        if (now.diff(this.upgradeFinishedTime) > 10000) {
          this.upgradeFinished = true;
        }
      } else {
        this.upgradeFinishedTime = undefined;
        this.upgradeFinished = false;
      }

      if (this.upgradeFinished) {
        setTimeout(() => {
          self.dev.upgrade.clearUpgradeState();
          useAppStore().gotoDefaultRoute();
        }, 10000)
        return;
      }

      setTimeout(() => {
        self.doCheckUpgrading()
      }, 1000)
    }
  },

}
</script>

<style scoped></style>
