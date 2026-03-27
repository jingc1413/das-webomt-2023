import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus';
import { useDasDevices } from './das-devices';
import Cookies from 'js-cookie'

export const useDasTopo = defineStore('dasTopo', {
  state: () => ({
    graphTopoData: undefined,
    treeTopoData: [],
    topoTypeModel: Cookies.get('topeType') || "Graph",
    treeDeviceName: Cookies.get('secondContent') || "Device Name",
  }),
  getters: {
  },
  actions: {
    deleteTopoNode: async function (sub) {
      const deviceInfo = await useDasDevices().getDeviceInfo(sub);
      if (deviceInfo) {
        const id = Number(deviceInfo.SubID);
        if (!isNaN(id)) {
          const routeAddress = deviceInfo.RouteAddress.split(".").map(v => {
            const v2 = Number(v);
            if (!isNaN(v2)) {
              return v2.toString(16).padStart(2, "0");
            }
            return "00";
          })
          const value = id.toString(16).padStart(2, "0") + routeAddress.join("");
          const result = await useDasDevices().setDeviceParameterValue({ sub: 0, oid: "TB4.P0B17", value: value })
          if (result && result.code != "00" && result.msg) {
            if (result.msg.erorr) {
              ElMessage.error(result.msg.error);
            } else if (result.msg.warning) {
              ElMessage.warning(result.msg.warning);
            }
            return false;
          }
          if (sub == 0) {
            ElMessage.success("Reset all device nodes successfully")
          } else {
            ElMessage.success("Delete device node successfully")
          }
          setTimeout(() => {
            this.refreshTopo(true, false);
          }, 100);
          return true;
        }
      }
      return false;
    },
    deleteTopoRootNode: async function () {
      return this.deleteTopoNode(0);
    },
    refreshTopo: async function (force = false, showMessage = false) {
      const self = this;
      const dasDevices = useDasDevices();
      await dasDevices.updateDeviceInfos(force, showMessage);

    },
    setupTopoData: function (infos) {
      console.log("setup topo data", infos.length)
      try {
        const treeTopoData = getDeviceTreeTopoData(infos);
        if (treeTopoData === undefined || treeTopoData.length < 1) {
          return;
        }
        const root = treeTopoData[0];
        transferTopoNodeForGraph(root)
        this.graphTopoData = root;
        this.treeTopoData = getDeviceTreeTopoData(infos);
      } catch (e) {
        console.log(e);
      }

    },
    changeTopoTypeModel(value) {
      this.topoTypeModel = value
      Cookies.set('topeType', value);
    },
    changeTreeDeviceName(value) {
      this.treeDeviceName = value
      Cookies.set('secondContent', value);
    }
  },
});

function getDeviceGraphTopoData(infos) {
  if (infos === undefined || infos.length < 1) {
    return undefined;
  }
  infos.sort((a, b) => {
    try {
      return a.RouteAddress.localeCompare(b.RouteAddress);
    } catch (e) {
      return 0;
    }
  }); const rootInfo = infos.find(v => String(v.SubID) === "0");
  if (rootInfo == undefined) {
    return undefined;
  }
  const root = {
    ID: String(rootInfo.SubID),
    id: String(rootInfo.SubID),
    info: rootInfo,
    children: [],
  }
  getDeviceGraphNodeChild(root, infos);
  return root;
}

function getDeviceGraphNodeChild(root, infos) {
  if (infos == undefined || infos.length == 0) {
    return;
  }

  for (const key in infos) {
    const info = infos[key];
    if (info.ParentAddress === root.info.RouteAddress && info.CascadingLevel === 1) {
      const child = { ID: String(info.SubID), id: String(info.SubID), info: info, children: [] };
      getDeviceGraphNodeChild(child, infos);
      root.children.push(child);
    } else if (info.ParentAddress === root.info.ParentAddress && info.CascadingLevel === root.info.CascadingLevel + 1) {
      const child = { ID: String(info.SubID), id: String(info.SubID), info: info, children: [] };
      getDeviceGraphNodeChild(child, infos);
      root.children.push(child);
    }
  }
}

function getDeviceTreeTopoData(infos) {
  if (infos === undefined || infos.length < 1) {
    return undefined;
  }
  if (infos.length > 1) {
    infos.sort((a, b) => {
      try {
        return a.RouteAddress.localeCompare(b.RouteAddress);
      } catch (e) {
        return 0;
      }
    });
  }

  const rootInfo = infos.find(v => String(v.SubID) === "0");
  if (rootInfo == undefined) {
    return undefined;
  }
  const root = {
    id: String(rootInfo.SubID),
    info: rootInfo,
    children: [],
  }
  getDeviceTreeNodeChild(root, infos);
  return [root];
}

function getDeviceTreeNodeChild(root, infos) {
  if (infos == undefined || infos.length == 0) {
    return;
  }
  for (const key in infos) {
    const info = infos[key];
    if (info.ParentAddress === root.info.RouteAddress) {
      const child = { id: String(info.SubID), info: info, children: [] };
      getDeviceTreeNodeChild(child, infos);
      root.children.push(child);
    }
  }
  if (root.children.length > 1) {
    root.children.sort((a, b) => {
      try {
        const opa = a?.info?.OpticalPort?.startsWith("OP") ? a.info.OpticalPort.substring(3) : "";
        const opb = b?.info?.OpticalPort?.startsWith("OP") ? b.info.OpticalPort.substring(3) : "";
        const la = a?.info?.CascadingLevel || "0";
        const lb = a?.info?.CascadingLevel || "0";
        const av = opa == "1/2" ? Number(la) : Number(opa) * 1000 + Number(la);
        const bv = opb == "1/2" ? Number(lb) : Number(opb) * 1000 + Number(lb);
        return av - bv;
      } catch (e) {
        console.log(e)
        return 0
      }
    })
  }
}

function transferTopoNodeForGraph(node) {
  if (node.children?.length > 0) {
    const children = node.children.filter(v => v.info?.CascadingLevel > 1 && v.info?.CascadingLevel > node.info.CascadingLevel);
    node.children = node.children.filter(v => !(v.info?.CascadingLevel > 1 && v.info?.CascadingLevel > node.info.CascadingLevel));
    node.children.forEach(child => {
      transferTopoNodeForGraph(child);
    })
    if (children?.length > 0) {
      const sortedChildren = children.sort((a, b) => {
        try {
          return a.info.CascadingLevel - b.info.CascadingLevel
        } catch (e) {
          return 0;
        }
      });
      for (let i = 0; i < sortedChildren.length - 1; i++) {
        const child = sortedChildren[i];
        transferTopoNodeForGraph(child);
        const next = sortedChildren[i + 1];
        if (next != undefined) {
          if (child.children) {
            child.children.push(next);
          } else {
            child.children = [next]
          }
        }
      }
      node.children.push(sortedChildren[0])
    }
  }
}