<template>
  <el-row style="padding: 20px 20px 0px 20px;background: #F2F3F5;min-height: 100%;">
    <el-col>
      <el-row align="middle" style="width: calc(100vw - 60px); height: 100%;">
        <el-col :span="16" :offset="4">
          <el-progress :percentage="loadingState?.progress" :stroke-width="10" striped striped-flow :duration="5"
            style="margin-bottom: 12px;" />
          <el-text style="display: block;text-align: center;">{{ loadingSubTitle }}</el-text>
        </el-col>
      </el-row>
    </el-col>
  </el-row>
</template>

<script>

import { useAppStore } from '@/stores/app';
import { useDasModel } from '@/stores/das-model';
import { useDasDevices } from '@/stores/das-devices';
import { useAuthStore } from '@/stores/auth';

export default {
  name: 'LoadingView',
  setup() {
    let appStore = useAppStore();
    let dasModel = useDasModel();
    let dasDevices = useDasDevices();
    let authStore = useAuthStore();
    return {
      appStore,
      dasModel,
      dasDevices,
      authStore,
    };
  },
  data() {
    return {
      loadingState: {
        resultStatus: '',
        resultText: '',
        progress: 0,
        logs: []
      }
    };
  },
  computed: {
    loadingSubTitle() {
      return `${this.loadingState?.resultText}`;
    },
  },
  created() {
  },
  mounted() {
    this.startLoadingData();
  },
  unmounted() {
    this.loadingState.progress = 100;
  },
  methods: {
    async startLoadingData() {
      this.loadingState.resultStatus = "info";
      this.loadingState.resultText = "Setup device models ...";
      this.loopAdProgress(0, 30);
      await this.dasModel.setup();
      this.loadingState.resultText = "Setup devices ...";
      this.loopAdProgress(30, 60);
      await this.dasDevices.setup();
      this.loadingState.resultText = "Done.";
      this.loadingState.progress = 100;
      this.appStore.setAppSetupDone();
      setTimeout(() => {
        this.appStore.gotoDefaultRoute();
      }, 300);
    },
    loopAdProgress(startProgress = 1, maxProgress = 99) {
      this.$nextTick(() => {
        if (this.loadingState.progress < startProgress) {
          this.loadingState.progress = startProgress;
        }
      });
      setTimeout(() => {
        if (this.loadingState.progress < maxProgress) {
          this.loadingState.progress++;
          this.loopAdProgress(startProgress, maxProgress);
        }
      }, 1000);
    },
  },
}
</script>

<style scoped></style>
