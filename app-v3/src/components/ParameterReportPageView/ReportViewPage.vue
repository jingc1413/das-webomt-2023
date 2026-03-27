<template>
  <el-main style="padding:0;">
    <el-row>
      <el-col style="padding-right:16px;">
        <report-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
      </el-col>
    </el-row>
  </el-main>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';


export default {
  name: 'ReportViewPage',
  components: {},
  props: {
    owner: Object,
    page: Object,
    title: String,
    closable: Boolean,
    handleClose: Function,
  },
  emits: ["close"],
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev
    };
  },
  data() {
    return {
      loading: false,
    }
  },
  methods: {
    closeViewPage() {
      this.dev.layout.closeViewPage({ page: this.page })
    },
  }
}
</script>
<style lang="scss" scoped></style>
