<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
  </el-row>
  <el-table :data="tableData" style="width: 100%;" table-layout="auto">
    <el-table-column prop="Name" label="Firmware" />
    <el-table-column prop="CRC" label="CRC" />
  </el-table>
</template>
  
<script>
import { useDasDevices } from '@/stores/das-devices';


export default {
  name: "MyFirmwaresTable",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
    };
  },
  data() {
    return {
      loading: false,
      deleteConfirmDialogVisible: false,
      deleteFirmwareData: {},
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
    handleGetFirmwareList: async function () {
      this.loading = true;
      await this.dev.firmwares.getFirmwareList();
      this.loading = false;
    },
  },
};
</script>
<style lang="scss" scoped></style>