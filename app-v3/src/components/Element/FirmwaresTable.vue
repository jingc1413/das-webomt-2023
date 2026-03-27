<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
    <el-col :span="12" v-if="viewMode !== provideKeys.viewModePrintValue"
      style="display: flex;justify-content: flex-end;">
      <div class="toolbar">
        <el-button type="primary" plain @click="handleGetFirmwareList(true)">Refresh</el-button>
      </div>
    </el-col>
  </el-row>
  <el-table v-loading="loading" :data="tableData" :border="false" style="width: 100%;" table-layout="auto" stripe
    :class="{ 'my-table-height': viewMode != provideKeys.viewModePrintValue }">
    <el-table-column prop="Name" label="Firmware" />
    <el-table-column prop="CRC" label="CRC" />
    <el-table-column prop="Action" label="" v-if="viewMode != provideKeys.viewModePrintValue">
      <template #default="scope">
        <template v-if="scope.row['Name']">
          <el-button @click="deleteFirmware(scope.row)" style="margin-left: 12px;">Delete</el-button>
        </template>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import { useAppStore } from '@/stores/app';
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: "MyFirmwaresTable",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const appStore = useAppStore();
    const dev = useDasDevices().currentDevice;
    return {
      appStore,
      dev,
      provideKeys,
    };
  },
  data() {
    return {
      loading: false,
    }
  },
  computed: {
    tableData() {
      return this.dev.firmwares.firmwares;
    }
  },
  mounted() {
    this.handleGetFirmwareList();
  },
  methods: {
    handleGetFirmwareList: async function (showMessage=false) {
      this.loading = true;
      await this.dev.firmwares.getFirmwareList(showMessage);
      this.loading = false;
    },
    deleteFirmware: async function (firmware) {
      const self = this;
      this.appStore.openConfirmDialog({
        title: 'Confirm',
        content: "Confirm to delete firmware " + firmware.Name,
        callback: ok => {
          if (ok) {
            self.appStore.openConfirmDialog({
              title: 'Confirm to delete the firmware',
              content: 'Please input force code',
              needInput: true,
              inputRule: [
                {
                  required: true,
                  validator: (rule, value, callback) => {
                    if (value !== "Image") {
                      return callback(new Error('Please input force code'));
                    } else {
                      callback();
                    }
                  },
                  trigger: 'blur',
                }
              ],
              callback: ok => {
                if (ok) {
                  self.doConfirmDeleteFirmwareStatus(firmware);
                }
              }
            });
          }
        },
      })

    },
    doConfirmDeleteFirmwareStatus: async function (firmware) {
      this.loading = true;
      await this.dev.firmwares.deleteFirmware(firmware)
      this.loading = false;
    },
  },
};
</script>
<style lang="scss" scoped></style>