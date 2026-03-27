<template>
  <template v-if="provideKeys.viewModePrintValue == viewMode">
    <my-param v-if="data.Type === 'Param'" key="param" :owner="owner" :data="data" />
    <my-param-group v-else-if="data.Type?.startsWith('Component:ParamGroup')" key="component" :owner="owner"
      :data='data' />
    <span v-else-if="data.Type === 'Text'">{{ data.Value }}</span>
  </template>
  <template v-else>
    <my-param-input v-if="data.Type === 'Param'" key="param" :owner="owner" :data="data" />
    <my-param-group-input v-else-if="data.Type?.startsWith('Component:ParamGroup')" key="component" :owner="owner"
      :data='data' />
    <my-button v-else-if="data.Type === 'Button'" key="button" :owner="owner" :data="data" />
    <span v-else-if="data.Type === 'Text'">{{ data.Value }}</span>
  </template>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'
export default {
  name: 'MyToolbarElement',
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
      provideKeys
    };
  },
  computed: {
    param() {
      if (this.data.Type === "Param") {
        return this.dev.params.getParam(this.data.OID)
      }
      return undefined;
    },
  },
}
</script>
<style lang="scss" scoped></style>
