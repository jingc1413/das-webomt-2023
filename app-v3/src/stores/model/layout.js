export function setupDeviceLayoutModel(sub, deviceModel, deviceParams) {
  const layout = JSON.parse(JSON.stringify(deviceModel.layout));
  setupLayoutElement(layout, deviceParams);
  // console.log("setup layout", sub);
  const deviceLayout = {
    sub: String(sub),
    layout: layout,
  };
  return deviceLayout;
}

function setupLayoutElement(elem, params) {
  if (Object.hasOwn(elem, "setup") && elem.setup) {
    return;
  }
  if (elem.Type == "Table") {
    elem.supportTableRowClick = false;
  }

  if (elem.OID) {
    elem.param = params.getParam(elem.OID);
  }
  if (elem.param) {
    if (elem.Access === "rw") {
      elem.InputRules = [
        {
          required: false,
          validator: (rule, value, callback) => {
            if (elem.InputDisabled) {
              callback();
              return;
            }
            const result = elem.param.validateInputValue(elem.param.InputValue);
            callback(result);
          },
          trigger: "blur" | "change",
        },
      ];
    }
  }

  if (elem.Style?.disableParam != undefined) {
    elem.styleDisableParam = params.getParam(elem.Style?.disableParam);
  }
  if (elem.Style?.enableParam != undefined) {
    elem.styleEnableParam = params.getParam(elem.Style?.enableParam);
  }
  if (elem.styleDisableParam != undefined || elem.styleEnableParam != undefined) {
    Object.defineProperty(elem, "InputDisabled", {
      get: function () {
        const result =
          (this.styleDisableParam &&
            this.Style?.disableValue != undefined &&
            this.Style?.disableValue === this.styleDisableParam.Value) ||
          (this.styleDisableParam &&
            this.Style?.disableInputValue != undefined &&
            this.Style?.disableInputValue === this.styleDisableParam.InputValue) ||
          (this.styleEnableParam &&
            this.Style?.enableValue != undefined &&
            this.Style?.enableValue !== this.styleEnableParam.Value) ||
          (this.styleEnableParam &&
            this.Style?.enableInputValue != undefined &&
            this.Style?.enableInputValue !== this.styleEnableParam.InputValue);
        if (this.param) {
          if (result) {
            this.param.resetInputValue();
          }
          this.param.InputDisabled = result;
        }
        return result;
      },
    });
  }

  elem.Data?.forEach((row, rowIndex) => {
    if (row) {
      Object.keys(row).forEach((key) => {
        setupLayoutElement(row[key], params);
      });
    }
    if (rowIndex == 0 && row._actions) {
      try {
        elem.supportTableRowClick = row._actions?.Actions?.click && true;
      } catch (e) {
        console.log(e);
      }
    }
  });
  if (elem.Actions) {
    Object.keys(elem.Actions).forEach((key) => {
      setupLayoutElement(elem.Actions[key], params);
    });
  }
  elem.Items?.forEach((item) => setupLayoutElement(item, params));
  elem.setup = true;
  return elem;
}
