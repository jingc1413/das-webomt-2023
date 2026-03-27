<template>
  <div id="app">
    <el-config-provider :locale="langList[locale]">
      <RouterView />
      <my-invalid-view />
    </el-config-provider>
  </div>
  <my-confirm-dialog />
  <my-view-file-dialog />
</template>
<script>

import en from "element-plus/dist/locale/en.min.mjs";
import { useI18n } from 'vue-i18n';

export default {
  setup() {
    const { locale } = useI18n()
    return { locale };
  },
  data() {
    return {
      count: 0,
      langList: {
        en
      }
    };
  },
  mounted() {
  },
  methods: {
  }
};
</script>
