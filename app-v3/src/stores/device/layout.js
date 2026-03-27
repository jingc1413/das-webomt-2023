import settings from "@/settings.js";
import { useDasDevices } from "@/stores/das-devices";
import permissions from "@/plugins/permissions.js";

const filterTreePages = ["Logs", "Upgrade", "Configuration", "Account", "Address Interface"];
export function newDeviceLayoutManager(sub, info, params, layoutModel) {
  const m = {
    sub,
    info,
    params,
    layoutModel,

    appLayout: {},
    originPageTreeData: [],
    originPageTreePaths: [],
    filterTreePages: filterTreePages,
  };

  m.pageTreeData = function () {
    return this.originPageTreeData;
  };
  m.getModule = function (moduleKey) {
    return this.appLayout.Items.find((v) => v.Key === moduleKey);
  };
  m.getPage = function (moduleKey, pageKey) {
    const module = this.getModule(moduleKey);
    if (module != undefined) {
      return module.Items.find((v) => v.Key === pageKey);
    }
    return undefined;
  };

  m.setup = function () {
    this.appLayout = getLayoutPermission(this.layoutModel.layout);
    try {
      setupLayoutAccess(this.appLayout);
    } catch (e) {
      console.log(e);
    }

    // this.appLayout = getLayoutPermission(this.layoutModel.layout, this.params);

    this.isMenuCollapsed = false;
    this.currentModule = null;
    this.currentPage = null;
    const pageTreeData = [];
    const pageTreePaths = [];

    let superModeOnly = false;
    this.appLayout.Items.forEach((module) => {
      module.superModeOnly = superModeOnly;
      if (module.Name === "Maintenance") {
        superModeOnly = true;
      }
    });

    this.appLayout.Items.forEach((module) => {
      const moduleChild = [];
      module.Items.forEach((page) => {
        page.tabs = page?.Items?.find((item) => item.Type === "Layout:Tabs");
        getLayoutParameterOidsAndDefaultValues(page, this.params);
        if (this.filterTreePages.includes(page.Name)) {
          return;
        }
        if (page.tabs) {
          const pageChild = [];
          page.tabs.Items.forEach((tabItem) => {
            if (tabItem.Type === "Page") {
              pageTreePaths.push([module.Key, page.Key, tabItem.Key]);
              pageChild.push({
                key: "/" + module.Key + "/" + page.Key + "/" + tabItem.Key,
                label: tabItem.Name,
                fullLabel: module.Name + " / " + page.Name + " / " + tabItem.Name,
                page: tabItem,
                superModeOnly: module.superModeOnly,
              });
            }
          });
          moduleChild.push({
            key: "/" + module.Key + "/" + page.Key,
            label: page.Name,
            fullLabel: module.Name + " / " + page.Name,
            children: pageChild,
            superModeOnly: module.superModeOnly,
          });
        } else {
          pageTreePaths.push([module.Key, page.Key]);
          moduleChild.push({
            key: "/" + module.Key + "/" + page.Key,
            label: page.Name,
            fullLabel: module.Name + " / " + page.Name,
            page: page,
            superModeOnly: module.superModeOnly,
          });
        }
      });
      pageTreeData.push({
        key: "/" + module.Key,
        label: module.Name,
        fullLabel: module.Name,
        children: moduleChild,
        superModeOnly: module.superModeOnly,
      });
    });
    this.originPageTreePaths = pageTreePaths;
    this.originPageTreeData = [
      {
        key: "/",
        label: "All",
        fullLabel: "All",
        children: pageTreeData,
      },
    ];

    this.updateTitle();
    console.log("setup layout", this.appLayout);
  };
  m.updateTitle = function () {
    let title = this.info.DeviceTypeName + " (" + settings.appInfo.version + "-" + settings.appInfo.build + ")";
    if (this.info.OpticalPort && this.info.OpticalPort !== "") {
      title = `[${this.info.OpticalPort}${this.info?.OpticalInputPort && ('>>'+this.info.OpticalInputPort)}${this.info.CascadingLevel ? ":L" + this.info.CascadingLevel : ""}] ` + title;
    }
    let title2 = "";
    if (this.info.DeviceName) {
      title2 = this.info.DeviceName;
    }
    if (this.info.DeviceTypeName) {
      title2 = this.info.DeviceTypeName + ":" + title2;
    }
    this.currentWebTitle = title;
    this.currentDeviceNameTitle = title2;
    // document.title = this.currentWebTitle;
  };
  m.matchPageTreePath = function (path) {
    const tmp = this.originPageTreePaths.find((v) => path.startsWith(v.join(".")));
    if (tmp) {
      return JSON.parse(JSON.stringify(tmp));
    }
  };
  m.collapseMenu = function ({ collapsed }) {
    this.isMenuCollapsed = collapsed;
  };
  m.selectPage = function (moduleKey, pageKey) {
    console.log("select page", this.sub, moduleKey, pageKey);
    const module = this.appLayout.Items.find((v) => v.Key === moduleKey);
    if (module) {
      let page = module.Items.find((v) => v.Key === pageKey);
      if (page) {
        this.currentModule = module;
        this.currentPage = page;
        return true;
      }
    }
    return false;
  };
  m.openViewPage = function ({ page, viewPage, onCloseViewPage }) {
    if (page.viewPage?.Key == viewPage.Key) {
      return;
    }
    page.viewPage = viewPage;
    page.onCloseViewPage = onCloseViewPage;
  };
  m.closeViewPage = function ({ page }) {
    page.viewPage = null;
    if (page.onCloseViewPage) {
      page.onCloseViewPage();
    }
  };

  m.setup();
  return m;
}

