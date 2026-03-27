<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <template v-if="pageType == 'name'">
    {{ param?.Name }}
  </template>
  <template v-else>
    <div style="display: inline-block;">
      <template v-if="paramValue === undefined">
        <span></span>
      </template>
      <template v-else-if="data?.Style?.hidden">
        {{ data.Value }}
      </template>
      <template v-else-if="data?.Style?.readonly">
        <el-text key="input-readonly" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'button'">
        <report-button key="input-button" :owner="owner" :data="data" />
      </template>
      <template v-else-if="data?.Style?.input === 'switch'">
        <el-text v-if="switchData.activeValue === paramValue" key="input-switch">
          {{ formatValueWithUnit }}
        </el-text>
        <el-text v-else-if="switchData.inactiveValue === paramValue" key="input-switch-info"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
        <el-text v-else :style="paramStyle" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'alarm'">
        <el-text v-if="paramValue === '00'" key="status-alarm-00" type="success">
          <el-icon>
            <Warning />
          </el-icon>
          <span :style="paramStyle"
            style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
            {{ formatValueWithUnit }}
          </span>
        </el-text>
        <el-text v-else-if="paramValue === '01'" key="status-alarm-01" type="danger">
          <el-icon>
            <Warning />
          </el-icon>
          <span :style="paramStyle"
            style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
            {{ formatValueWithUnit }}
          </span>
        </el-text>
        <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'sync'">
        <el-text v-if="paramValue === '01'" key="status-sync-01" type="danger">
          <el-icon>
            <Warning />
          </el-icon>
          <span :style="paramStyle"
            style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
            {{ formatValueWithUnit }}
          </span>
        </el-text>
        <el-text v-else-if="paramValue === '00'" key="status-sync-00" type="success">
          <el-icon>
            <Warning />
          </el-icon>
          <span :style="paramStyle"
            style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
            {{ formatValueWithUnit }}
          </span>
        </el-text>
        <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'Select'">
        <el-text key="select" :style="paramStyle" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <el-button-group v-else-if="data?.Style?.input === 'buttonGroup'">
        <report-button v-for="item in data.Items" :key="item.key" :owner="owner" :data="item" />
      </el-button-group>
      <template v-else-if="data?.Style?.input === 'number'">
        <el-text key="input-number" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'password'">
        <el-text key="input-password" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValue }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'binary'">
        <el-text key="input-binary" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValue }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'ipv4'">
        <el-text key="input-ipv4addr" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValue }}
        </el-text>
      </template>
      <template v-else-if="data?.Style?.input === 'datetime'">
        <el-text key="input-datetime" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValue }}
        </el-text>
      </template>
      <template v-else>
        <span key="input-default" :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
          {{ formatValueWithUnit }}
        </span>
      </template>
    </div>
  </template>
</template>

<script>

import { useDasDevices } from "@/stores/das-devices";
import { Warning } from "@element-plus/icons-vue"

export default {
  name: 'ReportParam',
  components: {},
  props: {
    owner: Object,
    data: Object,
    pageType: {
      default: 'param',
      type: String
    }
  },
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
    };
  },
  data() {
    return {
    }
  },
  computed: {
    param() {
      return this.dev.params.getParam(this.data?.OID)
    },
    paramStyle() {
      return this.param?.getShowStyle();
    },
    paramValue() {
      return this.param?.getValue({ defaultValue: this.data?.Value });
    },
    formatValue() {
      return this.param?.getShowValue({ defaultValue: this.data?.Value });
    },
    formatValueWithUnit() {
      return this.param?.getShowValue({ defaultValue: this.data?.Value, withUnit: true });
    },
    switchData() {
      return this.param?.getSwitchData(this.data?.Style);
    }
  }
}
</script>
<style lang="scss" scoped></style>
