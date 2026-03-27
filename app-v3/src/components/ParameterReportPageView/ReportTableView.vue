<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <report-firmware-table v-if="data.Type === 'Table:Firmwares'" :owner="owner" :data="data" />
  <template v-else>
    <el-row v-if="data.Name" style="height: 48px;">
      <el-col :span="12">
        <h4><span>{{ data.Name }}</span></h4>
      </el-col>
    </el-row>
    <report-toolbar-form  v-if="toolbarItems?.length > 0" :owner="owner" :data="toolbarItems" style="margin: 10px 0;"></report-toolbar-form>
    <el-table :data="filterTableData" table-layout="auto" :show-header="filterTableData.length == 0">
      <template v-for="column in data.Items">
        <el-table-column v-if="column.Key.toLocaleLowerCase() != 'action'" :key="column.Key" :prop="column.Key"
          :label="column.Name">
          <template #default="scope">
            <report-element v-if="scope.row[column.Key] && scope.$index !== 0" :owner="owner" :data="scope.row[column.Key]" />
            <span v-else-if="scope.$index == 0">{{column.Name}}</span>
          </template>
        </el-table-column>
      </template>
    </el-table>
  </template>
</template>
  
<script>
import { useDasDevices } from '@/stores/das-devices';
import ReportFirmwareTable from './ReportFirmwareTable.vue'
export default {
  name: "ReportTable",
  components: { ReportFirmwareTable },
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
    const tableInvalidKey = this.data?.Style?.invalidKey;
    const tableInvalidValue = this.data?.Style?.invalidValue;
    const toolbarItems = [];
    if (this.data?.Actions?.toolbar?.Items) {
      this.data?.Actions?.toolbar?.Items.forEach(v => {
        if (!(v.Type === "Button" && v.Actions?.click?.Name === "AddItem") && v.Access != 'wo') {
          toolbarItems.push(v);
        }
      })
    }

    const tableColumns = this.data?.Items?.filter(v => {
      if (v.Key == "action") {
        return false;
      }
      return true;
    }) || [];
    const tableData = this.data.Data;
    return {
      tableInvalidKey,
      tableInvalidValue,
      toolbarItems,
      tableColumns: tableColumns,
      tableData: tableData,
    };
  },
  computed: {
    filterTableData() {
      if (this.tableData.length == 0) {
        return []
      }
      if (!this.tableInvalidKey) {
        return ['',...this.tableData];
      }
      let temp = this.tableData.filter(row => {
        const item = row[this.tableInvalidKey];
        if (item) {
          const param = this.dev.params.getParam(item.OID)
          if (param && param.Value) {
            if (!this.tableInvalidValue) {
              return true;
            }
            if (param.Value !== this.tableInvalidValue) {
              return true;
            }
          }
        }
        return false;
      });
      if (temp.length == 0) {
        return []
      }
      return ['', ...temp]
    },
  },
};
</script>
  