function parseOid(oid) {
  let out = oid;
  const match = oid.match(/(.*)\[(.*?)\]/);
  if (match) {
    out = match[1];
  }
  return out;
}

function getLayoutParameterOidsAndDefaultValues(layout, params) {
  const oids = [];
  const rOids = [];
  const rwOids = [];
  const woOids = [];

  const defaultValues = [];
  if (layout.Actions) {
    const actionKeys = Object.keys(layout.Actions);
    actionKeys?.forEach((key) => {
      const item = layout.Actions[key];
      const result = getLayoutParameterOidsAndDefaultValues(item, params);
      result.oids.forEach((oid) => {
        if (!oids.includes(oid)) {
          oids.push(oid);
        }
      });
      result.rOids.forEach((oid) => {
        if (!rOids.includes(oid)) {
          rOids.push(oid);
        }
      });
      result.woOids.forEach((oid) => {
        if (!woOids.includes(oid)) {
          woOids.push(oid);
        }
      });
      result.rwOids.forEach((oid) => {
        if (!rwOids.includes(oid)) {
          rwOids.push(oid);
        }
      });
      result.defaultValues.forEach((value) => {
        const value2 = defaultValues.find((v) => {
          v.oid === value.oid;
        });
        if (!value2) {
          defaultValues.push(value);
        }
      });
    });
  }
  const ignoreOids = ["TB2.P0CCC"];
  switch (layout.Type) {
    case "Param": {
      const oid = parseOid(layout.OID);
      if (ignoreOids.includes(oid)) {
        break;
      }
      if (!oids.includes(oid)) {
        oids.push(oid);
      }
      const param = params.getParam(oid);
      if (param) {
        if (layout.Access === "rw" || layout.Access === "ro") {
          if (!rOids.includes(oid)) {
            rOids.push(oid);
          }
        }
        if (layout.Access === "rw") {
          if (!rwOids.includes(oid)) {
            rwOids.push(oid);
          }
        }
        if (layout.Access === "wo" && layout.Value !== undefined && layout.Style?.willReboot !== true) {
          if (!woOids.includes(oid)) {
            woOids.push(oid);
            defaultValues.push({ oid: oid, value: layout.Value });
          }
        }
      }
      break;
    }
    case "Table": {
      const keys = [];
      layout.Items.forEach((item) => {
        keys.push(item.Key);
      });
      layout.Data?.forEach((row) => {
        keys.forEach((key) => {
          const item = row[key];
          if (item) {
            const result = getLayoutParameterOidsAndDefaultValues(item, params);
            result.oids.forEach((oid) => {
              if (!oids.includes(oid)) {
                oids.push(oid);
              }
            });
            result.rOids.forEach((oid) => {
              if (!rOids.includes(oid)) {
                rOids.push(oid);
              }
            });
            result.rwOids.forEach((oid) => {
              if (!rwOids.includes(oid)) {
                rwOids.push(oid);
              }
            });
            result.woOids.forEach((oid) => {
              if (!woOids.includes(oid)) {
                woOids.push(oid);
              }
            });
            result.defaultValues.forEach((value) => {
              const value2 = defaultValues.find((v) => {
                v.oid === value.oid;
              });
              if (!value2) {
                defaultValues.push(value);
              }
            });
          }
        });
      });
      break;
    }
    default: {
      layout.Items?.forEach((item) => {
        const result = getLayoutParameterOidsAndDefaultValues(item, params);
        result.oids.forEach((oid) => {
          if (!oids.includes(oid)) {
            oids.push(oid);
          }
        });
        result.rOids.forEach((oid) => {
          if (!rOids.includes(oid)) {
            rOids.push(oid);
          }
        });
        result.rwOids.forEach((oid) => {
          if (!rwOids.includes(oid)) {
            rwOids.push(oid);
          }
        });
        result.woOids.forEach((oid) => {
          if (!woOids.includes(oid)) {
            woOids.push(oid);
          }
        });
        result.defaultValues.forEach((value) => {
          const value2 = defaultValues.find((v) => {
            v.oid === value.oid;
          });
          if (!value2) {
            defaultValues.push(value);
          }
        });
      });
      break;
    }
  }
  const types = ["Page", "Form", "Table"];
  if (types.includes(layout.Type)) {
    layout.oids = oids;
    layout.rOids = rOids;
    layout.rwOids = rwOids;
    layout.woOids = woOids;

    layout.defaultValues = defaultValues;
  }
  return { oids, rOids, rwOids, woOids, defaultValues };
}

