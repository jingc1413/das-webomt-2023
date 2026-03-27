<template>
  <el-table :data="tableData" table-layout="auto" style="max-width: 100%;" :show-header="false">
    <el-table-column prop="Name" label="Name" min-width="240"/>
    <el-table-column prop="Value" label="Value">
      <template #default="scope">
        <report-form-layout v-if="scope.row['Value'] && scope.$index !== 0" :owner="owner" :data="scope.row['Value']" />
        <span v-else-if="scope.$index == 0">{{scope.row['Value']}}</span>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
export default {
  name: 'ReportToolbarForm',
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    return {
    };
  },
  data() {
    function getColItem(colItems = []) {
      let temp = []
      colItems.forEach(colItem => {
          temp.push({
            Name: colItem.Name ?? '',
            Value: colItem,
          })
      })
      return temp
    }

    const tableData = [];
    if (this.data.length) {
      this.data.forEach(col => {
        tableData.push(...(getColItem([col])));
      })
    }

    tableData.unshift({Name:'Name', Value:'Value'})
    return {
      tableData: tableData
    };
  },
  methods: {
  }

}
</script>

<style lang="scss" scoped></style>
