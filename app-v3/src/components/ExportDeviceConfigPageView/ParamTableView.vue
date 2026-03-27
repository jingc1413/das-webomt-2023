<template>
  <el-table ref="tableRef" :data="tableList" table-layout="auto" style="padding-bottom: 16px;" @selection-change="changeSelection" row-key="oid">
    <el-table-column type="selection" width="55" />
    <el-table-column type="index" width="55" />
    <el-table-column label="Name" prop="">
      <template #default="scoped">
        <report-param :data="{OID:scoped.row.oid}" pageType="name"/>
      </template>
    </el-table-column>
    <el-table-column label="Value" prop="">
      <template #default="scoped">
        <report-param :data="{OID:scoped.row.oid, Value: scoped.row.value}" />
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup>
import { defineProps, defineEmits, defineExpose, ref, reactive, onMounted } from 'vue'

const props = defineProps(['oids', 'defaultValues']);
const emits = defineEmits(['changeSelect']);

let tableRef = ref(null);
let tableList = ref([]);

function changeSelection(selection) {
  emits('changeSelect', selection)
}

defineExpose({
  tableRef
})

onMounted(()=>{
  let defaultValues = {}
  props.defaultValues.forEach(item=>{
    defaultValues[item.oid]=item['value']
  })
  tableList.value = props.oids.map(item=>{
    return {
      oid: item,
      value: defaultValues[item]??undefined
    }
  })
})

</script>

<style scoped lang="scss"></style>