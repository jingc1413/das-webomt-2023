import { ElMessage, ElLoading } from "element-plus";
import { translator as t } from "@/i18n";

export function newDeviceConfigurationManager(sub, info, params, layout, files) {
  const m = {
    sub,
    info,
    params,
    layout,
    files,
  };

  m.exportConfigurationData = async function (data) {
    const body = {
      version: 2,
      deviceTypeName: this.info.DeviceTypeName,
      productModel: this.info.ProductModel,
      data: [],
    };
    if (data) {
      data.forEach((item) => {
        const groupData = [];
        item.data.forEach((item2) => {
          const param = this.params.getParam(item2.oid);
          if (param) {
            if (item2.value !== undefined) {
              const v = item2.value;
              groupData.push({ id: param.PrivOid, name: param.Name, value: v });
            } else {
              const v = param.Value;
              groupData.push({ id: param.PrivOid, name: param.Name, value: v });
            }
          }
        });
        body.data.push({
          id: item.path.join("."),
          data: groupData,
        });
      });
    }
    return body;
  };
  m.importConfigurationData = function (body) {
    const self = this;
    const data = [];
    if (!body) {
      return data;
    }
    if (body.version) {
      if (body.version === 2) {
        if (this.info.DeviceTypeName != body.deviceTypeName) {
          throw Error("mismatch device type name");
        }
        if (this.info.ProductModel != body.productModel) {
          throw Error("mismatch product model");
        }
        body.data.forEach((group) => {
          group.data.forEach((p) => {
            p.param = self.params.getParam(p.id);
            if (p.param) {
              const v = p.value;
              let item = data.find((v) => v.path.join(".") === group.id);
              if (!item) {
                item = {
                  path: group.id.split("."),
                  data: [],
                };
                data.push(item);
              }
              item.data.push({ oid: p.id, value: v, param: p.param });
            }
          });
        });
      }
    } else if (body.matchineroot) {
      const parseParam = function (module, page, group, p) {
        if (p.value !== undefined && p.limit === "nrw") {
          const oldPath = [module.name_en, page.name_en, group.name_en].join(",");
          const param = self.params.getParamByOldID(oldPath, p.id);
          if (param && param.Paths?.length > 0) {
            const path = this.layout.matchPageTreePath(param.Paths[0]);
            if (path) {
              let item = data.find((v) => v.path.join(".") === path.join("."));
              let v = p.value;
              if (String(v) === "true") {
                v = param.IsNumberDataType ? 1 : "01";
              } else if (String(v) === "false") {
                v = param.IsNumberDataType ? 0 : "00";
              } else if (param.Options) {
                v = Object.keys(param.Options).find(
                  (k) => String(k) === String(v) || String(param.Options[k]) === String(v)
                );
              } else if (param.IsNumberDataType) {
                v = Number(v);
              }
              if (!item) {
                item = { path: path, data: [] };
                data.push(item);
              }
              if (v !== undefined) {
                item.data.push({ oid: param.PrivOid, value: v, param: param });
              } else {
                console.log("import invalid parameter value", v, p);
              }
            } else {
              console.log("import unknow parameter path", param.Paths);
            }
          } else {
            console.log("import unknow parameter", oldPath, p);
          }
        }
      };
      Object.keys(body.matchineroot).forEach((key) => {
        const module = body.matchineroot[key];
        if (module.nodetype === "module") {
          Object.keys(module).forEach((key2) => {
            const page = module[key2];
            if (page.nodetype === "page" && page.group?.length > 0) {
              page.group.forEach((group) => {
                if ((group.type === "param" || group.type === "alarm") && group.dataTable?.length > 0) {
                  group.dataTable.forEach((p) => {
                    if (p.checkbox === true) {
                      parseParam(module, page, group, p);
                      if (p.col2) {
                        parseParam(module, page, group, p.col2);
                      }
                    }
                  });
                }
              });
            }
          });
        }
      });
    }
    return data;
  };

  m.supportCarrierConfig = async function () {
    const param = this.params.getParam("T02.P0880", "Carrinfo export");
    const param2 = this.params.getParam("T02.P0572", "Carrinfo upload");
    if (param && param2) {
      return true;
    }
    return false;
  };
  m.supportCustomCarrierConfig = async function () {
    const param = this.params.getParam("T02.P0883", "Customcarrier export");
    const param2 = this.params.getParam("T02.P0882", "Customcarrier import");
    if (param && param2) {
      return true;
    }
    return false;
  };
  m.getCarrierConfigFile = async function () {
    const param = this.params.getParam("T02.P0880", "Carrinfo export");
    if (param == undefined) {
      ElMessage.error("Not support to export the configuration file");
      return;
    }
    const result = await this.params.setParameterValue({
      oid: "T02.P0880",
      value: "00",
    });
    if (result && result.code === "00") {
      const self = this;
      const filename = result.value;
      const filetype = "ConfigFile";
      this.files.downloadFile(filetype, filename);
      setTimeout(() => {
        self.files.deleteFile(filetype, filename, false);
      }, 2000);
    } else {
      ElMessage.error("Export the configuration file failed");
      return;
    }
  };
  m.getCustomCarrierConfigFile = async function () {
    const param = this.params.getParam("T02.P0883", "Customcarrier export");
    if (param == undefined) {
      ElMessage.error("Not support to export the configuration file");
      return;
    }
    const result = await this.params.setParameterValue({
      oid: "T02.P0883",
      value: "00",
    });
    if (result && result.code === "00") {
      const self = this;
      const filename = result.value;
      const filetype = "ConfigFile";
      this.files.downloadFile(filetype, filename);
      setTimeout(() => {
        self.files.deleteFile(filetype, filename, false);
      }, 2000);
    } else {
      ElMessage.error("Export the configuration file failed");
      return;
    }
  };
  m.loadPAInitConfigFile = async function (filename) {
    const loadingInstance1 = ElLoading.service({
      fullscreen: true,
      text: "Loading PAInit Configuration ...",
    });
    try {
      const oid = "TB2-PA.P0060";
      const param = this.params.getParam(oid);
      if (param == undefined) {
        ElMessage.error("Not support to load the configuration file");
        return;
      }
      await this.params.setParameterValue({ oid: oid, value: filename });
      ElMessage.success("Configuration file load successfully");
    } catch (e) {
      console.error(e);
      ElMessage.error(t("tip.RequestFailed"));
    } finally {
      loadingInstance1.close();
    }
  };
  m.loadFlatnessCoefficientConfigFile = async function (filename) {
    const loadingInstance1 = ElLoading.service({
      fullscreen: true,
      text: "Loading Flatness Coefficient Configuration ...",
    });
    try {
      const oid = "T02.P0831";
      const param = this.params.getParam(oid);
      if (param == undefined) {
        ElMessage.error("Not support to load the configuration file");
        return;
      }
      const retult = await this.params.setParameterValue({
        oid: oid,
        value: filename,
      });
      if (retult?.code === "00") {
        ElMessage.success("Configuration file load successfully");
      } else {
        ElMessage.error("Configuration file load failed");
      }
    } catch (e) {
      console.error(e);
      ElMessage.error(t("tip.RequestFailed"));
    } finally {
      loadingInstance1.close();
    }
  };
  m.loadCarrierConfigFile = async function (filename) {
    const loadingInstance1 = ElLoading.service({
      fullscreen: true,
      text: "Loading Carrier Configuraion ...",
    });
    try {
      const oid = "T02.P0572";
      const param = this.params.getParam(oid, "Carrinfo Upload");
      if (param == undefined) {
        ElMessage.error("Not support to load the configuration file");
        return;
      }
      const retult = await this.params.setParameterValue({
        oid: oid,
        value: filename,
      });
      if (retult?.code === "00") {
        ElMessage.success("Configuration file load successfully");
      } else {
        ElMessage.error("Configuration file load failed");
      }
    } catch (e) {
      console.error(e);
      ElMessage.error(t("tip.RequestFailed"));
    } finally {
      loadingInstance1.close();
    }
  };
  m.loadCustomCarrierConfigFile = async function (filename) {
    const loadingInstance1 = ElLoading.service({
      fullscreen: true,
      text: "Loading Carrier Configuraion ...",
    });
    try {
      const oid = "T02.P0882";
      const param = this.params.getParam(oid, "Customcarrier import");
      if (param == undefined) {
        ElMessage.error("Not support to load the configuration file");
        return;
      }
      const retult = await this.params.setParameterValue({
        oid: oid,
        value: filename,
      });
      if (retult?.code === "00") {
        ElMessage.success("Configuration file load successfully");
      } else {
        ElMessage.error("Configuration file load failed");
      }
    } catch (e) {
      console.error(e);
      ElMessage.error(t("tip.RequestFailed"));
    } finally {
      loadingInstance1.close();
    }
  };
  return m;
}
