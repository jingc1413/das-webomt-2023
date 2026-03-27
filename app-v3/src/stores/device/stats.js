import apix from "@/api";
import { ElMessage } from "element-plus";
import { translator as t } from "@/i18n";
import model from "@/stores/model";

export function newDeviceStatsManager(sub, info, params, stats) {
  const m = {
    sub,
    info,
    params,
    stats,
    chartData: {},
    sourceData: [],
  };

  m.getMetrics = function () {
    return this.stats.metrics;
  };
  m.getTabs = function () {
    return this.stats.tabs;
  };
  m.getChartData = function () {
    return this.chartData;
  };

  m.queryStats = async function ({ sub, beginTime, endTime, keys, showMessage = false, isPrintMode = false }) {
    if (isPrintMode && Object.keys(this.chartDate).length > 0) {
      return JSON.parse(JSON.stringify(this.chartDate));
    }
    beginTime = model.nowTime2unixTimestampWithoutTimezones(beginTime);
    endTime = model.nowTime2unixTimestampWithoutTimezones(endTime);

    return new Promise((resolve, reject) => {
      apix
        .queryDeviceStats(sub, beginTime, endTime, [])
        .then((result) => {
          if (showMessage) {
            ElMessage.success("Statistics data retrieved successfully");
          }
          this.sourceData = result;
          let filterData = filterChartData(result);
          this.updateChartData(filterData);
          return resolve(filterData);
        })
        .catch((e) => {
          if (showMessage) {
            ElMessage.error(t("tip.RequestFailed"));
          }
          reject(e);
        });
    });
  };
  m.updateChartData = function (data = {}) {
    let keys = Object.keys(data);
    keys.forEach((key) => {
      this.chartData[key] = data[key];
    });
  };
  return m;
}

function filterChartData(res) {
  // console.log('filterChartData', { res })
  if (!Array.isArray(res)) {
    return {};
  }
  if (res.length == 0) {
    return {};
  }
  res.sort((a, b) => {
    a.t - b.t;
  });
  let dataJson = {};
  let len = res.length;
  let lastTime = res[0]["t"];
  for (let index = 0; index < len; index++) {
    let infoItem = res[index];
    let chartKeys = Object.keys(infoItem);
    let time = model.unixTimestampWithoutTimezones2nowTime(infoItem["t"], false);
    let addNullTime = model.unixTimestampWithoutTimezones2nowTime(infoItem["t"] - 300, false);
    let addNull = false;
    if (infoItem["t"] - lastTime > 600) {
      addNull = true;
    }
    chartKeys.forEach((chartKey) => {
      if (chartKey == "t") {
        return;
      }
      if (!dataJson[chartKey]) {
        dataJson[chartKey] = [];
      }
      if (addNull) {
        dataJson[chartKey].push([addNullTime, null]);
      }
      let lastIndex = dataJson[chartKey].length - 1;
      if (lastIndex >= 0 && dataJson[chartKey][lastIndex][0] == time) {
        dataJson[chartKey].splice(lastIndex, 1, [time, infoItem[chartKey]]);
      } else {
        dataJson[chartKey].push([time, infoItem[chartKey]]);
      }
    });
    lastTime = infoItem["t"];
  }
  // console.log('filterChartData', { dataJson })
  return dataJson;
}
