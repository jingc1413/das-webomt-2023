<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <template
    v-if="data.Items.length === 2 && param.UnitName == param2.UnitName && param.DataType == param2.DataType && param.UnitName === 'MHz'">
    <template v-for="(item, index) in data.Items">
      <template v-if="index === 0">
        <report-param :key="item.Key + '_' + index" :owner="owner" :data="item" />
      </template>
      <template v-else>
        <el-text :key="item.Key + '_' + index + '_span'" style="padding-inline: 4px;">~</el-text>
        <report-param :key="item.Key + '_' + index" :owner="owner" :data="item" />
      </template>
    </template>
  </template>
  <template v-else>
    <template v-for="(item, index) in data.Items">
      <template v-if="index === 0">
        <report-param :key="item.Key + '_' + index" :owner="owner" :data="item" />
      </template>
      <template v-else>
        {{ ", " }}
        <report-param :key="item.Key + '_' + index" :owner="owner" :data="item" />
      </template>
    </template>
  </template>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';


export default {
  name: 'MyParamGroup',
  components: {},
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
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
