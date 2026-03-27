<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  {{ buttonLabel }}
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';


export default {
  name: "ReportButton",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  emits: ["click"],
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
    };
  },
  data() {
    const param = this.dev.params.getParam(this.data?.OID);
    const layoutStyle = this.data?.Style;
    const buttonData = {};
    if (layoutStyle?.activeValue && layoutStyle?.inactiveValue) {
      buttonData.activeValue = layoutStyle?.activeValue
      buttonData.activeText = layoutStyle?.activeText
      buttonData.activeType = layoutStyle?.activeType
      buttonData.inactiveValue = layoutStyle?.inactiveValue
      buttonData.inactiveText = layoutStyle?.inactiveText
      buttonData.inactiveType = layoutStyle?.inactiveType
    }
    return {
      param,
      buttonData,
    };
  },
  computed: {
    buttonLabel() {
      const buttonData = this.buttonData;
      const param = this.param;
      if (buttonData && param) {
        if (buttonData?.activeText && buttonData?.activeValue === param?.Value) {
          return buttonData?.activeText;
        } else if (buttonData?.inactiveText && buttonData.inactiveValue === param?.Value) {
          return buttonData?.inactiveText;
        }
      }
      return this.data?.Value;
    }
  },
  methods: {
  },
};
</script>
<style lang="scss" scoped></style>
