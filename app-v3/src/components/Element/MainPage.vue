<template>
  <el-container>
    <el-main style="padding:0;"
      :style="{ 'overflow': mode == provideKeys.viewModePrintValue ? 'inherit' : 'auto', 'height': mode == provideKeys.viewModeDefaultValue ? 'calc(100vh - var(--header-height) - 10px)' : 'auto' }">
      <template v-if="page.Key === 'das_topo'">
        <my-topo :page="page" @clicked-device-event="callback" />
      </template>
      <template v-else-if="page.Key === 'statistics'">
        <my-statistics-page />
      </template>
      <template v-else-if="page.Key === 'ping_diagnostic'">
        <my-ping-diag-page />
      </template>
      <template v-else-if="page.Key === 'spectrum_diagnostic'">
        <span>TBD</span>
      </template>
      <template v-else>
        <el-row v-loading="loading" style="height:100%; padding-inline:16px; ">
          <template v-if="mode == provideKeys.viewModePrintValue">
            <el-col style="padding-right:16px;">
              <my-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
            </el-col>
          </template>
          <template v-else>
            <el-col :span="page.viewPage ? 14 : 24" style="height:100%;padding-right:16px;">
              <my-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
            </el-col>
            <el-col v-if="page.viewPage" :span="10"
              style="padding-inline:16px; border-left:solid 1px var(--el-menu-border-color);">
              <my-view-page :owner="page" :page="page.viewPage" :title="page.viewPage.Name" :closable="true"
                @close="closeViewPage" />
            </el-col>
          </template>
        </el-row>
      </template>
    </el-main>
  </el-container>
</template>

<script>
import { useDasDevices } from "@/stores/das-devices";

import { computed } from 'vue'
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: 'MyMainPage',
  components: {},
  provide() {
    // 使用函数的形式，可以访问到 `this`
    return {
      viewMode: computed(() => this.mode)
    }
  },
  props: {
    mode: {
      type: String,
      default: provideKeys.viewModeDefaultValue,
    },
    propPage: {
      type: Object,
      default: undefined,
    }
  },
  setup(prop) {
    const dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    let page = dev.layout.currentPage;
    if (prop.propPage) {
      page = prop.propPage
    }
    return {
      dev,
      dasDevices,
      provideKeys,
      page,
    };
  },
  data() {
    return {
      loading: false,
    }
  },
  methods: {
    callback(data) {
      console.log('maindata', data)
    },
    closeViewPage() {
      this.dev.layout.closeViewPage({ page: this.page })
    },
  },
}
</script>
<style lang="scss" scoped></style>
