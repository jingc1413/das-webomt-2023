
const InventoryTableProps = {
  "SubID": "Device Sub ID",
  "DeviceTypeName": "Device Type",
  "InstalledLocation": "Device Location",
  "DeviceName": "Device Name",
  "ElementSerialNumber": "Element Serial Number",
  "ElementModelNumber": "Element Model Number",
  "Version": "Software Version",
  "ConnectState": "Status",
  "IpAddress": "IP",
  "RouteAddress": "Route",
  "LifeTime": "Life Time",
  "ElementOperatingTemperature": "Element Operating Temperature"
}

const InVentoryTableQueryItem = [
  { key: 'DeviceTypeName', name: InventoryTableProps['DeviceTypeName'], type: 'select', options: [], value: null },
  { key: 'DeviceName', name: InventoryTableProps['DeviceName'], type: 'string', value: null },
  { key: 'SubID', name: InventoryTableProps['SubID'], type: 'string', value: null },
]


export default {
  InVentoryTableQueryItem,
  InventoryTableProps
}