export function setupDeviceStatsModel(sub, model) {
  const params = model.parameters;
  if (params === undefined) {
    return;
  }
  const systemItems = [
    {
      key: "cpu_usage",
      name: "CPU Usage",
      unit: "%",
      max: 100,
      min: 0,
      span: 12,
      height: "320px",
    },
    {
      key: "memory_usage",
      name: "Memory Usage",
      unit: "%",
      max: 100,
      min: 0,
      span: 12,
      height: "320px",
    },
    {
      key: "disk_usage",
      name: "Storage Usage",
      unit: "%",
      max: 100,
      min: 0,
      span: 12,
      height: "320px",
    },
  ];
  const radioModulePortInputPowerItems = [];
  const radioModuleULInputPowerItems = [];
  const radioModuleDLOutputPowerItems = [];

  const inputModulePortInputPowerItems = [];
  const amplifierModuleULInputPowerItems = [];
  const amplifierModuleDLOutputPowerItems = [];

  params.forEach((param) => {
    if (param.Name === "Element Operating Temperature") {
      systemItems.push({
        key: param.PrivOid,
        name: param.Name,
        unit: param.UnitName,
        min: 0,
        span: 12,
        height: "320px",
      });
    } else if (param.Name?.match(/^Radio Module (\d+) Port (\d+) Input Power$/)) {
      radioModulePortInputPowerItems.push(addParamItem(param, "Input Power"));
    } else if (param.Name?.match(/^Radio Module (\d+) UL Input Power$/)) {
      radioModuleULInputPowerItems.push(addParamItem(param, "UL Input Power"));
    } else if (param.Name?.match(/^Radio Module (\d+) DL Output Power$/)) {
      radioModuleDLOutputPowerItems.push(addParamItem(param, "DL Output Power"));
    } else if (param.Name?.match(/^Input Module (\d+) Port (\d+) Input Power$/)) {
      radioModulePortInputPowerItems.push(addParamItem(param, "Input Power"));
    } else if (param.Name?.match(/^Amplifier Module (\d+) UL Input Power$/)) {
      radioModuleULInputPowerItems.push(addParamItem(param, "UL Input Power"));
    } else if (param.Name?.match(/^Amplifier Module (\d+) DL Output Power$/)) {
      radioModuleDLOutputPowerItems.push(addParamItem(param, "DL Output Power"));
    }
  });

  const auRadioModuleItems = [];
  const ruRadioModuleItems = [];
  const auInputModuleItems = [];
  const ruAmplifierModuleItems = [];

  if (radioModulePortInputPowerItems.length > 0) {
    auRadioModuleItems.push(addModuleItem(radioModulePortInputPowerItems, {
      key: "port_input_power",
      name: "Input Power",
    }))
  }
  if (radioModuleULInputPowerItems.length > 0) {
    ruRadioModuleItems.push(addModuleItem(radioModuleULInputPowerItems, {
      key: "ul_input_power",
      name: "UL Input Power",
    }))
  }
  if (radioModuleDLOutputPowerItems.length > 0) {
    ruRadioModuleItems.push(addModuleItem(radioModuleDLOutputPowerItems, {
      key: "dl_output_power",
      name: "DL Output Power",
    }))
  }
  if (inputModulePortInputPowerItems.length > 0) {
    auInputModuleItems.push(addModuleItem(inputModulePortInputPowerItems, {
      key: "port_input_power",
      name: "Input Power",
    }))
  }
  if (amplifierModuleULInputPowerItems.length > 0) {
    ruAmplifierModuleItems.push(addModuleItem(amplifierModuleULInputPowerItems, {
      key: "ul_input_power",
      name: "UL Input Power",
    }))
  }
  if (amplifierModuleDLOutputPowerItems.length > 0) {
    ruAmplifierModuleItems.push(addModuleItem(amplifierModuleDLOutputPowerItems, {
      key: "dl_output_power",
      name: "DL Output Power",
    }))
  }

  const metrics = {};
  const tabs = {};
  let comparisonKeys = [];
  if (systemItems.length > 0) {
    let systemTabKey = [];
    systemItems.forEach((systemItem) => {
      metrics[systemItem.key] = systemItem;
      systemTabKey.push(systemItem.key);
    });
    tabs["system"] = {
      key: "system",
      name: "System",
      chartKeys: systemTabKey,
    };
    comparisonKeys.push(...systemTabKey);
  }

  if (auRadioModuleItems.length > 0) {
    let keys = [];
    auRadioModuleItems.forEach((v) => {
      metrics[v.key] = v;
      keys.push(v.key);
    });
    tabs["radio_module"] = {
      key: "radio_module",
      name: "Radio Module",
      chartKeys: keys,
    };
    comparisonKeys.push(...keys);
  } else if (ruRadioModuleItems.length > 0) {
    let keys = [];
    ruRadioModuleItems.forEach((v) => {
      metrics[v.key] = v;
      keys.push(v.key);
    });
    tabs["radio_module"] = {
      key: "radio_module",
      name: "Radio Module",
      chartKeys: keys,
    };
    comparisonKeys.push(...keys);
  }

  if (auInputModuleItems.length > 0) {
    let keys = [];
    auInputModuleItems.forEach((v) => {
      metrics[v.key] = v;
      keys.push(v.key);
    });
    tabs["input_module"] = {
      key: "input_module",
      name: "Input Module",
      chartKeys: keys,
    };
    comparisonKeys.push(...keys);
  } else if (ruAmplifierModuleItems.length > 0) {
    let keys = [];
    ruAmplifierModuleItems.forEach((v) => {
      metrics[v.key] = v;
      keys.push(v.key);
    });
    tabs["amplifier_module"] = {
      key: "amplifier_module",
      name: "Amplifier Module",
      chartKeys: keys,
    };
    comparisonKeys.push(...keys);
  }

  if (Object.keys(metrics).length > 0) {
    tabs["comparison"] = {
      key: "comparison",
      name: "Comparison",
      chartKeys: comparisonKeys,
    };
  }
  return {
    tabs,
    metrics,
  }
}

function addParamItem(param, replaceString) {
  return {
    key: param.PrivOid,
    name: param.Name.replaceAll(replaceString, ""),
    unit: param.UnitName,
    max: param.Ratio !== undefined ? param.Max / param.Ratio : param.Max,
    min: param.Ratio !== undefined ? param.Min / param.Ratio : param.Min,
  }
}

function addModuleItem(powerItems, options = {}) {
  let {key, name, yInterval=20, span=24} = options
  const item0 = powerItems[0];
    let min = -100;
    let max = 20;
    if (item0?.max && item0.max > max) {
      max = Math.ceil(item0.max / 20) * 20;
    }
    if (item0?.min && item0.min < min) {
      min = Math.ceil(item0.min / 20) * 20;
    }
    let height = 400 + Math.ceil(powerItems.length / 8) * 40;
    let bottom = 64 + Math.ceil(powerItems.length / 8) * 40;
    return {
      key: key,
      name: name,
      items: powerItems,
      unit: item0?.unit,
      max: max,
      min: min,
      yInterval: yInterval,
      span: span,
      height: `${height}px`,
      bottom: bottom,
    };
}
