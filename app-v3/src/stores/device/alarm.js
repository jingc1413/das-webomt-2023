import apix from "@/api";
import { ElMessage } from "element-plus";
import { translator as t } from "@/i18n";

export function newDeviceAlarms(sub, params) {
  const m = {
    sub,
    params,
    alarmLogs: [],
    alarmSeverity: [],
    sort: {
      prop: "EventTime",
      order: "descending"
    },
    queryAlarmLogs: [],
    query: {},
    showAlarmLogs: [],
  };

  m.getAlarmLogs = async function(showMessage = false) {
    await apix
    .getDeviceAlarmLogs(this.sub)
    .then((logs) => {
      this.alarmLogs = logs;
      let alarmSeverity = new Set();
      logs.forEach(item=>{
        alarmSeverity.add(item.AlarmSeverity)
      })
      this.alarmSeverity = Array.from(alarmSeverity);
    })
    .catch((e) => {
      console.log(e);
      if (showMessage) {
        ElMessage.error(t("tip.RequestFailed"));
      }
    });
  }

  m.alarmLogsStartSort = function() {
    const key = this.sort.prop;
    const order = this.sort.order;
    this.alarmLogs.sort((a, b) => {
      return compareValue(a[key], b[key], order === 'ascending');
    })
  }

  function compareValue(v1, v2, isAscending) {
    let flag = isAscending == true?1:-1;
    return flag*(String(v1).toUpperCase().localeCompare(String(v2).toUpperCase()))
  }

  m.alarmLogsStartFilter = function() {
    this.queryAlarmLogs = this.alarmLogs.filter(item=>{
      return Object.keys(this.query).some(keyItem=>{
        try {
          return this.query[keyItem] && (String(item[keyItem])?.indexOf(this.query[keyItem]) == -1) // find unmatched
        } catch (error) {
          console.log('getAlarmLogList alarmLogsStartFilter', {error});
        }
        return false;
      }) == false;
    })
  };

  m.queryIsChange = function(newQuery=[]) {
    let tempQuery = {};
    let flag = false;
    newQuery.forEach(item=>{
      let key = item.key;
      tempQuery[key] = item.value;
      if (String(item.value) != String(this.query[key])) {
        flag = true;
      }
    })

    if (flag) {
      this.query = tempQuery;
    }

    return flag;
  }

  m.getAlarmLogList = async function({pageNumber=1, pageSize=50, order, prop, query}, forcedRefresh = false) {
    if (forcedRefresh || this.alarmLogs.length === 0) {
      await this.getAlarmLogs();
    }

    if (order && prop) {
      if (order !== this.sort.order || prop !== this.sort.prop) {
        this.sort = {order, prop};
        this.alarmLogsStartSort();
      }
    }

    if (query && this.queryIsChange(query)) {
      this.alarmLogsStartFilter()
    }
    let offset = (pageNumber-1)*pageSize;
    this.showAlarmLogs = this.queryAlarmLogs.slice(offset, offset + pageSize);
  }

  return m;
}