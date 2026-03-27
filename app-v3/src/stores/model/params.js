import { dayjs } from "element-plus";
import * as utils from "./utils.js";

var utc = require("dayjs/plugin/utc");
dayjs.extend(utc);

export function setupDeviceParamsModel(sub, deviceModel) {
  const params = JSON.parse(JSON.stringify(deviceModel.parameters));
  params.forEach((item) => setupParam(item));

  const deviceParams = {
    sub: String(sub),
    params: params,
  };
  deviceParams.getParamByName = function (name) {
    const param = this.params.find((v) => v.Name === name);
    if (param) {
      setupParam(param);
    }
    return param;
  };
  deviceParams.getParam = function (oid) {
    const match = oid?.match(/(.*)\[(.*?)\]/);

    if (match) {
      const oid2 = match[1];
      const param = this.params.find((v) => v.PrivOid === oid2);

      if (param) {
        const childIndex = Number(match[2]);
        const child = param.Child ? param.Child[childIndex] : undefined;
        if (child) {
          setupParam(child);
        }
        return child;
      } else {
        return undefined;
      }
    }
    const param = this.params.find((v) => v.PrivOid === oid);
    if (param) {
      setupParam(param);
    }
    return param;
  };
  deviceParams.getParamByOldID = function (oldPath, oldId) {
    if (oldPath?.length > 0) {
      const tmp = this.params.find((p) => {
        let match = false;
        for (const i in p.Groups) {
          if (p.Groups[i] === oldPath) {
            match = true;
            break;
          }
        }
        return match;
      });
      if (tmp) {
        const oidPrefix = tmp.PrivOid.split(".")[0];
        const oid = oidPrefix + ".P" + oldId.toUpperCase();
        return this.getParam(oid);
      }
      return undefined;
    }
  };
  deviceParams.setParamValueUpdateCallback = function (oid, cb) {
    const param = this.getParam(oid);
    param?.setValueUpdateCallback(cb);
  };
  return deviceParams;
}

