<template>
  <my-layout v-if="data.Type?.startsWith('Layout')" :owner="owner" :data="data" />
  <my-table v-else-if="data.Type?.startsWith('Table')" :owner="owner" :data="data" />
  <my-form v-else-if="data.Type?.startsWith('Form')" :owner="owner" :data="data" />
  <my-param-group v-else-if="data.Type?.startsWith('Component:ParamGroup')" :owner="owner" :data="data" />
  <my-statistic-group v-else-if="data.Type?.startsWith('Component:StatisticGroup')" :owner="owner" :data="data" />
  <my-param v-else-if="data.Type === 'Param'" :owner="owner" :data="data" />
  <my-label v-else-if="data.Type === 'Label'" :owner="owner" :data="data" />
  <my-alert v-else-if="data.Type === 'Alert'" :owner="owner" :data="data" />
  <my-button v-else-if="data.Type === 'Button'" :owner="owner" :data="data" />
  <span v-else-if="data.Type === 'Text'">{{ data.Value }}</span>
</template>

<script>
import provideKeys from '@/utils/provideKeys.js'
export default {
  name: "myElement",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  data() {
    return {
      provideKeys
    }
  },
};
</script>
<style lang="scss" scoped></style>
