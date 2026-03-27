<!-- eslint-disable vue/no-v-for-template-key -->
<template>
  <el-button type="primary" plain @click="doExport()">Export</el-button>
</template>

<script>
import { useDasDevices } from "@/stores/das-devices";
import { read, utils, writeFile } from "xlsx";
export default {
  name: "TableDataExportButton",
  props: {
    tableData: {
      type: Array,
      default: () => [],
    },
    tableProps:{
      type: Object,
      default: () => {},
    },
    specialHandlingCallback: {
      type: Function,
      default: ({key,value})=>{return value},
    },
    exportFileName: {
      type: String,
      default: "",
    },
  },
  setup() {
    let dasDevices = useDasDevices()
    return {
      dasDevices
    }
  },
  data() {
    return {
    }
  },
  mounted() {
  },
  methods: {
    doExport() {
      if (this.tableData.length == 0) {
        return this.$msgModal.msgWarning("Exported content is no data");
      }
      let rows = this.tableData.map((row) => {
        let body = {};
        for (let key in this.tableProps) {
            let value = row[key];
            let name = this.tableProps[key];
            console.log({key, value, name});
            body[name] = this.specialHandlingCallback({key, value});
        }
        return body;
      });
      const worksheet = utils.json_to_sheet(rows);
      const workbook = utils.book_new();
      utils.book_append_sheet(workbook, worksheet, "Data");
      const filename = this.dasDevices.currentDeviceFileName(
        this.exportFileName, false)
        + ".xlsx";
      writeFile(workbook, filename, { compression: true });
    },
  }
}
</script>

<style lang="scss" scoped>
.table-search-form {
  // background: #eee;
  margin: 8px 0;
  padding: 8px 12px 0 12px;
  // border-radius: 4px;
  // box-shadow: 0px 0px 6px rgba(0, 0, 0, .12);
  border-bottom: 1px solid #eee;
}
</style>