function setupParam(param) {
  if (Object.hasOwn(param, "setup") && param.setup) {
    return;
  }
  param.Readable = param.Access === "rw" || param.Access === "ro";
  param.Writable = param.Access === "rw" || param.Access === "wo";

  param.IsDateTimeDataType = param.DataType === "datetime";
  param.IsIPv4DataType = param.DataType === "ipv4";
  param.IsNumberDataType = checkNumberDataType(param.DataType);
  param.IsBinaryDataType = param.DataType === "binary";
  param.IsStructDataType = param.DataType == "object";
  param.IsArray = checkArrayDataType(param);
  if (param.Options) {
    const sortOptions = [];
    if (param.IsNumberDataType) {
      const optionKeys = Object.keys(param.Options).sort((a, b) => (Number(a) > Number(b) ? 1 : -1));
      optionKeys.forEach((k) => {
        sortOptions.push({ k: Number(k), v: param.Options[k] });
      });
    } else {
      const optionKeys = Object.keys(param.Options).sort();
      optionKeys.forEach((k) => {
        if (param.MultipleOption && k === "0000") {
          return;
        }
        sortOptions.push({ k: k, v: param.Options[k] });
      });
    }
    param.SortOptions = sortOptions;

    param.findOption = function (v) {
      if (this.IsNumberDataType) {
        const option = this.SortOptions.find((item) => Number(item.k) === Number(v));
        return option;
      } else {
        const option = this.SortOptions.find((item) => item.k === v);
        return option;
      }
    };
  }

  param.Value = undefined;
  if (param.Access == "rw" || param.Access == "ro") {
    if (param.PrivOid === "TB0.P0AFF") {
      param.Value = "A5";
    } else if (param.MultipleOption) {
      param.Value = [];
    } else if (param.SortOptions) {
      param.Value = param.SortOptions[0].k;
    } else if (param.IsArray) {
      param.Value = [];
    } else if (param.IsStructDataType) {
      param.Value = [];
    } else if (param.IsNumberDataType) {
      param.NumberStep = param.Ratio ? 1 / param.Ratio : 1;
      if (param.Max !== undefined) {
        param.NumberMax = param.Ratio !== undefined ? param.Max / param.Ratio : param.Max;
      }
      if (param.Min !== undefined) {
        param.NumberMin = param.Ratio !== undefined ? param.Min / param.Ratio : param.Min;
      }
      param.Value = Number(0);
    } else if (param.IsDateTimeDataType) {
      // ignore
    } else {
      if (param.Max !== undefined && param.Max <= param.ByteSize) {
        param.TextMax = param.IsBinaryDataType ? param.Max * 2 : param.Max;
      } else {
        param.TextMax = param.IsBinaryDataType ? param.ByteSize * 2 : param.ByteSize;
      }
      if (param.Min !== undefined) {
        param.TextMin = param.IsBinaryDataType ? param.Min * 2 : param.Min;
      } else {
        param.TextMin = 0;
      }
      param.Value = "";
    }
  }

  if (param.Child && param.Child.length > 0) {
    const parent = param;
    param.Child.forEach((item) => {
      setupParam(item);
      item.onUpdateParentValue = function () {
        const parentInputValue = [];
        parent.Child.forEach((child) => {
          parentInputValue.push(child.InputValue);
        });
        parent.InputValue = parentInputValue;
      };
    });
  }

  param.formatRawFromValue = function (value) {
    if (value == undefined) {
      return undefined;
    }
    let out = value;
    if (this.PrivOid === "T02.P0150") {
      out = utils.dayjsToDeviceTimestamp(dayjs(value));
    } else if (this.IsStructDataType) {
      out = [];
      for (const i in value) {
        const v = value[i];
        const child = this.Child[i];
        out.push(child.formatRawFromValue(v));
      }
    } else if (this.IsNumberDataType) {
      out = Number(value);
      if (this.Ratio !== undefined && this.Ratio > 1) {
        out = out * this.Ratio;
      }
    } else if (this.IsBinaryDataType) {
      if (this.MultipleOption) {
        out = [];
        value.forEach((v) => {
          out.push(v.toUpperCase());
        });
      } else {
        out = value.toUpperCase();
      }
    }
    return out;
  };

  param.formatValueFromRaw = function (value) {
    if (value == undefined) {
      return undefined;
    }
    let out = value;
    if (this.PrivOid === "T02.P0150") {
      out = utils.deviceTimestampToDayJs(value).format("YYYY-MM-DD HH:mm:ss");
    } else if (this.IsStructDataType) {
      out = [];
      for (const i in value) {
        const v = value[i];
        const child = this.Child[i];
        out.push(child.formatValueFromRaw(v));
      }
    } else if (this.IsNumberDataType) {
      out = Number(value);
      if (this.Ratio !== undefined && this.Ratio > 1) {
        out = out / this.Ratio;
      }
    } else if (this.IsBinaryDataType) {
      if (this.MultipleOption) {
        out = [];
        value.forEach((v) => {
          out.push(v.toUpperCase());
        });
      } else {
        out = value.toUpperCase();
      }
    }

    return out;
  };

  param.inputValueStore = param.Value;
  Object.defineProperty(param, "InputValue", {
    configurable: true,
    enumerable: true,
    get: function () {
      return this.inputValueStore;
    },
    set: function (v) {
      const last = this.inputValueStore;
      if (this.PrivOid === "T02.P0101") {
        this.inputValueStore = v?.replace(/[^0-9A-Fa-f]/g, "")?.toUpperCase();
      } else if (this.PrivOid === "T02.P0F3E") {
        this.inputValueStore = v?.replace(/[^SUDsud]/g, "")?.toUpperCase();
      } else if (this.MultipleOption) {
        this.inputValueStore = v;
      } else if (this.DataType === "binary") {
        this.inputValueStore = v?.replace(/[^0-9A-Fa-f]/g, "")?.toUpperCase();
      } else {
        this.inputValueStore = v;
      }
      if (last !== this.inputValueStore) {
        if (this.onUpdateParentValue) {
          this.onUpdateParentValue();
        }
      }
    },
  });

  param.resetInputValue = function () {
    this.InputValue = this.Value;
    this.RespCode = "00";
    this.RespMsg = "";
  };
  param.validateInputValue = function (value) {
    let v = value;
    if (this.IsIPv4DataType) {
      //ignore
    } else if (this.IsDateTimeDataType) {
      if (!dayjs(v, "YYYY-MM-DD HH:mm:ss").isValid()) {
        return "Invalid date";
      }
    } else if (this.SortOptions) {
      if (this.MultipleOption) {
        for (const key in v) {
          if (Object.hasOwnProperty.call(value, key)) {
            const option = this.findOption(value[key]);
            if (option === undefined) {
              return `Invalid option value`;
            }
          }
        }
      } else {
        const option = this.findOption(v);
        if (option === undefined) {
          return `Invalid option value`;
        }
      }
    } else if (this.IsBinaryDataType) {
      if (this.TextMin != undefined && v?.length < this.TextMin) {
        return `Minimum of ${this.TextMin} characters in length`;
      } else if (this.TextMax != undefined && v?.length > this.TextMax) {
        return `Maximum of ${this.TextMax} characters in length`;
      }
    } else if (this.IsNumberDataType) {
      if (this.NumberMin != undefined && v < this.NumberMin) {
        return `Minimum of ${this.NumberMin}`;
      } else if (this.NumberMax != undefined && v > this.NumberMax) {
        return `Maximum of ${this.NumberMax}`;
      }
    } else if (this.TextMin != undefined && v?.length < this.TextMin) {
      return `Minimum of ${this.TextMin} characters in length`;
    } else if (this.TextMax != undefined && v?.length > this.TextMax) {
      return `Maximum of ${this.TextMax} characters in length`;
    }
    return undefined;
  };

  param.setRawValue = function (value, code = "00") {
    this.RespCode = code;
    this.RespMsg = getCodeMessage(this, code);
    this.RawValue = value;
    this.Value = this.formatValueFromRaw(this.RawValue);
    this.InputValue = this.Value;

    if (this.IsArray) {
      for (const i in this.RawValue) {
        let v = this.RawValue[i];
        const child = this.Child[i];
        if (child) {
          child.setRawValue(v, code);
        }
      }
    } else if (this.IsStructDataType) {
      for (const i in this.RawValue) {
        const v = this.RawValue[i];
        const child = this.Child[i];
        if (child) {
          child.setRawValue(v, code);
        }
      }
    }

    try {
      if (this.onValueUpdate) {
        this.onValueUpdate({
          raw: this.RawValue,
          value: this.Value,
          respCode: this.RespCode,
          respMsg: this.RespMsg,
        });
      }
    } catch (e) {
      console.error(e);
    }
    if (this.debug) {
      console.log("update value", this.PrivOid, this.Value);
    }
    return this.Value;
  };
  param.setValueUpdateCallback = function (cb) {
    this.onValueUpdate = cb;
  };

  param.getValue = function ({ defaultValue = undefined }) {
    if (defaultValue !== undefined) {
      return defaultValue;
    }
    if (this.Access === "wo") {
      return undefined;
    }
    if (this.IsNumberDataType) {
      if (this.Max && this.RawValue > this.Max) {
        return "++";
      }
      if (this.Min && this.RawValue < this.Min) {
        return "--";
      }
    }
    return this.Value;
  };

  param.getShowValue = function ({ withUnit = false, defaultValue = undefined }) {
    let v = this.getValue({ defaultValue });
    if (v === undefined || v === null) {
      return undefined;
    }
    // if (v === "invalid" || v=== "null") {
    //   return "invalid";
    // }
    if (this.MultipleOption) {
      try {
        if (v === undefined || v.length == 0) {
          return "NULL";
        }
        const out = [];
        v.forEach((v2) => {
          const option = this.findOption(v2);
          if (option != undefined) {
            out.push(option.v);
          }
        });
        if (out.length > 0) {
          return out.join(",");
        } else {
          return "NULL";
        }
      } catch (e) {
        console.error(e);
        return "";
      }
    } else if (this.SortOptions && v != undefined) {
      const option = this.findOption(v);
      return option?.v ?? v;
    }
    if (withUnit && this.UnitName && this.UnitName !== "") {
      return v + " " + this.UnitName;
    }
    return v;
  };

  param.getShowStyle = function () {
    let out = "";
    if (this.RespCode !== undefined && this.RespCode !== "00") {
      out = out + "color:red;";
    }
    return out;
  };

  param.getSwitchData = function (layoutStyle) {
    if (layoutStyle?.activeValue && layoutStyle?.inactiveValue) {
      const out = {
        activeValue: layoutStyle?.activeValue,
        activeText: layoutStyle?.activeText,
        activeType: layoutStyle?.activeType,
        inactiveValue: layoutStyle?.inactiveValue,
        inactiveText: layoutStyle?.inactiveText,
        inactiveType: layoutStyle?.inactiveType,
      };
      if (this.IsNumberDataType) {
        out.activeValue = Number(out.activeValue);
        out.inactiveValue = Number(out.inactiveValue);
      }
      return out;
    }
    if (this.SortOptions && this.SortOptions.length >= 2) {
      const activeOption = param.SortOptions[1];
      const inactiveOption = param.SortOptions[0];

      return {
        activeValue: activeOption.k,
        activeText: activeOption.v,
        inactiveValue: inactiveOption.k,
        inactiveText: inactiveOption.v,
      };
    }
    return undefined;
  };

  param.getTreeSelectData = function () {
    const out = [];
    if (this.SortOptions && this.SortOptions.length > 0) {
      this.SortOptions.forEach((opt) => {
        const args = opt.v.split("_");
        let node = undefined;
        let name = "";
        let num = 0;
        if (this.IsNumberDataType) {
          num = opt.k;
        } else {
          num = Number("0x" + opt.k);
        }

        args.forEach((arg, index) => {
          if (name == "") {
            name = arg;
          } else {
            name = name + " " + arg;
          }
          if (index == args.length - 1) {
            node.children.push({
              key: arg,
              label: name,
              value: opt.k,
              number: num,
            });
            return;
          }
          let node2 = undefined;
          if (index == 0) {
            node2 = out.find((item) => item.key === arg);
            if (node2 === undefined) {
              node2 = {
                key: arg,
                label: name,
                value: opt.k,
                number: num,
                children: [],
              };
              out.push(node2);
            } else {
              node2.number = node2.number + num;
              node2.value = node2.number.toString(16).padStart(this.ByteSize * 2, "0");
            }
          } else {
            node2 = node.children.find((item) => item.key === arg);
            if (node2 === undefined) {
              node2 = {
                key: arg,
                label: name,
                value: opt.k,
                number: num,
                children: [],
              };
              node.children.push(node2);
            } else {
              node2.number = node2.number + num;
              node2.value = node2.number.toString(16).padStart(this.ByteSize * 2, "0");
            }
          }
          node = node2;
        });
      });
    }
    return out;
  };

  param.setup = true;
  return param;
}

