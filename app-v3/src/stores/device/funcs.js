import apix from "@/api";
import { ElMessage } from "element-plus";

export function newDeviceFunctions(sub, params) {
  const m = {
    sub,
    params,
  };

  m.enterFactoryMode = async function (showMessage = false) {
    const result = await this.params.setParameterValue({ oid: "TB0.P0AFF", value: "5A" });
    if (result && result.code === "00" && result.value === "5A") {
      return true;
    }
    if (showMessage) {
      ElMessage.error("Facotry mode enter failed");
    }
    return false;
  };
  m.quitFactoryMode = async function (showMessage = false) {
    const result = await this.params.setParameterValue({ oid: "TB0.P0AFF", value: "A5" });
    if (result && result.code === "00" && result.value === "A5") {
      return true;
    }
    if (showMessage) {
      ElMessage.error("Facotry mode quit failed");
    }
    return false;
  };
  m.isFacotryMode = async function () {
    await this.params.getParameterValues({
      oids: ["TB0.P0AFF"],
    });
    const param = this.params.getParam("TB0.P0AFF");
    if (param && param.Value === "5A") {
      return true;
    }
    return false;
  };

  m.resetDeviceAlarmState = async function (showMessage = false) {
    const param = this.params.getParamByName("Alarm Initialization");
    if (!param) {
      if (showMessage) {
        ElMessage.error("Alarm state reset failed");
      }
      return false;
    }
    const result = await this.params.setParameterValue({ oid: param.PrivOid, value: "FF" });
    if (result && result.code === "00" && result.value === "FF") {
      return true;
    } else {
      if (showMessage) {
        ElMessage.error("Alarm state reset failed");
      }
      return false;
    }
  };
  m.resetCpriSyncLossCounter = async function (showMessage = false) {
    const param = this.params.getParamByName("Sync Loss Counter Reset");
    if (!param) {
      if (showMessage) {
        ElMessage.error("CPRI sync loss counter reset failed");
      }
      return false;
    }
    const result = await this.params.setParameterValue({ oid: param.PrivOid, value: "00" });
    console.log(result);
    if (result && result.code === "00" && result.value === "00") {
      return true;
    } else {
      if (showMessage) {
        ElMessage.error("CPRI sync loss counter reset failed");
      }
      return false;
    }
  };
  m.deleteDeviceKeyAndLogs = async function () {
    return new Promise((resolve, reject) => {
      apix
        .deleteDeviceKeyAndLogs(this.sub)
        .then(() => {
          return resolve(true);
        })
        .catch((e) => {
          return resolve(false);
        });
    });
  };
  m.readRegister = async function (module, offset, size) {
    return new Promise((resolve, reject) => {
      apix
        .readDeviceRegister(this.sub, module, offset, size)
        .then((data) => {
          const out = [];
          for (let i = 0; i < data.Size; i +=1) {
            out.push({
              offset: offset + i,
              value: parseInt(data.Buffer.substring(i*2, i*2 + 2), 16),
            });
          }
          return resolve(out);
        })
        .catch((e) => {
          console.error(e);
          return resolve(undefined);
        });
    });
  };
  m.writeRegister = async function (module, offset, size, buffer) {
    let input = "";
    for (let i = 0; i < size; i +=1) {
      input += buffer[i].value.toString(16).padStart(2, "0");
    }
    return new Promise((resolve, reject) => {
      apix
        .writeDeviceRegister(this.sub, module, offset, size, input)
        .then((data) => {
          const out = [];
          for (let i = 0; i < data.Size; i +=1) {
            out.push({
              offset: offset + i,
              value: parseInt(data.Buffer.substring(i*2, i*2 + 2), 16),
            });
          }
          return resolve(out);
        })
        .catch((e) => {
          console.error(e);
          return resolve(undefined);
        });
    });
  };
  return m;
}
