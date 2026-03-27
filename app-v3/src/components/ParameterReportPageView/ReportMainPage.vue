<template>
  <el-container>
    <el-main style="padding:0;overflow: inherit">
      <template v-if="redirectRoute[page.Key]">
        <component :is="redirectRoute[page.Key]" :page="page"/>
      </template>
      <template v-else>
        <el-row v-loading="loading" style="height=100%; padding-inline:16px; ">
          <el-col style="padding-right:16px;">
            <report-layout v-for="item in page.Items" :key="item.Key" :owner="page" :data="item" />
          </el-col>
        </el-row>
      </template>
    </el-main>
  </el-container>
</template>

<script>

import topo from "@/components/Topo/Topo.vue";
import systemStats from "@/components/SystemStats/Stats.vue";

import { computed } from 'vue'
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: 'ReportMainPage',
  components: {topo, systemStats},
  provide() {
    // 使用函数的形式，可以访问到 `this`
    return {
      viewMode: computed(() => this.mode)
    }
  },
  props: {
    page: Object,
    mode: {
      type: String,
      default: provideKeys.viewModeDefaultValue,
    }
  },
  data() {
    return {
      loading: false,
      redirectRoute: {
        'das_topo': 'topo',
        'stats': 'systemStats'
      },
    }
  },
  mounted() {
  },
  methods: {
  },
}
</script>
<style lang="scss" scoped></style>
