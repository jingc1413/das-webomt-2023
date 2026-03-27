<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <el-tooltip :disabled="!hasTips" :trigger-keys="[]" :show-after="500">
    <template #content>
      <span v-if="param.Tips">{{ `Tips: ${param.Tips}` }}<br></span>
      <span v-if="param.Value !== undefined">{{ `Value: ${param.Value}` }}</span>
      <span v-if="param.UnitName !== undefined">{{ `, Unit: ${param.UnitName}` }}</span>
      <span v-if="param.NumberMin != undefined || param.NumberMax != undefined">
        {{ `, Range: ${param.NumberMin != undefined ? param.NumberMin : ""}` }}
        {{ ` ~ ` }}
        {{ `${param.NumberMax != undefined ? param.NumberMax : " "}` }}
      </span>
      <!-- <span v-if="param.TextMin != undefined || param.TextMax != undefined">
        {{ `, Length: ${param.TextMin != undefined ? param.TextMin : ""}` }}
        {{ ` ~ ` }}
        {{ `${param.TextMax != undefined ? param.TextMax : " "}` }}
      </span> -->
      <p v-if="!auth.superTestDisabled && !appStore.debugTooltipDisabled">
        {{ param.Groups[0] || '' }}<br>
        {{ data.OID + ": " + param.Name + ", " + param.Access + ", " + param.DataType }}<br>
      </p>
    </template>
    <template v-if="paramValue === undefined">
      <span>{{ data?.Style?.nullValue || "" }}</span>
    </template>
    <template v-else-if="data?.Style?.hidden" />
    <template v-else-if="data?.Style?.readonly">
      <el-text key="input-readonly" :style="paramStyle"
        style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
        {{ formatValueWithUnit }}
      </el-text>
    </template>
    <template v-else-if="data?.Style?.input === 'button'">
      <my-button key="input-button" :owner="owner" :data="data" />
    </template>
    <template v-else-if="data?.Style?.input === 'switch'">
      <el-tag v-if="switchData.activeValue === paramValue" key="input-switch">
        {{ formatValueWithUnit }}
      </el-tag>
      <el-tag v-else-if="switchData.inactiveValue === paramValue" key="input-switch-info" type="info">
        {{ formatValueWithUnit }}
      </el-tag>
      <el-text v-else :style="paramStyle" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
        {{ formatValueWithUnit }}
      </el-text>
    </template>
    <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'alarm'">
      <el-text v-if="paramValue === '00'" key="status-alarm-00" type="success">
        <el-icon>
          <CircleCheck />
        </el-icon>
        <span :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
          {{ formatValue }}
        </span>
      </el-text>
      <el-text v-else-if="paramValue === '01'" key="status-alarm-01" type="danger">
        <el-icon>
          <Warning />
        </el-icon>
        <span :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
          {{ formatValue }}
        </span>
      </el-text>
      <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
        {{ formatValue }}
      </el-text>
    </template>
    <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'sync'">
      <el-text v-if="paramValue === '01'" key="status-sync-01" type="danger">
        <el-icon>
          <Warning />
        </el-icon>
        <span :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
          {{ formatValue }}
        </span>
      </el-text>
      <el-text v-else-if="paramValue === '00'" key="status-sync-00" type="success">
        <el-icon>
          <CircleCheck />
        </el-icon>
        <span :style="paramStyle"
          style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
          {{ formatValue }}
        </span>
      </el-text>
      <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
        {{ formatValue }}
      </el-text>
    </template>
    <template v-else-if="data?.Style?.input === 'Select'">
      <el-text key="select" :style="paramStyle" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
        {{ formatValue }}
      </el-text>
    </template>
    <el-button-group v-else-if="data?.Style?.input === 'buttonGroup'">
      <my-button v-for="item in data.Items" :key="item.key" :owner="owner" :data="item" />
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
        <span>{{ viewMode == provideKeys.viewModePrintValue ? formatValue : "******" }}</span>
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
  </el-tooltip>
</template>

<script>
import { useAppStore } from "@/stores/app";
import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import provideKeys from '@/utils/provideKeys.js'
import { Warning } from "@element-plus/icons-vue"

export default {
  name: 'MyParam',
  components: {},
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const auth = useAuthStore();
    const appStore = useAppStore();
    const dev = useDasDevices().currentDevice;
    return {
      auth,
      appStore,
      dev,
      provideKeys
    };
  },
  data() {
    return {}
  },
  computed: {
    param() {
      return this.dev.params.getParam(this.data?.OID);
    },
    hasTips() {
      if (this.param !== undefined) {
        if (this.param.Tips != undefined) {
          return true;
        }
        if (this.param.Value != undefined) {
          return true;
        }
        if (this.param.UnitName != undefined) {
          return true;
        }
        if (this.param.NumberMax != undefined || this.param.NumberMin != undefined) {
          return true;
        }
        if (this.param.TextMax != undefined || this.param.TextMin != undefined) {
          return true;
        }
        if (!this.auth.superTestDisabled && !this.appStore.debugTooltipDisabled) {
          return true;
        }
      }
      return false;
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
    },
  }
}
</script>
<style lang="scss" scoped></style>
