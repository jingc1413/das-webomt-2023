<template>
  <el-alert v-if="viewMode != provideKeys.viewModePrintValue && visible" :type="data?.Style?.type ?? 'info'"
    :title="data.Value" :show-icon="true" :closable="false" />
  <div v-else>
  </div>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'
export default {
  name: 'MyAlert',
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
    visible() {
      if (this.data?.Style?.visibleValue && this.data?.Style?.visibleParam) {
        const param = this.dev.params.getParam(this.data?.Style?.visibleParam);
        if (param) {
          return param?.Value === this.data?.Style?.visibleValue;
        }
      }
      return true;
    },
  },
}
</script>
<style lang="scss" scoped></style>
