import { defineStore } from "pinia";
import { useDasModel } from "./das-model";
import apix from "@/api";
import { ElMessage, dayjs } from "element-plus";
import { translator as t } from "@/i18n";
import settings from "@/settings";
import { useDasTopo } from "./topo";
import { newDevice } from "./device";
import { setupDeviceInfos } from "./model/infos";

export const useDasDevices = defineStore("dasDevices", {
  state: () => ({
    localInfo: {
      SubID: 0,
      DeviceTypeName: "Primary A3",
      ProductModel: "",
      SiteName: "null",
      DeviceName: "null",
      ElementSerialNumber: "null",
    },
    currentSub: undefined,
    currentDevice: undefined,
    allDevices: [],
  }),
  getters: {
    getDevice(state) {
      return function (sub) {
        if (sub === "local") {
          return state.localDevice;
        }
        return state.allDevices?.find((dev) => dev.sub == sub);
      };
    },

    getDeviceInfo(state) {
      return function (sub) {
        try {
          if (sub === "local") {
            return state.localInfo;
          }
          return state.deviceInfos.getInfo(sub);
        } catch (e) {
          return undefined;
        }
      };
    },
    getDeviceInfos(state) {
      return state.deviceInfos?.infos || [];
    },
    isLocalDevice(state) {
      return state.currentDeviceSub === "local";
    },
    isCurrentDevice(state) {
      return function (sub) {
        if (sub === "local") {
          return state.currentSub === sub;
        }
        const sub2 = state.currentDeviceInfo?.SubID;
        return String(sub2) === String(sub);
      };
    },
    currentDeviceSub(state) {
      if (state.currentSub == undefined) {
        return "local";
      }
      // if (String(state.currentSub) == String(state.localInfo.SubID)) {
      //   return "local";
      // }
      return String(state.currentSub);
    },
    currentDeviceInfo(state) {
      try {
        if (state.currentSub === "local") {
          return state.localInfo;
        }
        return state.deviceInfos.getInfo(state.currentSub);
      } catch (e) {
        return undefined;
      }
    },
    currentProductType(state) {
      return state.currentDeviceInfo?.ProductType || undefined;
    },
    currentDeviceFileName(state) {
      return function (name, withSN = true) {
        let title = state?.currentDeviceInfo?.ProductModel ?? "null";
        title = title + "_" + name;
        if (withSN) {
          title = title + "_" + state?.currentDeviceInfo?.ElementSerialNumber ?? "null";
        }
        title = title + "_" + dayjs().format("YYYYMMDD");
        return title;
      };
    },
    currentDevice(state) {
      return state.getDevice(state.currentDeviceSub);
    },
  },
  actions: {
    setup: async function () {
      const self = this;
      this.allDevices = [];
      
      await apix
        .getLocalInfo()
        .then((info) => {
          this.localInfo = info;
        })
        .catch((e) => {
          console.log(e);
        });

      this.localDevice = await this.setupDevice("local", this.localInfo);
      await this.localDevice.checkAvailable();

      this.deviceInfos = setupDeviceInfos({
        onChangeInfos: function (infos) {
          setTimeout(() => {
            useDasTopo().setupTopoData(infos);
          }, 200);
        },
        onRemoveInfo: function (sub) {
          self.removeDeviceLayoutAndParams(sub);
        },
        onUpdateInfo: function (sub) {
          self.setupDevice(sub);
        },
      });
      await this.updateDeviceInfos();
      // this.selectDeviceSub('local');
      await this.openWebSocket();
      console.log("setup done", this.localInfo);
    },
    updateDeviceInfos: async function (force = false, showMessage = false) {
      await this.deviceInfos.updateInfos(force, showMessage);
    },
    updateDeviceInfo: async function (sub, showMessage = false) {
      await this.deviceInfos.updateInfo(sub, showMessage);
      return this.getDeviceInfo(sub);
    },
    updateCurrentDeviceInfo: async function (showMessage = false) {
      return await this.updateDeviceInfo(this.currentDeviceSub, showMessage);
    },
    closeWebSocket: async function () {
      if (this.ws) {
        this.ws.close();
      }
    },
    openWebSocket: async function () {
      await this.closeWebSocket();
      const self = this;
      this.ws = await apix.connectToWebSocket({
        onOpen: function (evt) {
          console.log("onopen", { evt });
        },
        onClose: (evt) => {
          console.log("onclose", { evt });
        },
        onError: (evt) => {
          console.log("onerror", { evt });
          self.deviceInfos.updateQueryStatus({
            Loading: false,
            Finished: true,
            Success: false,
            Message: "WebSocket error",
          });
        },
        onMessage: (evt) => {
          // console.log('onmessage', { evt });
          const msg = JSON.parse(evt.data);
          switch (msg.Type) {
            case "queryDeviceData": {
              self.deviceInfos.setInfo(msg.Data);
              break;
            }
            case "queryDeviceStatus": {
              self.deviceInfos.updateQueryStatus(msg.Data);
              break;
            }
          }
        },
      });
    },
    setupDevice: async function (sub, info = undefined) {
      try {
        const dasModel = useDasModel();
        if (info == undefined) {
          info = this.getDeviceInfo(sub);
        }
        if (info == undefined) {
          console.error("setup device error, cant get device info", sub);
          return undefined;
        }
        info.invalidDeviceTypeName = !dasModel.validateDeviceTypeName(info.DeviceTypeName);
        if (info.invalidDeviceTypeName) {
          console.error("setup device error, invalid device type", sub, info);
          return undefined;
        }
        let dev = this.getDevice(sub);
        if (dev == undefined || dev.info.DeviceTypeName !== info.DeviceTypeName || dev.info.Version !== info.Version) {
          dev = await newDevice(sub, info);
          if (dev) {
            this.setDevice(dev);
          }
        }
        return dev;
      } catch (e) {
        console.error("setup device error", sub);
      }
    },
    removeDevice(sub) {
      this.allDevices = this.allDevices.filter((v) => String(v.sub) != String(sub));
    },
    setDevice(dev) {
      this.removeDevice(dev.sub);
      this.allDevices.push(dev);
    },

    getDeviceParameterValues: async function ({
      sub,
      oids,
      values = [],
      showMessage = false,
      statusEnable = false,
      options = {},
    }) {
      if (oids.includes("T02.P0006") && !oids.includes("TB0.P0AFF")) {
        oids.push("TB0.P0AFF");
      }
      const dev = this.getDevice(sub);
      const filterValues = [];
      const filterOids = oids.filter((oid) => {
        const param = dev.params.getParam(oid);
        if (param && param.Readable) {
          const v = values.find((v) => v.oid === oid);
          if (v != undefined) {
            v.value = param.formatRawFromValue(v.value);
            filterValues.push(v);
          }
          return true;
        }
        return false;
      });
      if (settings.nodeTest) {
        console.log("get values", sub, filterOids, filterValues);
        const returnValues = filterOids.map((oid) => {
          const param = dev.params.getParam(oid);
          return { oid: oid, value: param.Value, code: "00" };
        });
        return returnValues;
      }
      console.log("get values", sub, filterOids, filterValues);
      return new Promise((resolve, reject) => {
        apix
          .getParameterValues({
            sub,
            params: dev.params,
            oids: filterOids,
            values: filterValues,
          })
          .then((returnValues) => {
            const updateValues = returnValues
              .map((v) => {
                const param = dev.params.getParam(v.oid);
                v.value = param?.setRawValue(v.value, v.code);
                return v;
              })
              .filter((v) => v.value !== undefined);
            if (showMessage) {
              ElMessage.success(t("tip.GetParameterValuesSuccessfully"));
            }
            return resolve(updateValues);
          })
          .catch((e) => {
            console.error(e);
            if (showMessage) {
              ElMessage.error(t("tip.RequestFailed"));
            }
            return resolve([]);
          });
      });
    },
    setDeviceParameterValues: async function ({
      sub,
      values,
      showMessage = false,
      statusEnable = false,
      options = {},
    }) {
      const dev = this.getDevice(sub);
      const filterValues = values.filter((v) => {
        const param = dev.params.getParam(v.oid);
        if (param && param.Writable && !param.InputDisabled && v.value !== undefined) {
          v.value = param.formatRawFromValue(v.value);
          return true;
        }
        return false;
      });
      if (settings.nodeTest) {
        console.log("set values", sub, filterValues);
        const returnValues = filterValues.map((v) => {
          const param = dev.params.getParam(v.oid);
          param.setRawValue(v.value);
          return { oid: v.oid, value: param.Value, code: "00" };
        });
        return returnValues;
      }
      console.log("set values", sub, filterValues);
      return new Promise((resolve, reject) => {
        apix
          .setParameterValues({
            sub,
            params: dev.params,
            values: filterValues,
          })
          .then((returnValues) => {
            const filterValues = returnValues.filter((v) => v.value !== undefined);
            const updateValues = filterValues.map((v) => {
              const param = dev.params.getParam(v.oid);
              v.value = param?.setRawValue(v.value, v.code);
              return v;
            });
            if (showMessage) {
              const faultValues = updateValues.filter((v) => v.code !== "00");
              if (faultValues.length > 0) {
                ElMessage.warning(t("tip.SetParameterValuesSuccessfully") + ", " + t("tip.ResponseWithFaultCode"));
              } else {
                ElMessage.success(t("tip.SetParameterValuesSuccessfully"));
              }
            }
            return resolve(updateValues);
          })
          .catch((e) => {
            console.error(e);
            if (showMessage) {
              ElMessage.error(t("tip.RequestFailed"));
            }
            return resolve([]);
          });
      });
    },
    getDeviceParameterValue: async function ({
      sub,
      oid,
      value = undefined,
      showMessage = false,
      statusEnable = false,
      options = {},
    }) {
      const returnValues = await this.getDeviceParameterValues({
        sub: sub,
        oids: [oid],
        values: [{ oid: oid, value: value }],
        showMessage: showMessage,
        statusEnable: statusEnable,
        options: options,
      });
      if (returnValues != undefined && returnValues.length > 0) {
        const result = returnValues.find((v) => v.oid === oid);
        return result;
      }
      return undefined;
    },
    setDeviceParameterValue: async function ({
      sub,
      oid,
      value,
      showMessage = false,
      statusEnable = false,
      options = {},
    }) {
      const returnValues = await this.setDeviceParameterValues({
        sub: sub,
        values: [{ oid: oid, value: value }],
        showMessage: showMessage,
        statusEnable: statusEnable,
        options: options,
      });
      if (returnValues && returnValues.length > 0) {
        const result = returnValues.find((v) => v.oid === oid);
        return result;
      }
      return undefined;
    },
    selectDeviceSub: async function (sub) {
      if (this.currentSub != sub) {
        const info = this.getDeviceInfo(sub);
        if (info == undefined) {
          ElMessage.error("Invalid sub-device");
          return;
        }
        this.currentSub = sub;
        console.log("current", this.currentDeviceSub, this.currentDeviceInfo);
        await this.setupDevice(sub, info);
        await this.currentDevice.checkAvailable();
        document.title = this.currentDevice.layout.currentWebTitle;
      }
    },
  },
});
