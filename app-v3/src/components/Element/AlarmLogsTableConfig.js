
const TableProps = {
  "EventTime": "Event Time",
  "AlarmName": "Alarm Name",
  "AlarmSeverity": "Alarm Severity",
  "AlarmStatus": "Alarm Status",
  "DeviceTypeName": "Device Type Name",
  "DeviceSubID": "Device Sub",
  "SiteName": "Site Name",
  "DeviceName": "Device Name",
  "SerialNumber": "Serial Number",
  "SoftwareVersion": "Software Version",
}

const TableQueryItem = [
  { key: 'DeviceTypeName', name: TableProps['DeviceTypeName'], type: 'select', options: [], value: null },
  { key: 'DeviceName', name: TableProps['DeviceName'], type: 'string', value: null },
  { key: 'AlarmName', name: TableProps['AlarmName'], type: 'string', value: null },
  { key: 'AlarmSeverity', name: TableProps['AlarmSeverity'], type: 'select', options: [], value: null },
]

const TableDefaultSort = { prop: 'EventTime', order: 'descending' }


export default {
  TableQueryItem,
  TableProps,
  TableDefaultSort
}