function checkArrayDataType(param) {
  const parts = param.DataType.split(":");
  if (parts.length === 3 && parts[1] === "array") {
    const out = { DataType: parts[0], Size: Number(parts[2]) };
    return out;
  }
  return undefined;
}

function checkNumberDataType(dataType) {
  if (
    dataType === "int8" ||
    dataType === "int16" ||
    dataType === "int32" ||
    dataType === "int64" ||
    dataType === "uint8" ||
    dataType === "uint16" ||
    dataType === "uint32" ||
    dataType === "uint64"
  ) {
    return true;
  }
  return false;
}

function getCodeMessage(def, code) {
  const msg = { warning: undefined, error: undefined, info: undefined };
  switch (code) {
    case "00":
      break;
    case "01":
      msg.error = "Invalid parameter";
      break;
    case "02":
      if ((def.IsNumberDataType && def.NumberMin) || def.NumberMax) {
        msg.error = "Value out of range" + ` ${def.NumberMin} ~ ${def.NumberMax}`;
        // } else if (def.TextMin || def.TextMax) {
        //   msg.error = 'Value out of range' + ` ${def.TextMin} ~ ${def.TextMax}`
      } else {
        msg.error = "Invalid value";
      }
      break;
    case "03":
      msg.error = "Illegal value";
      break;
    case "04":
      msg.error = "Value length incorrect";
      break;
    case "05":
      msg.warning = "Value is low";
      break;
    case "06":
      msg.warning = "Value is high";
      break;
    case "07":
      msg.warning = "Value cannot be detected";
      break;
    case "09":
      msg.error = "Other error";
      break;
    case "0a":
    case "0A":
      msg.error = "Not in factory mode";
      break;
    default:
      msg.warning = "Error " + code;
      break;
  }
  return msg;
}
