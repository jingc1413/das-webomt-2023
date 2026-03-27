<template>
  <el-form  v-if="deviceInfo" :model="deviceInfo" size="small" label-width="120" label-position="right"
    style="width: 85%; max-width: 800px;">
    <el-form-item label="Device Type:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.DeviceTypeName" :maxlength="20" :rwValue="false"/>
    </el-form-item>
    <el-form-item label="Device name:">
      <my-form-editable-input :subID="subID" :data="{OID:'T02.P0030', Value: deviceInfo.DeviceName}" :maxlength="40" :rwValue="isOnline"/>
    </el-form-item>
    <el-form-item label="ID:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.SubID" :rwValue="false"/>
    </el-form-item>
    <el-form-item label="IP:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.IpAddress" :rwValue="false"/>
    </el-form-item>
    <el-form-item label="Device Location:">
      <my-form-editable-input :subID="subID" :data="{OID:'T02.P0023', Value: deviceInfo.InstalledLocation}" :maxlength="40" :rwValue="isOnline"/>
    </el-form-item>
    <el-form-item label="Alarm:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.AlarmState == 1 ? 'Yes' : 'No'" :rwValue="false"/>
    </el-form-item>
    <el-form-item label="Version:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.Version" :rwValue="false"/>
    </el-form-item>
    <el-form-item label="Element model:">
      <my-form-editable-input :subID="subID" :value="deviceInfo.ElementModelNumber" :rwValue="false"/>
    </el-form-item>
  </el-form>
</template>
<script>
import { useDasDevices } from '@/stores/das-devices';

export default {
  name: 'topoInfo',
  inject: ['viewMode'],
  props: {
    subID: String,
  },
  setup() {
    const dasDevices = useDasDevices();
    return {
      dasDevices,
    }
  },
  data() {
    return {
    }
  },
  computed: {
    deviceInfo: function () {
      const info = this.dasDevices.getDeviceInfo(this.subID);
      return info;
    },
    isOnline() {
      return this.deviceInfo.ConnectState < 6;
    }
  },
}
</script>