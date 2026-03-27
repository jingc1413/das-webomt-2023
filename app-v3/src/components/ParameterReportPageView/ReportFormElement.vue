<template>
  <report-form-layout v-if="data.Type === 'Layout'" :owner="owner" key="layout" :data="data" />
  <report-param v-else-if="data.Type === 'Param'" key="param" :owner="owner"
    :data="data"/>
  <report-param-group v-else-if="data.Type?.startsWith('Component:ParamGroup')" key="param-group"
    :owner="owner" :data='data' />
  <span v-else-if="data.Type === 'Text'">{{ data.Value }}</span>
</template>

<script>

export default {
  name: 'ReportFormElement',
  props: {
    owner: Object,
    data: Object,
  },
  computed: {
    paramLabel() {
      let label = this.data.Name;
      return label;
    },
  },
  created() {
  },
}
</script>
<style lang="scss" scoped></style>
