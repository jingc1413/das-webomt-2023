import Cookies from "js-cookie";
import { defineStore } from "pinia";
import { useDasDevices } from "./das-devices";
import router from "@/router";
import apix from "@/api";
import settings from "@/settings";

export const useAppStore = defineStore("app", {
  state: () => ({
    /** @type { opened: boolean, withoutAnimation: boolean } */
    sidebar: {
      opened: Cookies.get("sidebarStatus") ? !!+Cookies.get("sidebarStatus") : false,
      withoutAnimation: false,
    },
    /** @type boolean  */
    device: "desktop",
    _isAppSetupDone: false,
    appInfo: {},

    confirmDialog: {
      visible: false,
      title: "Confirm",
      content: undefined,
      callback: undefined,
      needInput: false,
      formModel: {
        inputValue: null,
      },
      inputRule: [],
    },
    viewFileDialog: {
      visible: false,
      title: "File View",
      supportCancel: true,
      handleCancel: undefined,
      supportSave: false,
      handleSave: undefined,
      handleLoadingData: undefined,
    },

    debugTooltipDisabled: false,
  }),
  getters: {
    isAppSetupDone(state) {
      return state._isAppSetupDone;
    },
    loginTitle(state) {
      return state.appInfo?.DeviceTypeName ?? "OMT";
    },
    appTitle(state) {
      if (state.appInfo?.Schema == "corning") {
        return "CORNING EVERON™ 6200"
      }
      return undefined;
    }
  },
  actions: {
    setup: async function () {
      await apix
        .getAppInfo()
        .then((info) => {
          this.appInfo = info;
        })
        .catch((e) => {
          console.log(e);
          this.appInfo = {
            Schema: "default",
            DeviceTypeName: "Primary A3",
          };
        });
    },
    toggleSideBar() {
      this.sidebar.opened = !this.sidebar.opened;
      this.sidebar.withoutAnimation = false;
      if (this.sidebar.opened) {
        Cookies.set("sidebarStatus", 1);
      } else {
        Cookies.set("sidebarStatus", 0);
      }
    },
    closeSideBar({ withoutAnimation }) {
      Cookies.set("sidebarStatus", 0);
      this.sidebar.opened = false;
      this.sidebar.withoutAnimation = withoutAnimation;
    },
    toggleDevice(device) {
      this.device = device;
    },
    setAppSetupDone() {
      this._isAppSetupDone = true;
    },
    clearAppSetupDone() {
      this._isAppSetupDone = false;
    },
    gotoLoginRoute: async function (redirect = undefined) {
      const currentRoute = router.currentRoute.value;
      if (currentRoute.path === "/login") {
        return;
      }
      if (redirect === undefined) {
        redirect = currentRoute?.fullPath;
      }
      if (redirect) {
        router.push(`/login?redirect=${redirect}`);
      } else {
        router.push(`/login`);
      }
    },
    gotoUpgradeRoute: async function (sub, redirect = undefined) {
      const currentRoute = router.currentRoute.value;
      if (currentRoute.path === `/${sub}/upgrading`) {
        return;
      }
      const dasDevices = useDasDevices();
      if (sub === undefined) {
        sub = dasDevices.currentDeviceSub || "local";
      }
      if (redirect === undefined) {
        redirect = currentRoute?.fullPath;
      }
      if (redirect) {
        router.push(`/${sub}/upgrading?redirect=${redirect}`);
      } else {
        router.push(`/${sub}/upgrading`);
      }
    },
    gotoDefaultRoute: async function (path = undefined) {
      const currentRoute = router.currentRoute.value;
      if (currentRoute.path === "/loading") {
        return;
      }
      const redirect = currentRoute?.query?.redirect;
      if (path) {
        router.push(path);
      } else if (redirect) {
        router.push(redirect);
      } else {
        router.push(`/loading`);
      }
      return;
    },
    getLoginBackground: function (Schema) {
      let fileName = Schema ? Schema.toLowerCase() : null;
      return apix.getCurrentBackground(fileName).catch((e) => {
        return apix.getCurrentBackground().catch((e) => {});
      });
    },
    openConfirmDialog({ title, content, callback, supportCancel = true, needInput = false, inputRule = [] }) {
      this.confirmDialog.visible = true;
      this.confirmDialog.title = title;
      this.confirmDialog.content = content;
      this.confirmDialog.supportCancel = supportCancel;
      this.confirmDialog.callback = callback;
      this.confirmDialog.needInput = needInput;
      this.confirmDialog.formModel.inputValue = "";
      this.confirmDialog.inputRule = inputRule;
    },
    closeConfirmDialog(ok = false) {
      this.confirmDialog.visible = false;
      if (this.confirmDialog.callback) {
        this.confirmDialog.callback(ok, this.confirmDialog.inputModel);
      }
    },
    openViewFileDialog({ title, supportCancel = true, handleCancel = ()=>{return null}, supportSave=false, handleSave, handleLoadingData=()=>{return Promise.resolve(null)} }) {
      this.viewFileDialog.title = title;
      this.viewFileDialog.supportCancel = supportCancel;
      this.viewFileDialog.handleCancel = handleCancel;
      this.viewFileDialog.supportSave = supportSave;
      this.viewFileDialog.handleSave = handleSave;
      this.viewFileDialog.handleLoadingData = handleLoadingData;
      this.viewFileDialog.visible = true;
    },
    closeViewFileDialog(ok = false) {
      this.viewFileDialog.visible = false;
      if (this.viewFileDialog.handleCancel) {
        this.viewFileDialog.handleCancel(ok);
      }
    },
  },
});
