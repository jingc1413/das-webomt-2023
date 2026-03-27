<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
        <span v-if="item.redirect === 'noRedirect' || index == levelList.length - 1" class="no-redirect max-width-title"
          :title="item.meta.title">{{ item.meta.title
          }}</span>
        <span v-else-if="index > 0" class="max-width-title" :title="item.meta.title">{{ item.meta.title }}</span>
        <span v-else class="max-width-title" :title="item.meta.title">{{ item.meta.title}}</span>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices'
import { match, pathToRegexp } from 'path-to-regexp'

export default {
  setup() {
    const dasDevices = useDasDevices();
    return {
      dasDevices,
    }
  },
  data() {
    return {
      levelList: null
    }
  },
  watch: {
    $route() {
      this.getBreadcrumb();
    },
    'dasDevices.currentDevice'()  {
      this.getBreadcrumb();
    }
  },
  created() {
    this.getBreadcrumb()
  },
  methods: {
    getBreadcrumb() {
      let { params, name } = this.$route;
      let matched = this.$route.matched.filter(item => item.meta && item.meta.title)
      matched = matched.filter(item => item.meta && item.meta.title && item.meta.breadcrumb !== false)
      matched[0] = {
        meta: {
          title: this.dasDevices.currentDevice.layout.currentDeviceNameTitle,
        },
        path: this.dasDevices.currentDeviceSub,
      }
      if (name.endsWith('Page') == true) {
        let replaceLastItem = this.filterPageTitle(params);
        matched.splice(-1, 1, ...replaceLastItem);
      }
      this.levelList = matched;
    },
    pathCompile(path) {
      // To solve this problem https://github.com/PanJiaChen/vue-element-admin/issues/561
      const { params } = this.$route
      var toPath = pathToRegexp.compile(path)
      return toPath(params)
    },
    handleLink(item) {
      const { redirect, path } = item
      if (redirect) {
        this.$router.push(redirect)
        return
      }
      this.$router.push(this.pathCompile(path))
    },
    filterPageTitle(params) {
      let { page, module } = params;
      let tempModule = this.dasDevices.currentDevice.layout.getModule(module);
      let tempPage = this.dasDevices.currentDevice.layout.getPage(module, page);
      if (!tempModule) {
        tempModule = { Name: this.toUpperCaseFirst(module.split('_')), Key: 'tempModule' };
      }
      if (!tempPage) {
        tempPage = { Name: this.toUpperCaseFirst(page.split('_')), Key: 'tempPage' };
      }
      return [{ meta: { title: tempModule.Name }, path: tempModule.Key }, { meta: { title: tempPage.Name }, path: tempPage.Key }]
    },
    toUpperCaseFirst(charts = []) {
      return charts.map(item => {
        return item.charAt(0).toUpperCase() + item.slice(1)
      }).join(' ')
    }
  }
}
</script>

<style lang="scss" scoped>
.app-breadcrumb.el-breadcrumb {
  display: inline-block;
  font-size: 14px;
  line-height: var(--header-height);
  margin-left: 8px;

  .no-redirect {
    color: #97a8be;
    cursor: text;
  }

  .max-width-title {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
    max-width: 20vw;
    display: block;
  }
}
</style>
