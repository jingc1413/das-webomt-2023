<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
    <el-col :span="12" v-if="viewMode !== provideKeys.viewModePrintValue"
      style="display: flex;justify-content: flex-end;">
      <div class="toolbar">
        <my-table-data-export-button :table-data="tableData" :table-props="inventoryTableProps" export-file-name="Inventory" />
        <el-button plain @click="handleGetTableData()">Refresh</el-button>
      </div>
    </el-col>
    <el-col v-if="viewMode !== provideKeys.viewModePrintValue">
      <my-table-search-bar v-model="tableQuery"/>
    </el-col>
  </el-row>

  <el-table v-loading="loading" :data="tableData" :border="false" style="width: 100%;" table-layout="auto" stripe
    :class="{ 'my-table-height': viewMode != provideKeys.viewModePrintValue }">
    <el-table-column prop="SubID" label="Device Sub ID" />
    <el-table-column prop="DeviceTypeName" label="Device Type" />
    <el-table-column prop="InstalledLocation" label="Device Location" />
    <el-table-column prop="DeviceName" label="Device Name" />
    <el-table-column prop="ElementSerialNumber" label="Element Serial Number" />  
    <el-table-column prop="ElementModelNumber" label="Element Model Number" />
    <el-table-column prop="Version" label="Software Version" />

    <el-table-column prop="ConnectState" label="Status" />
    <el-table-column prop="IpAddress" label="IP" />
    <el-table-column prop="RouteAddress" label="Route" />

    <el-table-column prop="LifeTime" label="Life Time" />
    <el-table-column prop="ElementOperatingTemperature" label="Element Operating Temperature" />

  </el-table>
</template>

<script>
import { useAppStore } from '@/stores/app';
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'
import InVentoryTableConfig from "./InVentoryTableConfig";
import { useDasModel } from '@/stores/das-model';

export default {
  name: "MyInventoryTable",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dasModel =useDasModel();
    return {
      appStore,
      dasDevices,
      provideKeys,
      dasModel,
    };
  },
  data() {
    let cpQueryItems = [...(InVentoryTableConfig.InVentoryTableQueryItem)];
    cpQueryItems[0].options = this.dasModel.deviceTypeNames;
    return {
      loading: false,
      tableQuery: cpQueryItems,
      inventoryTableProps: InVentoryTableConfig.InventoryTableProps,
    }
  },
  computed: {
    tableData() {
      const infos = this.dasDevices.getDeviceInfos.filter(item=>{
        return this.tableQuery.some(queryItem=>{
          try {
            return queryItem.value && (String(item[queryItem.key])?.indexOf(queryItem.value) == -1) // find unmatched
          } catch (error) {
            console.log({error});
          }
          return false;
        }) == false;
      });
      return infos;
    }
  },
  mounted() {
    this.handleGetTableData();
  },
  methods: {
    handleGetTableData: async function () {
      this.loading = true;
      await this.dasDevices.updateDeviceInfos(false, true);
      this.loading = false;
    },
  },
};
</script>
<style lang="scss" scoped></style>