const allAccessKeys = ["get", "set"];
function setupLayoutAccess(app) {
  app.Items?.forEach((module) => {
    if (module.Type == "Module") {
      const moduleAccessKeys = {};
      allAccessKeys.forEach((key) => {
        if (permissions.hasPermission(`page.${module.Key}.${key}`)) {
          moduleAccessKeys[key]=true;
        }
      });
      setElementAccess(module, moduleAccessKeys);
      module.Items?.forEach((page) => {
        if (page.Type == "Page") {
          const pageAccessKeys = {};
          allAccessKeys.forEach((key) => {
            if (permissions.hasPermission(`page.${module.Key}.${page.Key}.${key}`)) {
              pageAccessKeys[key]=true;
            }
          });
          setElementAccess(page, pageAccessKeys);
          page.Items?.forEach((tabs) => {
            if (tabs.Type === "Layout:Tabs") {
              tabs.Items?.forEach((tab) => {
                if (tab.Type == "Page") {
                  const tabAccessKeys = {};

                  allAccessKeys.forEach((key) => {
                    if (permissions.hasPermission(`page.${module.Key}.${page.Key}.${tab.Key}.${key}`)) {
                      tabAccessKeys[key]=true;
                    }
                  });
                  setElementAccess(tab, tabAccessKeys);
                }
              });
            }
          });
        }
      });
    }
  });
}

function setElementAccess(elem, accessKeys = {}) {
  if (!elem) {
    return;
  }
  if (Object.keys(accessKeys).length == 0) {
    return;
  }
  if (!elem.accessKeys) {
    elem.accessKeys = {};
  }

  Object.keys(accessKeys).forEach((key) => {
    elem.accessKeys[key]=accessKeys[key];
  });

  setElementAccess(elem.Actions, accessKeys);
  elem.Items?.forEach((item) => {
    setElementAccess(item, accessKeys);
  });
  elem.Data?.forEach((row) => {
    for (const key in row) {
      if (Object.hasOwnProperty.call(row, key) && key != "accessKeys") {
        setElementAccess(row[key], accessKeys);
      }
    }
  });

  if (elem.Actions) {
    Object.keys(elem.Actions).forEach((key) => {
      if (elem.Actions[key] && key != "accessKeys") {
        setElementAccess(elem.Actions[key], accessKeys);
      }
    });
  }
}


const filterApiPages = {
  "system_settings.account.users": "api.iam.users.list"
};
function getLayoutPermission(layout) {
  layout.Items = layout.Items.filter((module) => {
    return permissions.hasPermission(`page.${module.Key}`)
  });

  layout.Items = layout.Items.map((module) => {
    module.Items = module.Items.filter((page) => {
      page.tabs = page?.Items?.find((item) => item.Type === "Layout:Tabs");
      if (page.tabs) {
        page.tabs.Items = page.tabs.Items.filter((tabItem) => {
          let keyString = `${module.Key}.${page.Key}.${tabItem.Key}`;
          if (tabItem.Type === "Page") {
            if (filterApiPages[keyString]) {
              return permissions.hasPermission(`page.${keyString}.get`) && permissions.hasPermission(filterApiPages[keyString]);
            }
            return permissions.hasPermission(`page.${keyString}.get`)
          }
          return true;
        });
        return page.tabs.Items.length > 0;
      } else {
        let keyString = `${module.Key}.${page.Key}`;
        if (filterApiPages[keyString]) {
          return permissions.hasPermission(`page.${keyString}.get`) && permissions.hasPermission(filterApiPages[keyString]);
        }
        return permissions.hasPermission(`page.${keyString}.get`)
      }
    });
    return module;
  });

  layout.Items = layout.Items.filter((module) => {
    return module.Items.length > 0;
  });

  return layout;
}
