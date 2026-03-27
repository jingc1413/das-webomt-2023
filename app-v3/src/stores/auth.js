
import apix from "@/api"
import { clearSession } from "@/utils/auth";
import { defineStore } from "pinia";
import { useAppStore } from "./app";
import settings from "@/settings";
import Cookies from 'js-cookie'
import { ElMessage } from "element-plus";
import { translator as t } from '@/i18n';

export const useAuthStore = defineStore("user", {
  state: () => ({
    name: "",
    avatar: "",
    currentUser: undefined,
    superModeDisabled: settings.nodeTest ? false : Cookies.get('superModeDisabled') ? !!+Cookies.get('superModeDisabled') : true,
    superTestDisabled: settings.nodeTest ? false : Cookies.get('superTestDisabled') ? !!+Cookies.get('superTestDisabled') : true,
    permissions: [],
    extensionsPermissions: [],
  }),
  getters: {
    loginUserName(state) {
      return state.currentUser?.Name || "";
    },
    loginUserRoles(state) {
      return state.currentUser?.Roles || [];
    },
  },
  actions: {
    setup: async function () {
      await this.checkLogin();
      if (settings.nodeTest) {
        this.getCurrentUser();
      }
    },
    login: async function (username, password, redirectPath = undefined) {
      return apix.login(username, password).catch(e => {
        ElMessage.error(e.message);
      });
    },
    logout: async function (message = undefined) {
      if (this.currentUser !== undefined) {
        console.log("logout");
        const appStore = useAppStore();

        await apix.logout().catch(e => {
          ElMessage.error(e.message);
        });
        this.clearPermissions()
        this.currentUser = undefined;
        clearSession();
        appStore.clearAppSetupDone();
        await appStore.gotoLoginRoute();
        if (message) {
          ElMessage.warning(message);
        }
      }
      return Promise.resolve(true);
    },
    isLogin: async function () {
      if (settings.nodeTest) {
        this.currentUser = { Name: 'admin' };
        return true;
      }
      const isLogin = this.currentUser != undefined;
      console.log("isLogin", isLogin)
      return isLogin;
    },
    checkLogin: async function (showMessage = false) {
      await this.getCurrentUser();
      const isLogin = this.currentUser != undefined;
      if (!isLogin) {
        clearSession();
        if (showMessage) {
          ElMessage.error(t("tip.NotLogin"));
        }
      }
    },
    getCurrentUser: function () {
      return new Promise((resolve, reject) => {
        apix.getCurrentUser().then(user => {
          this.currentUser = user;
          this.clearPermissions()
          this.setupAllows(user.Rules)
          resolve(this.currentUser);
        }).catch(e => {
          this.clearPermissions()
          resolve(null)
          console.log(e);
        }).finally(()=>{
        })
      })
    },
    setupAllows: function (rules = []) {
      this.permissions = rules;
      let componentRules = new Set();
      for (let ruleItem of rules) {
        if (ruleItem.includes('.*')) {
          this.extensionsPermissions.push(ruleItem);
        } else {
          let ruleArray = ruleItem.split(".");
          if (ruleArray[0] == 'page') {
            componentRules.add(`page.${ruleArray[1]}`)
            if (ruleArray[2] != 'get' || ruleArray[2] != 'set') {
              componentRules.add(`page.${ruleArray[1]}.${ruleArray[2]}`)
            }
            if (ruleArray[3] != 'get' || ruleArray[3] != 'set') {
              componentRules.add(`page.${ruleArray[1]}.${ruleArray[2]}.${ruleArray[3]}`)
            }
          }
        }
        
      }
      this.permissions.push(...componentRules);
    },
    setSuperModeEnable() {
      this.superModeDisabled = false;
      Cookies.set('superModeDisabled', 0);
      ElMessage.success("Super mode entered");
    },
    setSuperModeDisable() {
      this.superModeDisabled = true;
      Cookies.set('superModeDisabled', 1);
      ElMessage.success("Super mode exited");
    },
    setSuperTestEnable() {
      this.superTestDisabled = false;
      Cookies.set('superTestDisabled', 0);
      ElMessage.success("Test mode entered");
    },
    setSuperTestDisable() {
      this.superTestDisabled = true;
      Cookies.set('superTestDisabled', 1);
      ElMessage.success("Test mode exited");
    },
    clearPermissions() {
      this.permissions = [];
      this.extensionsPermissions = [];
    },
  },
});
