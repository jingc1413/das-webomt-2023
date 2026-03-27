import { defineStore } from "pinia";
import apix from "@/api";
// import model from "@/stores/model";

export const useDasModel = defineStore("dasModel", {
  state: () => ({
    auDeviceTypeName: "Primary A3",
    deviceTypeNames: ["E3-O", "N3-RU", "Primary A3", "Secondary A3", "X3-RU", "M3-RU-L", "M3-RU-H"],
    deviceModels: [],
    productModels: [],
  }),
  getters:{
    productModel(state){
      return function(deviceTypeName) {
        return state.productModels.find( model => model.DeviceTypeName === deviceTypeName)
      }
    },
  },
  actions: {
    setup: async function () {
      await apix
        .getDeviceTypeNames()
        .then((names) => {
          this.deviceTypeNames = names;
        })
        .catch((e) => {
          console.log(e);
        });
      await apix
        .getProductModels(this.deviceTypeNames)
        .then((models) => {
          this.productModels = models;
        })
        .catch((e) => {
          console.log(e);
        });
      // this.deviceModels = await apix.getDeviceModels(this.deviceTypeNames);
    },
    validateDeviceTypeName(deviceTypeName) {
      return this.deviceTypeNames.includes(deviceTypeName);
    },
    getModel: async function (deviceTypeName, version="latest") {
      if (version == "") {
        version = "latest"
      }
      let model = this.deviceModels.find((v) => v.type === deviceTypeName && v.version === version);
      if (model === undefined) {
        await apix.getDeviceModel(deviceTypeName, version).then((data) => {
          this.deviceModels.push(data);
        }).catch((e)=>{
          // console.log(e);
        });
        model = this.deviceModels.find((v) => v.type === deviceTypeName && v.version === version);
      }
      return model;
    },
  },
});
