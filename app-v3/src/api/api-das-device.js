import * as base from "./api-base";
import settings from "@/settings";

export async function getQueryDevicesProgress() {
  const u = `${base.DasApiBase}/query-devices/progress`;
  return base.httpGet(u);
}
export async function getDeviceInfos(force = false) {
  if (settings.nodeTest) {
    return require("@/components/Topo/info.json");
  }
  const u = force ? `${base.DasApiBase}/devices?force` : `${base.DasApiBase}/devices`;
  return base.httpGet(u);
}
export async function getDeviceInfo(sub) {
  return base.httpGet(base.DasApiBase + "/devices/" + sub);
}

export async function getLocalInfo() {
  if (settings.nodeTest) {
    return {
      DeviceTypeName: settings.deviceInfo.currentType,
      Version: settings.deviceInfo.currentVersion,
    };
  }
  return base.httpGet(base.DasApiBase + "/devices/local");
}

export async function queryDeviceStats(sub, beginTime, endTime, keys = []) {
  if (settings.nodeTest) {
    return base.httpGet("/mock/systemStats/mock.json");
  }
  const body = {
    beginTime,
    endTime,
    keys,
  };
  return base.httpPost(base.DasApiBase + "/devices/" + sub + "/metrics/data/query", body);
}

export async function getDeviceAlarmLogs(sub) {
  if (settings.nodeTest) {
    return base.httpGet("/mock/alarm/alarmLog.json");
  }
  return base.httpGet(base.DasApiBase + "/devices/" + sub + "/alarm/logs");
}
