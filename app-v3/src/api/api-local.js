import * as base from "./api-base";
import settings from "@/settings.js";

export function getAppInfo() {
  if (settings.nodeTest) {
    return Promise.resolve({
      // AppName: "CORNING | Everon™ 6200",
      AppShortName: "OMT",
      AppVersion: "0.0.0",
      AppBuild: "FFFFFF",
      Schema: "corning",
      DeviceTypeName: settings.DeviceTypeName,
    });
  }
  return base.httpGet(base.ApiBase + "/app/info");
}
export function getCurrentBackground(fileName) {
  if (settings.nodeTest) {
    return base.httpGet("mock/login/bg.jpg", null, { responseType: "blob" });
  }
  let url =  `img/login-bg${fileName?('.'+fileName):''}.jpg`;
  return base.httpGet(url, null, { responseType: "blob" });
}
