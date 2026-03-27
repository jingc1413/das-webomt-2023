<template>
  <el-row>
    <el-col>
      <el-form ref="infoFormRef" :model="jobInfo" :disabled="jobLoading" :inline="true"
        style="margin: 18px 10px 18px 18px;" label-width="auto" label-position="right">

        <el-form-item label="Address" prop="Address" :rules="[{ required: true, trigger: 'blur' }]">
          <el-input v-model.trim="jobInfo.Address" type="text" style="width: 220px;" clearable />
        </el-form-item>

        <el-form-item label="Interval" prop="Interval" :rules="[{ required: true, trigger: 'blur' }]">
          <el-input-number v-model.trim="jobInfo.Interval" :step="1" :min="1" :max="15" :step-strictly="true"
            :controls="false" style="width: 80px;" />
        </el-form-item>

        <el-form-item label="Count" prop="Count" :rules="[{ required: true, trigger: 'blur' }]">
          <el-input-number v-model.trim="jobInfo.Count" :step="1" :min="1" :max="100" :step-strictly="true"
            :controls="false" style="width: 80px;" />
        </el-form-item>


        <!-- <el-form-item label="Timeout" prop="Timeout" :rules="[{ required: true, trigger: 'blur' }]">
          <el-input-number v-model.trim="jobInfo.Timeout" :step="1" :min="1" :max="60" :step-strictly="true"
            :controls="false" style="width: 80px;" />
        </el-form-item> -->

        <el-form-item v-show="!jobRunning">
          <div class="dialog-footer" v-hasPermissionAnd="['api.diag.ping.jobs.create']"
            style="display: flex;align-items: center;">
            <el-button type="primary" :loading="jobLoading" @click="submitPingTest()">Ping</el-button>
          </div>
        </el-form-item>

        <el-form-item v-show="jobRunning">
          <div class="dialog-footer" style="display: flex;align-items: center;">
            <el-button type="primary" :loading="jobLoading" @click="cancelPingTest()">Cancel</el-button>
          </div>
        </el-form-item>

      </el-form>
    </el-col>

    <el-col>
      <div class="mainBox">
        <el-scrollbar style="height:calc(100vh - 226px);">
          <div class="boxBorderCls">
            <div class="table-main-con">
              <template v-for="  messageItem   in   messageList  ">
                <pre style="margin: 8px 0;">{{ messageItem }}</pre>
              </template>
            </div>
          </div>
        </el-scrollbar>
      </div>
    </el-col>
  </el-row>
</template>

<script >

import { useCurrentPing } from '@/stores/current-ping.js';

export default {
  name: 'MyPingDiagPage',
  setup(prop) {
    let currentPing = useCurrentPing();
    return {
      currentPing
    }
  },
  data() {
    return {
      jobInfo: {
        Address: '',
        Count: 4,
        Interval: 1,
        Timeout: -1
      },
      jobLoading: false,
    }
  },
  computed: {
    jobRunning() {
      return this.currentPing.jobStats.IsRunning;
    },
    messageList() {
      return this.currentPing.getWsMessageList;
    }
  },
  methods: {
    submitPingTest() {
      this.$refs['infoFormRef'].validate((valid, fields) => {
        if (valid) {
          this.jobLoading = true;
          this.currentPing.createPingJob({ ...this.jobInfo }).then(res => {
            if (res) {
              this.currentPing.openWebsocket(this.currentPing.token, {
                openWebsocketCallback: () => {
                  this.currentPing.runPingJob(this.currentPing.token)
                }
              })
            }
          }).finally(() => {
            this.jobLoading = false;
          })
        }
      })
    },
    cancelPingTest() {
      this.jobLoading = true;
      this.currentPing.stopPingJob(this.currentPing.token).then(res => {
        console.log(res)
      }).catch(e => {
        console.log(e)
      }).finally(() => {
        this.jobLoading = false;
      });
    }
  }
}


</script>

<style lang='scss' scoped>
.mainBox {
  height: calc(100vh - 216px);
  max-width: 800px;
  width: 64%;
  margin: 0 10px;
  padding: 8px 8px 0;
  overflow: auto;
  background-color: #f4f4f5;
  border-top: 2px Solid #a0cfff;
  box-shadow: var(--el-box-shadow-light);

  .boxBorderCls {
    min-height: calc(100vh - 226px);
    border-radius: 4px;
    overflow: auto;
  }
}
</style>