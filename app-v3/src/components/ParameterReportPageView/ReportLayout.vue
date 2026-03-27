<!-- eslint-disable vue/html-indent -->
<!-- eslint-disable vue/html-closing-bracket-newline -->
<!-- eslint-disable vue/max-attributes-per-line -->
<!-- eslint-disable vue/first-attribute-linebreak -->
<template>
  <template v-if="data?.Type === 'Layout:Tabs'">
    <report-view-page v-for="(item, index) in data.Items" :key="index" :owner="owner" :page="item" />
  </template>
  <el-col v-if="data?.Type === 'Layout:Col'" :span="data?.Style?.span || 24" class="flex-table" style="height: auto;">
    <report-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-col>
  <el-row v-else-if="data?.Type === 'Layout:Row'">
    <report-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
  <el-row v-else>
    <report-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
</template>

<script>
export default {
  name: 'ReportLayout',
  // 通过 props 接收父组件传过来的值(可以直接在模板中使用)
  // props 也可以是对象接收
  props: {
    owner: Object,
    data: Object,
  },
  data() {
    return {
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
