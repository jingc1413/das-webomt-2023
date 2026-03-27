import apix from "@/api";
import { ElMessage, ElLoading } from "element-plus";
import { sleep } from "@/utils";
import { useAppStore } from "@/stores/app";
import { useDasModel } from "@/stores/das-model";

export function newDeviceUpgradeManager(sub, info, params) {
  const m = {
    sub,
    info,
    params,
    upgradeState: undefined,
  };

  m.upgradeProgress = function () {
    if (this.upgradeState?.result !== undefined) {
      return 100;
    } else if (this.upgradeState?.resultByArm !== undefined) {
      return this.upgradeState.upgradeProgress;
    }
    return undefined;
  };

  m.clearUpgradeState = function () {
    this.upgradeState = undefined;
  };
  m.startUpgrade = async function (filename, force = false) {
    const self = this;
    const isSubDevice = this.sub != 0;
    const appStore = useAppStore();
    const dasModel = useDasModel();

    if (isSubDevice && !force) {
      appStore.openConfirmDialog({
        title: "Upgrade is not allowed",
        content: `Please upgrade on ${dasModel.auDeviceTypeName}, otherwise will be synchronized by the  ${dasModel.auDeviceTypeName}`,
        supportCancel: false,
      });
      return;
    } else if (force) {
      appStore.openConfirmDialog({
        title: "Force Code",
        content: "Please input force code to continue",
        needInput: true,
        inputRule: [
          {
            required: true,
            validator: (rule, value, callback) => {
              if (value === "iDas") {
                return callback();
              }
              return callback(new Error("Please input force code"));
            },
            trigger: "blur",
          },
        ],
        callback: (ok) => {
          if (ok) {
            self.doStartDeviceUpgrade(filename, true);
          }
        },
      });
    } else {
      self.doStartDeviceUpgrade(filename, false);
      return;
    }
  };

  m.doStartDeviceUpgrade = async function (filename, force = false) {
    const dasModel = useDasModel();
    const appStore = useAppStore();

    const loadingInstance1 = ElLoading.service({ fullscreen: true, text: "Start upgrade ..." });

    const isSnmp = filename.indexOf("SNMP") > -1;
    const byArm = this.info.DeviceTypeName === dasModel.auDeviceTypeName;
    console.log("start upgrade", { sub: this.sub, upgrade: filename, force, byArm, isSnmp });

    this.upgradeState = {
      filename: filename,
      force: force,
      byArm: byArm,
      upgrading: true,
      finished: false,
      resultStatus: undefined,
      resultText: undefined,

      resp: undefined,
      respByArm: undefined,
      progress: 0,
      stats: undefined,
      logs: [],
    };

    return new Promise((resolve, reject) => {
      apix
        .startDeviceUpgrade(this.sub, filename, force, byArm)
        .then((resp) => {
          let cmdOid = "TB4.P0B19";
          if (byArm) {
            if (force) {
              cmdOid = "TB4.P0B22";
            } else if (isSnmp) {
              cmdOid = "TB4.P0B25";
            }
            this.params
              .seteParameterValue({ oid: cmdOid, value: "00" })
              .then((v) => {
                if (v.code !== "00") {
                  appStore.openConfirmDialog({
                    title: "Upgrade failed",
                    content: "The upgrade package is not found, please try again",
                    supportCancel: false,
                  });
                  this.clearDeviceUpgradeState();
                  return resolve(false);
                }
              })
              .catch((e) => {
                console.log(e);
                appStore.openConfirmDialog({
                  title: "Upgrade failed",
                  content: e.message,
                  supportCancel: false,
                });
                this.clearDeviceUpgradeState();
                return resolve(false);
              });
          } else {
            apix
              .startDeviceUpgrade(this.sub, filename, force, false)
              .then((resp) => {
                if (resp.Status === "success") {
                  this.upgradeState.resp = resp;
                  appStore.gotoUpgradeRoute(this.sub);
                  return resolve(true);
                } else {
                  appStore.openConfirmDialog({ title: "Upgrade failed", content: resp.Message, supportCancel: false });
                  this.upgradeState = undefined;
                  return resolve(false);
                }
              })
              .catch((e) => {
                console.log(e);
                appStore.openConfirmDialog({
                  title: "Upgrade failed",
                  content: e.message,
                  supportCancel: false,
                });
                this.clearDeviceUpgradeState();
                return resolve(false);
              });
          }
        })
        .catch((e) => {
          console.log(e);
          ElMessage.error("upgrade error");
          this.clearDeviceUpgradeState();
          return resolve(false);
        })
        .finally(() => {
          loadingInstance1.close();
        });
    });
  };

  m.setUpgradeFinished = function (status, text, progress) {
    if (this.upgradeState) {
      if (progress !== undefined) {
        this.upgradeState.progress = progress;
      }
      this.upgradeState.resultStatus = status;
      this.upgradeState.resultText = text;
      this.upgradeState.finished = true;
    }
  };

  m.checkUpgrading = async function () {
    if (this.upgradeState?.finished) {
      return;
    }
    if (this.upgradeState?.resp !== undefined) {
      const resp = this.upgradeState.resp;
      const result = resp.result;
      if (result.needUpdateCRC && resp.crc != undefined) {
        const returnValue = await this.params.setParameterValue({ oid: "TB4.P0B20", value: resp.crc });
        if (returnValue?.code === "00") {
          this.upgradeState.logs.push({ status: "success", text: "The package CRC updated successfully" });
        } else {
          this.upgradeState.logs.push({ status: "error", text: "The package CRC updated error" });
        }
      }
      this.upgradeState.progress = 40;
      sleep(200);

      if (result.needUpdateHostPackage) {
        const returnValue = await this.params.setParameterValue({ oid: "TB4.P0B19", value: "01" });
        if (returnValue?.code === "00") {
          this.upgradeState.logs.push({
            status: "success",
            text: "The package file of host device updated successfully",
          });
        } else {
          this.upgradeState.logs.push({ status: "error", text: "The package file of host device updated error" });
        }
      } else if (result.needUpdateSubPackage) {
        const returnValue = await this.params.setParameterValue({ oid: "TB4.P0B21", value: "01" });
        if (returnValue?.code === "00") {
          this.upgradeState.logs.push({
            status: "success",
            text: "The package file of sub-device updated successfully",
          });
        } else {
          this.upgradeState.logs.push({ status: "error", text: "The package file of sub-device updated error" });
        }
      }
      this.upgradeState.progress = 80;
      sleep(200);

      if (result.needTurnOffFPGA) {
        console.log("need to turn off fpga");
      }
      this.upgradeState.progress = 90;
      sleep(200);

      if (result.needReboot) {
        await apix.setDeviceUpgradeToReboot(this.sub).catch((e) => {
          this.upgradeState.logs.push({ status: "error", text: "The device called to reboot error" });
        });
      }
      this.upgradeState.progress = 100;

      if (result.needReboot) {
        this.setUpgradeFinished("success", "The package upgrade is complete. The device will restart!");
      } else {
        this.setUpgradeFinished("success", "The package upgrade is complete.");
      }
      return;
    } else if (this.upgradeState?.force) {
      const result = await this.getForceUpgradeProgress();
      if (result === true) {
        this.setUpgradeFinished("success", "The package upgrade is complete.", 100);
      } else if (result === false) {
        this.setUpgradeFinished("error", "The package upgrade is failed.");
      }
    } else {
      const resp = await this.getUpgradeProgress();
      if (resp === undefined) {
        return;
      }
      console.log("upgrade progress", resp);
      if (this.upgradeState === undefined) {
        this.upgradeState = {
          upgrading: true,
          finished: false,
          resultStatus: undefined,
          resultText: undefined,
          resp: undefined,
          respByArm: undefined,
          progress: 0,
          stats: {},
          logs: [],
        };
      }
      this.upgradeState.respByArm = resp;
      this.upgradeState.progress = resp.upgradeProgress;
      this.upgradeState.stats.totalDeviceNum = resp.totalDeviceNum;
      this.upgradeState.stats.timeoutDeviceNum = resp.timeoutDeviceNum;
      this.upgradeState.stats.failedDeviceNum = resp.failedDeviceNum;
      // this.upgradeState.stats.successDeviceNum = resp.totalDeviceNum - resp.failedDeviceNum;

      if (!resp.isUpgrading) {
        switch (resp.errorCode) {
          case 0:
            if (this.upgradeState.stats?.failedDeviceNum > 0) {
              this.setUpgradeFinished("warning", "Upgrade progress completed with failed devices.", 100);
            } else {
              this.setUpgradeFinished("success", "Upgrade progress completed.", 100);
            }
            break;
          case 1:
            this.setUpgradeFinished("error", "Upgrade progress failed, the package CRC is invalid!");
            break;
          case 3:
            this.setUpgradeFinished("error", "Upgrade progress failed, the package file is error!");
            break;
          case 4:
            this.setUpgradeFinished("error", "Upgrade progress failed, timeout");
            break;
          case 5:
            this.setUpgradeFinished(
              "warning",
              "Device does not exist, or the sub-devices in the system are the same as the upgrade package version.",
              100
            );
            break;
        }
      }
    }
  };
  m.getForceUpgradeProgress = async function () {
    const returnValue = await this.params.getParameterValue({ oid: "TB4.P0B24" });
    if (returnValue != undefined && returnValue.code === "00") {
      if (returnValue.value === "00") {
        return true;
      } else if (returnValue.value === "01") {
        return undefined;
      }
    }
    return false;
  };
  m.getUpgradeProgress = async function () {
    const returnValue = await this.params.getParameterValue({ oid: "TB4.P0B23" });
    if (returnValue != undefined && returnValue.code === "00") {
      const args = returnValue.value;
      const result = {};
      if (args.length == 8) {
        result.timeoutDeviceNum = args[0];
        result.failedDeviceNum = args[1];
        result.upgradeProgress = args[2];
        result.errorCode = args[3];
        result.isUpgrading = args[4];
        result.downloadStatus = args[5];
        result.totalDeviceNum = args[6];
      }
      return result;
    }
    return undefined;
  };
  return m;
}
