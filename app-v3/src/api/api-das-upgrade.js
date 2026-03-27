import * as base from "./api-base";

export async function setDeviceUpgradeToReboot(sub) {
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/upgrade/reboot", {});
}

export async function startDeviceUpgrade(sub, filename, force = false, byArm = false) {
  const data = {
    Filename: filename,
    Force: force,
    ByArm: byArm,
  };
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/upgrade/start", data);
}

// let UPGRADE_FAILED = -1;
// let UPGRADE_SUCCESS = 1;
// let UPGRADE_SUCCESS_DAS = 2;
// let UPGRADE_SUCCESS_DAS_REBOOT = 3;
// let UPGRADE_SUCCESS_WEBOMT = 4;
// let UPGRADE_SUCCESS_FPGA = 5;
// let UPGRADE_SUCCESS_ARM = 6;
// let UPGRADE_SUCCESS_PA = 7;
// let UPGRADE_SUCCESS_DAS_CRC = 8;
// let UPGRADE_NORMAL_CANT = 9;
// let UPGRADE_SUCCESS_AUTO = 10;
// let UPGRADE_SUCCESS_SNMP = 11;
// let UPGRADE_SUCCESS_LINUX = 12;

// function parseStartUpgradeResponseData(data) {
//   //,sRstText,sDebug
//   //<!debug>>...-->   <!rsttext>>...-->
//   let stemp;
//   let npos1, npos2, npos3;
//   let sRstText = "";
//   let sDebug = "";
//   let sImg = "";
//   let sMainText = "";
//   let sResultCode = "";

//   let sSrc = data
//   while (sSrc.indexOf("<!") != -1) {
//     npos1 = sSrc.indexOf("<!");
//     npos2 = sSrc.indexOf(">>", npos1);
//     if ((npos3 = sSrc.indexOf("-->", npos2)) == -1) {
//       break;
//     }
//     stemp = sSrc.substring(npos1 + 2, npos2);
//     if (stemp == "debug") {
//       sDebug += sSrc.substring(npos2 + 2, npos3) + "<br>";
//     } else if (stemp == "rsttext") {
//       sRstText += sSrc.substring(npos2 + 2, npos3);
//     } else if (stemp == "img") {
//       sImg += sSrc.substring(npos2 + 2, npos3);
//     } else if (stemp == "rstmain") {
//       sMainText += sSrc.substring(npos2 + 2, npos3);
//     } else if (stemp == "resultcode") {
//       sResultCode = sSrc.substring(npos2 + 2, npos3);
//     }
//     sSrc = sSrc.substring(npos3, sSrc.length);
//   }

//   let iCRC = "";
//   let crc_npos1 = data.indexOf("(CRC:");
//   let crc_npos2 = data.indexOf(")", crc_npos1 + 1);
//   if (crc_npos1 >= 0 && crc_npos2 > crc_npos1) {
//     iCRC = data.substring(crc_npos1 + 5, crc_npos2);
//   }

//   const resp = {
//     code: parseInt(sResultCode),
//     statusText: sMainText,
//     infoText: sRstText,
//     debugText: sDebug,
//     crc: parseInt(iCRC),
//   }
//   parseStartUpgradeResponseResult(resp);
//   return resp;
// }

// function parseStartUpgradeResponseResult(resp) {
//   const result = {}
//   result.success = false;
//   switch (resp.code) {
//     case UPGRADE_NORMAL_CANT:
//       result.title = "Upgrade not allow!";
//       break;
//     case UPGRADE_FAILED:
//       result.title = "Upgrading Failed!";
//       result.msg = "The device will restart after 2 minutes";
//       result.needReboot = true;
//       break;

//     case UPGRADE_SUCCESS: //old version
//       result.title = "Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       result.needTurnOffFPGA = true;
//       break;

//     case UPGRADE_SUCCESS_DAS: // au upgrade slave
//       result.title = "Upgrading Succeeded!"
//       result.msg = "Upgrade slave device";
//       result.success = true;
//       result.needUpdateSubPackage = true;//0b21

//       // this.OnUpdateDasAll(upgradeDialog, 0);
//       break;
//     case UPGRADE_SUCCESS_DAS_CRC: // slave upgrade slave
//       result.title = "Upgrading Succeeded!"
//       result.msg = "Saving Current Software CRC check";
//       result.msg += "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       result.needUpdateCRC = true;
//       break;
//     case UPGRADE_SUCCESS_DAS_REBOOT: // au upgrade au
//       result.title = "Upgrading Succeeded!"
//       result.msg = "Local Device Upgrading Succeeded, Upgrading Slave Devices,";
//       result.msg += "The device will restart, please wait!";
//       result.success = true;
//       result.needTurnOffFPGA = true;
//       result.needReboot = true;
//       result.needUpdateHostPackage = true; //0b19
//       // this.OnUpdateDasAll(upgradeDialog, 1);
//       break;
//     case UPGRADE_SUCCESS_LINUX:
//       result.title = "LINUX Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       result.needUpdateHostPackage = true;
//       // this.OnUpdateDasAll(upgradeDialog, 1);
//       break;

//     case UPGRADE_SUCCESS_WEBOMT:
//       result.title = "WEBOMT Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       break;
//     case UPGRADE_SUCCESS_FPGA:
//       result.title = "FPGA Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       break;
//     case UPGRADE_SUCCESS_ARM:
//       result.title = "ARM Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       break;
//     case UPGRADE_SUCCESS_SNMP:
//       result.title = "SNMP Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       break;
//     case UPGRADE_SUCCESS_PA:
//       result.title = "485 sub module Upgrading Succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       result.needReboot = true;
//       break;
//     case UPGRADE_SUCCESS_AUTO:
//       result.title = "LONG enhanced version upgrading succeeded!"
//       result.msg = "The device will restart, please wait!";
//       result.success = true;
//       break;
//   }
//   resp.result = result;
//   // this.upgradeDialogText(upgradeDialog, result, sRet, sText)
// }

// function parseUpgradeLogsData(data) {
//   const lines = data.split("\n");
//   const result = [];
//   for (var i = 1; i < lines.length; i++) {
//     if (lines[i] === "") {
//       continue;
//     }
//     const fields = lines[i].split("|");
//     const fullLabel = fields[2];

//     const record = {
//       key: fields[0],
//       step: parseInt(fields[1]),
//       fullLabel: fullLabel,
//     };
//     const subs = fullLabel.match(/resultcode:(\d+)/);
//     if (subs && subs[1]) {
//       record.status = parseInt(subs[1]);
//     }
//     record.crc = fullLabel.match(/CRC:([0-9a-fA-F]{4})/)[1];
//     const subs2 = fullLabel.match(/upgradeResultCode:([0-9a-fA-F]+)/);
//     if (subs2 && subs2[1]) {
//       record.code = parseInt(subs[1], 16);
//     }
//     if (fullLabel.indexOf('The upgrade is complete') != -1) {
//       record.done = true;
//     }

//     result.push(record);
//   }
//   return result;
// }