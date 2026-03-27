<template>
  <el-col v-if="data?.Type === 'Layout:Col'" :span="data?.Style?.span || 24">
    <my-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-col>
  <el-row v-else-if="data?.Type === 'Layout:Row'">
    <my-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
  <el-row v-else>
    <my-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
</template>

<script>

export default {
  name: 'MyFormLayout',
  props: {
    owner: Object,
    data: Object,
  },
}
</script>

<style lang="scss" scoped></style>
