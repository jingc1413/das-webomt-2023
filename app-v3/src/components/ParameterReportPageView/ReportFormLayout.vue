<template>
  <el-col v-if="data?.Type === 'Layout:Col'" :span="data?.Style?.span || 24">
    <report-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-col>
  <el-row v-else-if="data?.Type === 'Layout:Row'">
    <report-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
  <el-row v-else>
    <report-form-element :owner="owner" :data="data" />
  </el-row>
</template>

<script>

export default {
  name: 'ReportFormLayout',
  props: {
    owner: Object,
    data: Object,
  },
}
</script>

<style lang="scss" scoped></style>
