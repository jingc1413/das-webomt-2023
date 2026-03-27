import * as base from "./api-base";
import settings from "@/settings";

export async function getProductModel(deviceTypeName) {
  return new Promise((resolve, reject) => {
    base
      .httpGet(base.DasApiBase + "/products/" + deviceTypeName)
      .then((resp) => {
        return resolve(resp);
      })
      .catch((e) => {
        if (deviceTypeName === "Primary A3") {
          return resolve({ Type: "A", ProductTypeName: "AU" });
        } else if (deviceTypeName.endsWith("RU")) {
          return resolve({ Type: "R", ProductTypeName: "RU" });
        }
        return reject(e);
      });
  });
}

export async function getProductModels(deviceTypeNames) {
  const list = [];
  deviceTypeNames.forEach((name) => {
    list.push(getProductModel(name));
  });
  return new Promise((resolve, reject) => {
    Promise.all(list)
      .then((all) => {
        return resolve(all);
      })
      .catch((e) => {
        return reject(e);
      });
  });
}

let modelCache = [];
export async function getDeviceModel(type, version, force = false) {
  if (!type) {
    return Promise.reject(undefined);
  }
  if (!version) {
    version = "latest";
  }
  let layoutUrl = base.DasApiBase + "/device-types/" + type + "/" + version + "/model/layout";
  let paramsUrl = base.DasApiBase + "/device-types/" + type + "/" + version + "/model/parameters";
  if (settings.nodeTest) {
    layoutUrl = "/mock/models/" + type + "/" + version + "/layout.json";
    paramsUrl = "/mock/models/" + type + "/" + version + "/parameters.json";
  }
  if (!force) {
    const model = modelCache.find((v) => v.type === type && v.version === version);
    if (model) {
      return Promise.resolve(model);
    }
  }
  return new Promise((resolve, reject) => {
    Promise.all([base.httpGet(layoutUrl), base.httpGet(paramsUrl)])
      .then((all) => {
        let model = modelCache.find((v) => v.type === type && v.version === version);
        if (model) {
          model.layout = all[0];
          model.parameters = all[1];
        } else {
          model = {
            type: type,
            version: version,
            layout: all[0],
            parameters: all[1],
          };
          modelCache.push(model);
          console.log("load", type, version, model);
        }
        return resolve(model);
      })
      .catch((e) => {
        return reject(e);
      });
  });
}

export async function getDeviceModels(infos) {
  const promises = [];
  for (const info of infos) {
    promises.push(getDeviceModel(info.type, info.version));
  }
  const models = [];
  await Promise.allSettled(promises).then((results) => {
    results.forEach((result) => {
      if (result.status === "fulfilled" && result.value) {
        models.push(result.value);
      } else if (result.status === "rejected") {
        console.log(result.reason);
      }
    });
  });
  return models;
}

export async function getDeviceTypeNames() {
  return base.httpGet(base.DasApiBase + "/device-types");
}
