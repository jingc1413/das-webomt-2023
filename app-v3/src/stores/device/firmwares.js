import { ElMessage } from "element-plus";
import apix from "@/api";
import { translator as t } from "@/i18n";

export function newDeviceFirmwareManager(sub, params) {
  const m = {
    sub,
    params,
    firmwares: [],
  };
  m.getFirmwareList = async function (showMessage=false) {
    await apix
      .getDeviceFirmwareList(this.sub)
      .then((data) => {
        this.firmwares = data;
        if (showMessage) {
          ElMessage.success(t("tip.getFirmwareListSuccessfully"));
        }
      })
      .catch((e) => {
        console.log(e);
        if (showMessage) {
          ElMessage.error(t("tip.getFirmwareListFailed"));
        }
      });
  };
  m.deleteFirmware = async function (firmware) {
    await apix
      .deleteDeviceFirmware(this.sub, firmware.Name)
      .then((data) => {
        ElMessage.success(t("tip.firmwareDeletedSuccessfully"));
        this.getFirmwareList();
      })
      .catch((e) => {
        console.error(e);
        ElMessage.error(t("tip.firmwareDeletedFailed"));
        this.getFirmwareList();
      });
  };
  return m;
}
