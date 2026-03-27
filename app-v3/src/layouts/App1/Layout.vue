<template>
  <div :class="classObj" class="app-wrapper">
    <div v-if="device==='mobile'&&sidebar.opened" class="drawer-bg" @click="handleClickOutside" />
    
    <div class="main-container">
      <div :class="{'fixed-header':fixedHeader}">
        <navbar />
      </div>
      <div style="display: flex;flex-direction: row;flex-wrap: nowrap;">
        <el-scrollbar style="height:calc(100vh - var(--header-height));">
          <sidebar class="sidebar-container" />
        </el-scrollbar>
        
        <el-scrollbar style="height:calc(100vh - var(--header-height)); width: 100%;">
          <app-main />
        </el-scrollbar>
      </div>
    </div>
  </div>
</template>

<script>
import { Navbar, Sidebar, AppMain } from './components'
// import ResizeMixin from './mixin/ResizeHandler'
import { mapState, mapActions } from 'pinia'
import {useAppStore} from '@/stores/app.js'
import {useSettingsStore} from '@/stores/settings.js'
import { useDocumentVisibility } from '@vueuse/core'
import { useAuthStore } from "@/stores/auth"

export default {
  components: {
    Navbar,
    Sidebar,
    AppMain
  },
  mixins: [],
  setup() {
    const authStore = useAuthStore()
    const documentVisibility = useDocumentVisibility();
    return {
      authStore,
      documentVisibility
    }
  },
  computed: {
    ...mapState(useAppStore, ['sidebar','device','fixedHeader']),
    // ...mapState(useAppStore, ['device']),
    // ...mapState(useSettingsStore, ['fixedHeader']),
    classObj() {
      return {
        hideSidebar: !this.sidebar.opened,
        openSidebar: this.sidebar.opened,
        withoutAnimation: this.sidebar.withoutAnimation,
        mobile: this.device === 'mobile'
      }
    }
  },
  watch: {
    documentVisibility(value) {
      if (value == 'visible') {
        this.authStore.checkLogin();
      }
    }
  },
  methods: {
    ...mapActions(useAppStore, ['closeSideBar']),
    handleClickOutside() {
      this.closeSideBar(true);
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "@/assets/styles/mixin.scss";
  @import "@/assets/styles/variables.module.scss";

  .app-wrapper {
    @include clearfix;
    position: relative;
    height: 100%;
    width: 100%;
    &.mobile.openSidebar{
      position: fixed;
      top: 0;
    }
  }
  .drawer-bg {
    background: #000;
    opacity: 0.3;
    width: 100%;
    top: 0;
    height: 100%;
    position: absolute;
    z-index: 999;
  }

  .fixed-header {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 9;
    width: calc(100% - #{$sideBarWidth});
    transition: width 0.28s;
  }

  .hideSidebar .fixed-header {
    width: calc(100% - 54px)
  }

  .mobile .fixed-header {
    width: 100%;
  }
</style>
