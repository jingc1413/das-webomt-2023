<template>
  <template v-if="page.Key === 'statistics'">
    <my-statistics-page />
  </template>
  <template v-else-if="page.Key === 'ping_diagnostic'">
    <my-ping-diag-page />
  </template>
  <template v-else-if="page.Key === 'spectrum_diagnostic'">
    <span>TBD</span>
  </template>
  <template v-else>
    <el-main style="padding:0;">
      <el-icon v-if="closable" style="position: absolute; right: 8px; top: 16px; z-index: 100; cursor: pointer;">
        <Close @click="$emit('close')" />
      </el-icon>
      <el-row style="height:100%">
        <template v-if="viewMode == provideKeys.viewModePrintValue">
          <el-col :span="24" style="height: 100%;padding-right:16px;">
            <my-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
          </el-col>
        </template>
        <template v-else>
          <el-col :span="page.viewPage ? 14 : 24" style="height: 100%;padding-right:16px;">
            <my-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
          </el-col>
          <el-col v-if="page.viewPage" :span="10"
            style="padding-inline:16px; border-left:solid 1px var(--el-menu-border-color);">
            <my-view-page :owner="page" :page="page.viewPage" :title="page.viewPage.Name" :closable="true"
              @close="closeViewPage" />
          </el-col>
        </template>
      </el-row>
    </el-main>
  </template>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: 'MyViewPage',
  inject: ['viewMode'],
  props: {
    owner: { type: Object, default: undefined, },
    page: { type: Object, default: undefined, },
    title: { type: String, default: "", },
    closable: { type: Boolean, default: false, },
    handleClose: { type: Function, default: undefined, },
  },
  emits: ["close"],
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
      provideKeys,
    };
  },
  data() {
    return {
      loading: false,
      supportRefresh: false,
      toolbarVisible: true,
    }
  },
  methods: {
    getParameterValues: async function () {
      if (!this.supportRefresh) return;
      if (this.page?.oids?.length > 0) {
        this.loading = true;
        await this.dev.params.getParameterValues({
          oids: this.page.oids,
          showMessage: true,
        })
        this.loading = false;
      }
    },
    closeViewPage() {
      this.dev.layout.closeViewPage({ page: this.page })
    },
  }
}
</script>
<style lang="scss" scoped></style>
