<template>
  <el-table ref="tableRef" :data="tableData" table-layout="auto" style="padding-bottom: 16px;" @selection-change="changeSelection" row-key="oid">
    <el-table-column type="selection" width="55" :selectable="handleTableRowSelectable"/>
    <el-table-column type="index" width="55" />
    <el-table-column label="Parameter Name" prop="">
      <template #default="scoped">
        <report-param :data="{OID:scoped.row.oid}" pageType="name" />
      </template>
    </el-table-column>
    <el-table-column label="Current Value" prop="">
      <template #default="scoped">
        <report-param :data="{OID:scoped.row.oid}" />
      </template>
    </el-table-column>
    <el-table-column label="Update Value" prop="">
      <template #default="scoped">
        <div v-if="scoped.row.value !== undefined" :class="{'cell-changed': scoped.row.isChanged, 'cell-invalid': !scoped.row.isValidate}" style="padding: 0 4px;">
          <report-param :data="{OID:scoped.row.oid, Value: scoped.row.value}" />
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup>
import { defineProps, defineEmits, defineExpose, ref, reactive, onMounted } from 'vue'

const props = defineProps(['oids']);
const emits = defineEmits(['changeSelect']);


let tableRef = ref();
let tableData = ref([]);

function changeSelection(selection) {
  emits('changeSelect', selection)
}

function handleTableRowSelectable(row, index) {
  return row.isValidate;
}

defineExpose({
  tableRef
})



onMounted(()=>{
  tableData.value = props.oids
})

</script>

<style scoped lang="scss"></style>