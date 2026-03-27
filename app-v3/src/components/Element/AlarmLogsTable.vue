<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
    <el-col :span="12" v-if="viewMode !== provideKeys.viewModePrintValue"
      style="display: flex;justify-content: flex-end;">
      <div class="toolbar">
        <my-table-data-export-button :table-data="tableData" :table-props="tableProps" export-file-name="Alarm"
          :specialHandlingCallback="specialHandlingCallback" />
        <el-button plain @click="handleGetTableData()">Refresh</el-button>
      </div>
    </el-col>
    <el-col v-if="viewMode !== provideKeys.viewModePrintValue">
      <my-table-search-bar v-model="tableQuery" :loading="loading"/>
    </el-col>
  </el-row>
  <el-table v-loading="loading" :data="tableData" :border="false" style="width: 100%;" table-layout="auto" stripe
    :class="{ 'my-table-height': viewMode != provideKeys.viewModePrintValue }" :default-sort="tableDefaultSort">
    <el-table-column prop="EventTime" label="Event Time">
      <template #default="scope">
        <span v-if="scope.row['EventTime']">{{
          deviceTimestampToDayJs(scope.row["EventTime"]).format("YYYY-MM-DD HH:mm")
        }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="AlarmName" label="Alarm Name" />
    <el-table-column prop="AlarmSeverity" label="Alarm Severity" />
    <el-table-column prop="AlarmStatus" label="Alarm Status" />
    <el-table-column prop="DeviceTypeName" label="Device Type Name" />
    <el-table-column prop="DeviceSubID" label="Device Sub" />
    <el-table-column prop="SiteName" label="Site Name" />
    <el-table-column prop="DeviceName" label="Device Name" />
    <el-table-column prop="SerialNumber" label="Serial Number" />
    <el-table-column prop="SoftwareVersion" label="Software Version" />
  </el-table>
  <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="tableDataTotal"
    layout="prev, pager, next, sizes, jumper, total" :page-sizes="[20, 50, 100]" :pager-count="5"
    style="margin-top: 4px;" />
</template>

<script>
import { useAppStore } from '@/stores/app';
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'
import model from "@/stores/model"
import AlarmLogsTableConfig from './AlarmLogsTableConfig';
import { useDasModel } from '@/stores/das-model';

export default {
  name: "MyAlarmLogsTable",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dasModel = useDasModel();
    return {
      appStore,
      dasDevices,
      provideKeys,
      dasModel,
      deviceTimestampToDayJs: model.deviceTimestampToDayJs,
    };
  },
  data() {
    let cpQueryItems = [...(AlarmLogsTableConfig.TableQueryItem)];
    cpQueryItems[0].options = this.dasModel.deviceTypeNames;
    return {
      loading: false,
      tableQuery: cpQueryItems,
      tableProps: AlarmLogsTableConfig.TableProps,
      tableDefaultSort: AlarmLogsTableConfig.TableDefaultSort,
      currentPage: 1,
      pageSize: 50,
    }
  },
  computed: {
    tableData() {
      return this.dasDevices.currentDevice.alarm.showAlarmLogs;
    },
    tableDataTotal() {
      return this.dasDevices.currentDevice.alarm.queryAlarmLogs.length;
    },
  },
  watch: {
    tableQuery() {
      if (this.currentPage != 1) {
        this.currentPage = 1;
      } else {
        this.handleGetTableData();
      }
    },
    currentPage() {
      this.handleGetTableData();
    },
    pageSize() {
      this.handleGetTableData();
    }
  },
  mounted() {
    this.initPage();
  },
  methods: {
    initPage: async function () {
      this.loading = true;
      await this.dasDevices.currentDevice.alarm.getAlarmLogList({ pageNumber: this.currentPage, pageSize: this.pageSize, query: this.tableQuery });
      this.loading = false;
      this.tableQuery[3].options = this.dasDevices.currentDevice.alarm.alarmSeverity;
      this.tableQuery = [...this.tableQuery];
    },
    handleGetTableData: async function () {
      this.loading = true;
      await this.dasDevices.currentDevice.alarm.getAlarmLogList({ pageNumber: this.currentPage, pageSize: this.pageSize, query: this.tableQuery });
      this.loading = false;
    },
    specialHandlingCallback({ key, value }) {
      if (key == 'EventTime' && value) {
        return this.deviceTimestampToDayJs(value).format("YYYY-MM-DD HH:mm")
      }
      return value;
    },
  },
};
</script>
<style lang="scss" scoped></style>