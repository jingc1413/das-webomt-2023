import { ElMessage } from "element-plus";
import apix from "@/api";
import { translator as t } from "@/i18n";

export function setupDeviceInfos({ onUpdateInfo = undefined, onRemoveInfo = undefined, onChangeInfos = undefined }) {
  const deviceInfos = {
    loadingState: {},
    infos: [],
    onRemoveInfo: onRemoveInfo,
    onUpdateInfo: onUpdateInfo,
    onChangeInfos: onChangeInfos,
  };

  deviceInfos.getInfo = function (sub) {
    return this.infos.find((v) => v.SubID !== undefined && v.SubID !== "" && String(v.SubID) === String(sub));
  };

  deviceInfos.setInfo = function (info, upsert = false) {
    if (info == undefined || info.SubID == undefined) {
      return;
    }
    const sub = String(info.SubID);
    const last = this.infos.find(
      (v) => v.SubID !== undefined && v.SubID !== "" && String(v.SubID) === String(info.SubID)
    );
    if (last === undefined && !upsert) {
      return;
    }
    if (last === undefined) {
      this.infos.push(info);
    } else {
      Object.keys(info).forEach((k) => {
        last[k] = info[k];
      });
    }

    if (this.onUpdateInfo != undefined) {
      this.onUpdateInfo(info.SubID);
    }
  };
  deviceInfos.removeInfo = function (sub) {
    this.infos = this.infos.filter((v) => String(v.SubID) != String(sub));
    if (this.onRemoveInfo != undefined) {
      this.onRemoveInfo(sub);
    }
  };

  deviceInfos.updateQueryStatus = function (status) {
    this.loadingState.total = status?.Total;
    this.loadingState.index = status?.Index;
    if (!status.Finished) {
      this.loadingState.loading = true;
      this.loadingState.finished = false;
      if (status?.Index > 0 && status?.Total > 0) {
        this.loadingState.progress = parseFloat((status.Index * 100) / status.Total);
      } else {
        this.loadingState.progress = 0;
      }
      setTimeout(() => {
        this.getQueryProgress();
      }, 1000);
    } else {
      if (this.loadingState.loading && !this.loadingState.finished) {
        this.loadingState.finished = true;
        if (!status.Success) {
          ElMessage.error("Query device failed");
        }
        this.updateInfos();
        setTimeout(() => {
          this.loadingState.loading = false;
        }, 1000);
      }
    }
  };

  deviceInfos.getQueryProgress = async function () {
    return;
    // apix.getQueryDevicesProgress().then(data => {
    //   this.updateQueryStatus(data);
    // }).catch(e => {
    //   console.log(e);
    // });
  };

  deviceInfos.updateInfo = async function (sub, showMessage = false) {
    apix
      .getDeviceInfo(sub)
      .then((info) => {
        this.setInfo(info);
      })
      .catch((e) => {
        if (showMessage) {
          ElMessage.error(t("tip.RequestFailed"));
        }
      });
  };

  deviceInfos.updateInfos = async function (force = false, showMessage = false) {
    try {
      const lastInfos = this.infos;
      let updateInfos = [this.localInfo];
      await apix
        .getDeviceInfos(force)
        .then((infos) => {
          updateInfos = infos;
        })
        .catch((e) => {
          console.log(e);
        });

      const removedDeviceInfos = lastInfos.filter((v) => {
        return (
          updateInfos.find((v2) => String(v2.SubID) === String(v.SubID) && v2.DeviceTypeName === v.DeviceTypeName) ===
          undefined
        );
      });
      const newDeviceInfos = updateInfos.filter((v) => {
        return (
          lastInfos.find((v2) => String(v2.SubID) === String(v.SubID) && v2.DeviceTypeName === v.DeviceTypeName) ===
          undefined
        );
      });
      // console.log({ last: lastInfos, removed: removedDeviceInfos, added: newDeviceInfos })

      removedDeviceInfos.forEach((v) => {
        this.removeInfo(v.SubID);
      });
      updateInfos.forEach((v) => {
        this.setInfo(v, true);
      });
      if (this.onChangeInfos != undefined) {
        this.onChangeInfos(this.infos);
      }

      if (force) {
        this.getQueryProgress();
      }
    } catch (e) {
      if (showMessage) {
        ElMessage.error(t("tip.RequestFailed"));
      }
    }
  };

  return deviceInfos;
}
