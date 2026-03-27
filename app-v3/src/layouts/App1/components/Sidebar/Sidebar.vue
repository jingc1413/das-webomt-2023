<!-- eslint-disable vue/no-v-for-template-key -->
<template>
  <div :class="{ 'has-logo': showLogo }">
    <logo v-if="showLogo" :collapse="isCollapse" />
    <el-scrollbar wrap-class="scrollbar-wrapper" style="height: calc(100vh - var(--header-height) - 40px);">
      <el-menu :default-active="activeMenu" :collapse="isCollapse" :background-color="variables.menuBg"
        :text-color="variables.menuText" :unique-opened="true" :active-text-color="variables.menuActiveText"
        :collapse-transition="false" mode="vertical">
        <template v-for="module in appLayout.Items">
          <el-sub-menu v-if="module.superModeOnly === false || (module.superModeOnly && !auth.superModeDisabled)"
            :key="module.Name" :index="module.Key">
            <template #title>
              <el-icon>
                <component v-if="module.Style?.icon" :is="module.Style?.icon" />
                <Menu v-else />
              </el-icon>
              <span>
                {{ module.Name }}
              </span>
            </template>
            <el-menu-item-group>
              <el-menu-item v-for="page in module.Items" :key="page.Name"
                :index="`/${dasDevices.currentDeviceSub}/${module.Key}/${page.Key}`"
                @click="selectModulePage(module.Key, page.Key)">
                {{ page.Name }}
              </el-menu-item>
            </el-menu-item-group>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-scrollbar>
    <div style="position: absolute;bottom: 10px;width: 100%;margin-left: calc(40% - 20px);">
      <hamburger :is-active="appStore.sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />
    </div>
  </div>
</template>

<script>
import Logo from "./Logo.vue";
import Hamburger from '@/components/Hamburger/Hamburger.vue'
import variables from "@/assets/styles/variables.module.scss";
import { useSettingsStore } from "@/stores/settings";
import { useAppStore } from "@/stores/app";
import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import { Cpu, Tools, OfficeBuilding, Monitor, Document, Menu, Operation } from "@element-plus/icons-vue";
import { useRouter } from "vue-router";
import { getPageVersion } from '@/utils/index.js';

export default {
  components: { Logo, Hamburger },
  setup() {
    const auth = useAuthStore();
    const appStore = useAppStore();
    const settingsStores = useSettingsStore();
    const dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    const router = useRouter();
    function selectModulePage(moduleKey, pageKey) {
      router.push({
        name: "app1Page",
        params: {
          module: moduleKey,
          page: pageKey,
        },
      });
    }
    return {
      settingsStores,
      auth,
      appStore,
      selectModulePage,
      dasDevices,
      dev,
    };
  },
  data() {
    return {
      routes: this.$router.options.routes,
    };
  },
  computed: {
    appLayout() {
      return this.dasDevices?.currentDevice?.layout?.appLayout;
    },
    activeMenu() {
      const route = this.$route;
      const { meta, path } = route;
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu;
      }
      let version = getPageVersion();
      return path.replace(`/${version}`, '');
    },
    showLogo() {
      return this.settingsStores.sidebarLogo;
    },
    variables() {
      return variables;
    },
    isCollapse() {
      return !this.appStore.sidebar.opened;
    },
  },
  watch: {
    app() {
      let newRoute = this.$router.getRoutes().filter((item) => {
        if (item.path.search(/\/.*\/.*/g) != -1) {
          return false;
        }        // }
        return true;
      });
      this.routes = newRoute;
    },
    appLayout() {
      if (this.$route.params?.module && this.$route.params?.page && (this.$route.params.sub == this.dasDevices.currentDeviceSub)) {
        this.selectModulePage(this.$route.params.module, this.$route.params.page);
      }
    },
    activeMenu() {
      console.log({activeMenu:this.activeMenu});
    }
  },
  mounted() {
    if (this.$route.params?.module && this.$route.params?.page) {
      this.selectModulePage(this.$route.params.module, this.$route.params.page);
    }
  },
  methods: {
    toggleSideBar() {
      this.appStore.toggleSideBar()
    },
  }
};
</script>
