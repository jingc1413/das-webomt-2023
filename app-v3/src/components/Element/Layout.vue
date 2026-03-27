<!-- eslint-disable vue/html-indent -->
<!-- eslint-disable vue/html-closing-bracket-newline -->
<!-- eslint-disable vue/max-attributes-per-line -->
<!-- eslint-disable vue/first-attribute-linebreak -->
<template>
  <template v-if="data?.Type === 'Layout:Tabs'">
    <el-tabs v-if="provideKeys.viewModePrintValue != viewMode" :model="activeTab"
      :tab-position="data?.Style?.position || 'left'" style="height:100%" :class="'el-tabs-view-'+(data?.Style?.position || 'left')">
      <template v-for="item in data.Items ">
        <el-tab-pane v-if="item.Type === 'Page'" :key="item.Key" :label="item.Name">
          <el-scrollbar v-if="data?.Style?.position === 'top'"
            style="height: calc(100vh - var(--header-height) - var(--main-page-header-height) -  var(--view-page-header-height) - 10px);">
            <my-view-page :owner="owner" :page="item" style="height: calc(100vh - var(--header-height) - var(--main-page-header-height) -  var(--view-page-header-height) - 10px);"/>
          </el-scrollbar>
          <el-scrollbar v-else
            style="height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height));">
            <my-view-page :owner="owner" :page="item" style="height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height));"/>
          </el-scrollbar>
        </el-tab-pane>
      </template>
    </el-tabs>
    <template v-else-if="provideKeys.viewModePrintValue == viewMode">
      <my-view-page v-for="(item, index) in data.Items" :key="index" :owner="owner" :page="item" />
    </template>
  </template>
  <el-col v-if="data?.Type === 'Layout:Col'" :span="data?.Style?.span || 24" :class="{'flex-table':viewMode == provideKeys.viewModeDefaultValue}">
    <my-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-col>
  <el-row v-else-if="data?.Type === 'Layout:Row'" align="top">
    <my-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
  <el-row v-else>
    <my-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-row>
</template>

<script>
import provideKeys from '@/utils/provideKeys.js'
export default {
  name: 'MyLayout',
  // 通过 props 接收父组件传过来的值(可以直接在模板中使用)
  // props 也可以是对象接收
  props: {
    owner: Object,
    data: Object,
  },
  inject: ['viewMode'],
  data() {
    return {
      provideKeys,
      activeTab: undefined,
    }
  },
}
</script>
<style lang="scss" scoped>
.el-tabs-view-top :deep(.my-table-height) {
  height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--main-page-tabs-height) - var(--view-page-header-height) - 10px);
}

.el-tabs-view-top :deep(.my-form-height .el-scrollbar__wrap) {
  max-height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--main-page-tabs-height) - var(--view-page-header-height) - 36px) !important;
}
.el-tabs-view-left :deep(.my-table-height) {
  max-height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height) - 36px) !important;
}

.el-tabs-view-left :deep(.my-form-height .el-scrollbar__wrap) {
  max-height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height) - 36px) !important;
}
</style>
