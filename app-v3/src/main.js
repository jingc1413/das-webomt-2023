// Object.hasOwn 低版本浏览器兼容
if (!Object.hasOwn) {
  Object.defineProperty(Object.prototype, "hasOwn", {
    value: (obj, key) => {
      return Object.prototype.hasOwnProperty.call(obj, key);
    },
    enumerable: false,
  });
}

import settings from "@/settings.js";
console.info(
  "version: " + settings.appInfo.version + ", build: " + settings.appInfo.build
);

import "element-plus/theme-chalk/el-message-box.css";
import "element-plus/theme-chalk/el-message.css";
import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "@/router";
import App from "./App.vue";
import ElementPlus, { dayjs } from "element-plus";
import "element-plus/dist/index.css";
import "@/assets/styles/index.scss"; // global css
import "@/icons"; // icon
import SvgIcon from "@/components/SvgIcon"; // svg component

import importComponentView from "./components/importComponentView.js";

import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import { useAuthStore } from "./stores/auth";
import ImportReportViewComponents from "./components/ParameterReportPageView/importReportView";

let utc = require("dayjs/plugin/utc");
dayjs.extend(utc);

let { nodeTest, deviceInfo } = settings;
const app = createApp(App);

// i18n 国际化
import i18n from "@/i18n";
import { useAppStore } from "./stores/app";
app.use(i18n);
//

const store = createPinia();

app.component("svg-icon", SvgIcon);
importComponentView(app);
ImportReportViewComponents(app);

app.use(store);
await useAuthStore().setup();
await useAppStore().setup();

app.use(router);
app.use(ElementPlus, { size: "small" });

app.mount("#app");

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// plugins
import plugins from './plugins' 
app.use(plugins);
// directive
import directive from './directive'
directive(app)
