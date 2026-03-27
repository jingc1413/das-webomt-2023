import model from "@/stores/model";
import { newDeviceLayoutManager } from "./layout";
import { newDeviceParametersManager } from "./parameters";
import { newDeviceStatsManager } from "./stats";
import { newDeviceFunctions } from "./funcs";
import { newDeviceFilesManager } from "./files";
import { newDeviceFirmwareManager } from "./firmwares";
import { newDeviceUpgradeManager } from "./upgrade";
import { newDeviceImportor } from "./importer";
import { newDeviceAlarms } from "./alarm";

import { newDeviceConfigurationManager } from "./configuration";
import settings from "@/settings";
import { ElMessage, dayjs } from "element-plus";
import { useDasModel } from "../das-model";
import { number } from "echarts";

export async function newDevice(sub, info) {
  const dev = {
    sub,
    info,
  };
  if (dev.info.ConnectState >= 6){
    return undefined;
  }

  const deviceModel = await useDasModel().getModel(dev.info.DeviceTypeName, dev.info.Version);
  if (!deviceModel){
    return undefined;
  }

  const paramsModel = model.setupDeviceParamsModel(dev.sub, deviceModel);
  const statsModel = model.setupDeviceStatsModel(dev.sub, deviceModel);
  const layoutModel = model.setupDeviceLayoutModel(dev.sub, deviceModel, paramsModel);

  if (!paramsModel || !statsModel || !layoutModel){
    return undefined;
  }

  dev.params = newDeviceParametersManager(dev.sub, dev.info, paramsModel);
  dev.layout = newDeviceLayoutManager(dev.sub, dev.info, dev.params, layoutModel);
  dev.stats = newDeviceStatsManager(dev.sub, dev.info, dev.params, statsModel);
  dev.files = newDeviceFilesManager(dev.sub);
  dev.cfg = newDeviceConfigurationManager(dev.sub, dev.info, dev.params, dev.layout, dev.files);
  dev.firmwares = newDeviceFirmwareManager(dev.sub, dev.params);
  dev.upgrade = newDeviceUpgradeManager(dev.sub, dev.info, dev.params);
  dev.importer = newDeviceImportor(dev.sub, dev.params);
  dev.funcs = newDeviceFunctions(dev.sub, dev.params);
  dev.alarm = newDeviceAlarms(dev.sub, dev.params);
  dev.setup = function () {
    const self = this;
    this.params.setParamValueUpdateCallback("T02.P0007", function ({ value }) {
      self.info.ElementSerialNumber = value;
    });
    this.params.setParamValueUpdateCallback("T02.P0023", function ({ value }) {
      self.info.InstalledLocation = value;
    });
    this.params.setParamValueUpdateCallback("T02.P0024", function ({ value }) {
      self.info.SiteName = value;
      self.layout.updateTitle();
    });
    this.params.setParamValueUpdateCallback("T02.P0030", function ({ value }) {
      self.info.DeviceName = value;
      self.layout.updateTitle();
    });
    this.params.setParamValueUpdateCallback("T02.P0102", function ({ value }) {
      self.info.SubID = value;
      self.layout.updateTitle();
    });
    this.params.setParamValueUpdateCallback("T02.P0150", function ({ raw, value }) {
      self.info.SystemTime = raw;
    });
    this.params.setParamValueUpdateCallback("T02.P002B", function ({ raw, value }) {
      self.info.LifeTime = Number(value);
    });
    this.params.setParamValueUpdateCallback("T02.P002C", function ({ value }) {
      const args = value.split(":");
      if (args?.length === 3) {
        let uptime = 0;
        uptime += 86400 * Number(args[0]);
        uptime += 3600 * Number(args[1]);
        uptime += 60 * Number(args[2]);
        self.info.UpTime = uptime;
      }
    });
    this.params.setParamValueUpdateCallback("TB0.P0AEF", function ({ value }) {
      self.resetInputValue();
    });
    // setTimeout(() => {
    //   self.funcs.isFacotryMode();

    //   self.params.getParameterValues({
    //     oids: ["T02.P0007", "T02.P0024", "T02.P0024", "T02.P0030", "T02.P0102"],
    //     statusEnable: true,
    //   });
    // }, 500);
  };

  dev.setAvailable = function (available) {
    const now = dayjs();
    if (this.info.state === undefined) {
      this.info.state = {};
    }
    if (this.info.state.available === undefined || available != this.info.state.available) {
      if (available === true && this.info.state.available === false) {
        if (this.info.state.availableTime === undefined) {
          this.info.state.availableTime = now;
        }
      }
      if (this.info.state.availableTime && now.diff(this.info.state.availableTime) < 30 * 1000) {
        return;
      }
      this.info.state.availableTime = undefined;
      this.info.state.available = available;
      console.log("set available", this.sub, available);
      if (!this.info.state?.available && !settings.nodeTest) {
        setTimeout(() => {
          this.checkAvailable(sub);
        }, 1000);
      }
    }
    return;
  };
  dev.checkAvailable = async function () {
    const self = this;
    // console.log("!!!! check available", sub, info)
    if (this.info.state === undefined) {
      this.info.state = {};
    }
    try {
      this.info.state.availableChecking = true;
      const start = dayjs();
      await self.doCheckAvailable(sub);
      // const diff = dayjs().diff(start);
      // if (5000 > diff) {
      //   await sleep(5000 - diff);
      // }
      if (!this.info.state.available) {
        setTimeout(() => {
          self.checkAvailable(sub);
        }, 5000);
      }
    } catch (e) {
      console.log(e);
      setTimeout(() => {
        self.checkAvailable(sub);
      }, 5000);
    } finally {
      this.info.state.availableChecking = false;
    }
  };
  dev.doCheckAvailable = async function () {
    if (settings.nodeTest) {
      this.setAvailable(true);
      return;
    }
    await this.params.getParameterValue({
      oid: "T02.P0101",
      options: {
        timeout: 3000,
      },
      statusEnable: true,
    });
  };

  dev.doAction = async function (action, options = {}) {
    if (action == undefined) {
      return;
    }
    switch (action.Name) {
      case "SetDeviceUnavailable": {
        this.setAvailable(false);
        break;
      }
      case "GetParameterValues": {
        const oids = [];
        const values = [];
        action?.Items?.forEach((v) => {
          const param = this.params.getParam(v.OID);
          if (param) {
            oids.push(v.OID);
            if (v.Value !== undefined) {
              values.push({ oid: v.OID, value: v.Value });
            }
          }
        });
        await this.params.getParameterValues({
          oids,
          values: values,
          showMessage: action?.Style?.showMessage,
        });
        break;
      }
      case "SetParameterValues": {
        const values = [];
        action?.Items?.forEach((v) => {
          const param = this.params.getParam(v.OID);
          if (v.Value !== undefined) {
            values.push({ oid: v.OID, value: v.Value });
          } else {
            values.push({ oid: v.OID, value: param.InputValue });
          }
        });
        await this.params.setParameterValues({ values, showMessage: action?.Style?.showMessage });
        break;
      }
      case "ViewPage": {
        const page = action.Items[0];
        if (page && options.owner) {
          this.layout.openViewPage({
            page: options.owner,
            viewPage: page,
            onCloseViewPage: options["onCloseViewPage"],
          });
        }
        break;
      }
      case "MultipleActions": {
        for (const i in action.Items) {
          await this.doAction(action.Items[i], options);
        }
        break;
      }
      case "StateAction": {
        const k = options.param?.Value;
        if (k) {
          const v = action.Actions[k];
          if (v) {
            this.doAction(v, options);
          }
        }
        break;
      }
    }
  };
  dev.setup();
  return dev;
}
