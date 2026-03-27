<template>
  <el-row v-if="data.Name" style="height: 48px;">
    <el-col :span="12">
      <h4><span>{{ data.Name }}</span></h4>
    </el-col>
  </el-row>
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
  name: 'ReportForm',
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
        if (colItem.Type != 'Button') {
          temp.push({
            Name: colItem.Name ?? '',
            Value: colItem,
            
          })
        }
      })
      return temp
    }

    const tableData = [];
    if (this.data.Items) {
      this.data.Items.forEach(col => {
        if (col?.Type?.startsWith('Layout:Col')) {
          tableData.push(...(getColItem(col.Items)));
        } else if ((col?.Type?.startsWith('Layout:Row'))) {
          col.Items.forEach(rowItem => {
            if (rowItem.Type?.startsWith('Layout:Col')) {
              tableData.push(...(getColItem(rowItem.Items)));
            } else {
              tableData.push(...(getColItem([rowItem])));
            }
          })
        } else {
          tableData.push(...(getColItem([col])));
        }
      })
    }

    tableData.unshift({Name:'Name', Value:'Value'})
    return {
      tableData: tableData.filter(item=>{
        let flag = true;
        try {
          flag = item.Value.Access != 'wo'
        } catch (error) {
          console.error(error);
        }
        return flag;
      }),
    };
  },
  methods: {
  }

}
</script>

<style lang="scss" scoped></style>
