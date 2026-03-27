import * as base from "./api-base";
import settings from "@/settings";

export async function getDeviceFileList(sub, fileType) {
  if (settings.nodeTest) {
    if (fileType === "UpgradeFile") {
      return [
        {
          "FileName": "iDAS_N3RU_R404_V0.9_2AFE_20230829_Everon_6200_v3.0_build11.zip",
          "FileSize": 13182976,
          "ModTime": -4098780
        }
      ];
    } else {
      return [];
    }
  }
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/files/" + fileType);
}

export async function getDeviceFile(sub, fileType, fileName, options) {
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/files/" + fileType + "/" + fileName, undefined, options);
}

export async function deleteDeviceFile(sub, fileType, fileName) {
  return base.httpDelete(base.DasApiBase + "/devices/" + sub + "/files/" + fileType + "/" + fileName);
}

export async function getDeviceCurrentPacketInfo(sub) {
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/version/packet-info");
}


export async function getDeviceUpgradeFilePacketInfo(sub, fileName) {
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/files/UpgradeFile/" + fileName + "/packet-info");
}

export async function deleteDeviceKeyAndLogs(sub) {
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/delete-key-and-logs");
}

export async function getDeviceFirmwareList(sub) {
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/firmwares");
}

export async function deleteDeviceFirmware(sub, name) {
  return base.httpDelete(base.DasApiBase + "/devices/" + sub + "/firmwares/" + name);
}

export async function readDeviceRegister(sub, module, offset, size) {
  const body ={
    Module: module,
    Offset: offset,
    Size: size,
  }
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/register/read", body);
}

export async function writeDeviceRegister(sub, module, offset, size, buffer) {
  const body ={
    Module: module,
    Offset: offset,
    Size: size,
    Buffer: buffer,
  }
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/register/write", body);
}