import { ElMessage } from "element-plus";
import { translator as t } from "@/i18n";
import { exportFileByALink } from "@/utils/request";
import apix from "@/api";

export function newDeviceFilesManager(sub) {
  const m = {
    sub: sub,
    allFiles: {},
  };
  m.setup = function () {
    this.allFiles = {};
  };

  m.fileList = function (fileType) {
    return this.allFiles[fileType];
  };

  m.getFileList = async function (filetype, showMessage = false) {
    await apix
      .getDeviceFileList(this.sub, filetype)
      .then((fileList) => {
        if (filetype === "UpgradeFile" && fileList !== undefined) {
          fileList.forEach((v) => {
            const name = v.FileName;
            const subs = name.match(
              // eslint-disable-next-line no-useless-escape
              /^(iDAS|DDAS)_([a-zA-Z0-9\-]+)_([a-zA-Z0-9]+)_([a-zA-Z0-9\-.]+)_([a-fA-F0-9]{4})_.*$/
            );
            if (subs?.length == 6) {
              v.ProductType = subs[2];
              v.ProductModel = subs[3];
              v.Version = subs[4];
              v.CRC = subs[5];
            }
          });
        }
        this.allFiles[filetype] = fileList;
      })
      .catch((e) => {
        console.log(e);
        if (showMessage) {
          ElMessage.error(t("tip.RequestFailed"));
        }
      });
  };
  m.deleteFile = async function (filetype, filename, showMessage = false) {
    return new Promise((resolve, reject) => {
      apix
        .deleteDeviceFile(this.sub, filetype, filename)
        .then(() => {
          this.getDeviceFileList(filetype);
          if (showMessage) {
            ElMessage.success(t("tip.fileDeletedSuccessfully"));
          }
          return resolve(true);
        })
        .catch((e) => {
          if (showMessage) {
            ElMessage.error(t("tip.RequestFailed"));
          }
          return resolve(false);
        });
    });
  };
  m.getFile = async function (filetype, filename, showMessage = false) {
    return new Promise((resolve, reject) => {
      apix
        .getDeviceFile(this.sub, filetype, filename)
        .then((content) => {
          return resolve(content);
        })
        .catch((e) => {
          if (showMessage) {
            ElMessage.error(t("tip.RequestFailed"));
          }
          return resolve(undefined);
        });
    });
  };
  m.downloadFile = async function (filetype, filename) {
    const url = apix.DasApiBase + "/devices/" + this.sub + "/files/" + filetype + "/" + filename;
    exportFileByALink(url, filename);
  };
  m.getCurrentPacketInfo = async function (showMessage = false) {
    return new Promise((resolve, reject) => {
      apix
        .getDeviceCurrentPacketInfo(this.sub)
        .then((content) => {
          return resolve(content);
        })
        .catch((e) => {
          if (showMessage) {
            ElMessage.error(t("tip.RequestFailed"));
          }
          return resolve(undefined);
        });
    });
  };
  m.getUpgradeFilePacketInfo = async function (filename = undefined, showMessage = false) {
    return new Promise((resolve, reject) => {
      apix
        .getDeviceUpgradeFilePacketInfo(this.sub, filename)
        .then((content) => {
          return resolve(content);
        })
        .catch((e) => {
          if (showMessage) {
            ElMessage.error(t("tip.RequestFailed"));
          }
          return resolve(undefined);
        });
    });
  };
  return m;
}
