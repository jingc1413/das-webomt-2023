<template>
  <template
    v-if="data.Items.length === 2 && param.DataType == param2.DataType && param.UnitName === param2.UnitName && param.UnitName === 'MHz'">
    <el-col :span="11">
      <my-param-input :owner="owner" :data="data.Items[0]" />
    </el-col>
    <el-col :span="2" style="text-align: center;">
      <span>~</span>
    </el-col>
    <el-col :span="11">
      <my-param-input :owner="owner" :data="data.Items[1]" />
    </el-col>
  </template>
  <template v-else-if="data.Key.endsWith('_digital_signal_bandwidth')">
    <el-row>
      <my-param-input :owner="owner" :data="data.Items[0]" />
    </el-row>
    <el-row style="margin-top: 8px;">
      <my-param-input :owner="owner" :data="data.Items[1]" />
    </el-row>
  </template>
  <template v-else>
    <template v-for="(item, index) in data.Items">
      <template v-if="index === 0">
        <el-col :key="item.Key" :span="11">
          <my-param-input :owner="owner" :data="item" />
        </el-col>
      </template>
      <template v-else>
        <el-col :key="item.Key" :span="2" style="text-align: center;">
          <span>{{ "" }}</span>
        </el-col>
        <el-col :key="item.Key" :span="11">
          <my-param-input :owner="owner" :data="item" />
        </el-col>
      </template>
    </template>
  </template>
</template>

<script>

import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: 'MyParamGroupInput',
  props: {
    owner: Object,
    data: Object,
  },
  inject: ['viewMode'],
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
      provideKeys
    };
  },
  computed: {
    param() {
      return this.dev.params.getParam(this.data.Items[0].OID);
    },
    param2() {
      return this.dev.params.getParam(this.data.Items[1].OID);
    },
  },
}
</script>
<style lang="scss" scoped></style>
