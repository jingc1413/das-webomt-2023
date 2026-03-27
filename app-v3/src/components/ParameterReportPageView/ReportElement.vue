<template>
  <report-layout v-if="data.Type?.startsWith('Layout')" :owner="owner" :data="data" />
  <report-table v-else-if="data.Type?.startsWith('Table')" :owner="owner" :data="data" />
  <report-form v-else-if="data.Type?.startsWith('Form')" :owner="owner" :data="data" />
  <report-param-group v-else-if="data.Type?.startsWith('Component:ParamGroup')" :owner="owner"
    :data="data" />
  <report-param v-else-if="data.Type === 'Param'" :owner="owner" :data="data" style="height: 24px;"/>
  <span v-else-if="data.Type === 'Label'" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">{{data.Value}}</span>
  <el-text v-else-if="data.Type === 'Button'">{{ data.Value??data.Name }}</el-text>
  <span v-else-if="data.Type === 'Text'" style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">{{ data.Value }}</span>
</template>

<script>

export default {
  name: "ReportElement",
  props: {
    owner: Object,
    data: Object,
  },
  data() {
    return {
    }
  },
  setup() {
    return {};
  },
  mounted() { },
};
</script>
<style lang="scss" scoped></style>
