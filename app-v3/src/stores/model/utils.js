import { dayjs } from "element-plus";
var utc = require("dayjs/plugin/utc");
dayjs.extend(utc);

export function deviceTimestampToDayJs(ts) {
  const text = dayjs.unix(ts).utc().format("YYYY-MM-DD HH:mm:ss");
  return dayjs(text);
}

export function dayjsToDeviceTimestamp(t) {
  const text = t.format("YYYY-MM-DD HH:mm:ss");
  return dayjs.utc(text).utc().unix();
}

export function deviceTimeStringToDayJs(text) {
  const text2 = dayjs.utc(text).utc().format("YYYY-MM-DD HH:mm:ss");
  return dayjs(text2);
}

export function dayjsToDeviceTimeString(t) {
  return t.local().format("YYYY-MM-DD HH:mm:ss");
}

/**
 * nowTime2unixTimestampWithoutTimezones
 * @param {number} nowTime timestamp s
 * @param {boolean} isReturnUnix
 */
export function nowTime2unixTimestampWithoutTimezones(nowTime) {
  return dayjsToDeviceTimestamp(dayjs.unix(nowTime));
}

/**
 * unixTimestampWithoutTimezones2nowTime
 * @param {number} timestamp timestamp s
 * @param {boolean} isReturnUnix
 */
export function unixTimestampWithoutTimezones2nowTime(timestamp, isReturnUnix = true) {
  let temp = deviceTimestampToDayJs(timestamp);
  if (isReturnUnix) {
    return temp.unix();
  }
  return temp.format("YYYY-MM-DD HH:mm:ss");
}

