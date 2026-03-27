<!-- eslint-disable vue/max-attributes-per-line -->
<template>
  <template v-if="viewMode == provideKeys.viewModePrintValue || (accessSet == false)">
    <el-button link :disabled="data.InputDisabled">
      {{ buttonLabel }}
    </el-button>
  </template>
  <template v-else-if="data?.Type === 'Param'">
    <el-link v-if="buttonType === 'text'" type="primary" :disabled="data.InputDisabled" @click="doClick()">
      {{ buttonLabel }}
    </el-link>
    <el-button v-else :type="buttonType" :plain="data?.Style?.plain" :disabled="data.InputDisabled" @click="doClick()">
      {{ buttonLabel }}
    </el-button>
  </template>
  <template v-else-if="data?.Actions?.click">
    <el-link v-if="buttonType === 'text'" type="primary" :disabled="data.InputDisabled" @click="doClick()">
      {{ buttonLabel }}
    </el-link>
    <el-button v-else :type="buttonType" :plain="data?.Style?.plain" :disabled="data.InputDisabled" @click="doClick()">
      {{ buttonLabel }}
    </el-button>
  </template>
  <template v-else>
    <el-link v-if="buttonType === 'text'" type="primary" :disabled="data.InputDisabled">
      {{ buttonLabel }}
    </el-link>
    <el-button v-else :type="buttonType" :plain="data?.Style?.plain" :disabled="data.InputDisabled">
      {{ buttonLabel }}
    </el-button>
  </template>
</template>

<script>
import { useAppStore } from '@/stores/app';
import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: "myButton",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  emits: ["click"],
  setup() {
    const appStore = useAppStore();
    const dev = useDasDevices().currentDevice;
    return {
      appStore,
      dev,
      provideKeys,
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
    const accessSet = this.data?.accessKeys?.set || false;
    return {
      param,
      buttonData,
      accessSet
    };
  },
  computed: {
    buttonLabel() {
      const buttonData = this.buttonData;
      const param = this.param;
      if (param) {
        if (buttonData) {
          if (buttonData?.activeText && buttonData?.activeValue === param?.Value) {
            return buttonData?.activeText;
          } else if (buttonData?.inactiveText && buttonData.inactiveValue === param?.Value) {
            return buttonData?.inactiveText;
          }
        }
        return param.getShowValue({ defaultValue: this.data?.Value })
      }

      return this.data?.Value;
    },
    buttonType() {
      const buttonData = this.buttonData;
      const param = this.param;
      if (buttonData && param) {
        if (buttonData?.activeType && buttonData?.activeValue === param?.Value) {
          return buttonData?.activeType;
        } else if (buttonData?.inactiveType && buttonData?.inactiveValue === param?.Value) {
          return buttonData?.inactiveType;
        }
      }
      return this.data?.Style?.type || "info";
    },
  },
  methods: {
    doClick() {
      const self = this;
      const action = this.data?.Actions?.click;
      if (action) {
        if (this.data?.Style?.confirmMessage) {
          this.appStore.openConfirmDialog({
            title: this.data?.Style?.confirmTitle || 'Confirm',
            content: this.data?.Style?.confirmMessage || '',
            callback: ok => {
              if (ok) {
                self.doClickAction(action);
              }
            }
          })
        } else {
          this.doClickAction(action);
        }
      }
    },
    doClickAction(action) {
      this.doAction(action);
      this.$emit("click");
    },
    doAction: async function (action) {
      this.dev.doAction(action, {
        owner: this.owner,
        param: this.param,
      })
    },
  },
};
</script>
<style lang="scss" scoped></style>
