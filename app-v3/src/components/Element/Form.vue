<template>
  <el-row style="height: 48px;">
    <el-col :span="20">
      <h4 v-if="data.Name" style="margin-left: 32px;">
        <span>{{ data.Name }}</span>
      </h4>
    </el-col>
    <el-col v-if="toolbarVisible" :span="4" style="display: flex;justify-content: flex-end;">
      <div class="toolbar" v-if="data?.Actions?.toolbar" style="margin-right: 16px;">
        <my-element v-for="item in data?.Actions?.toolbar?.Items" :key="item.Key" :owner="owner" :data="item" />
      </div>
      <div class="toolbar">
        <el-button v-if="supportRefresh && data.Key != 'address_interface'" type="primary" plain
          @click="getParameterValues(true)">
          Refresh
        </el-button>
      </div>
    </el-col>
  </el-row>
  <my-address-interface-form v-if="data.Key == 'address_interface'" :owner="owner" :data="data" />
  <el-scrollbar v-else-if="scrollable" class="my-form-height"
    max-height="calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height) - 34px)">
    <el-form v-loading="loading" ref="formRef" :model="formModel" :rules="rules" size="small"
      style="width: calc(100% - 16px); max-width: 800px;" label-position="right" label-width="40%"
      :disabled="data.InputDisabled || viewMode == provideKeys.viewModePrintValue" class="form-view-page">
      <my-form-element v-for="item in data.Items" :key="item.key" :owner="owner" :data="item" />
      <el-form-item v-if="submitAction && viewMode == provideKeys.viewModeDefaultValue " prop="Submit" label=" ">
        <el-button type="primary" :disabled="data.InputDisabled" @click="handleSubmit()">
          Submit
        </el-button>
      </el-form-item>
    </el-form>
  </el-scrollbar>
  <el-form v-else v-loading="loading" size="small" label-width="40%" label-position="right" class="form-view-page"
    style="width: 80%; max-width: 800px;" :disabled="viewMode == provideKeys.viewModePrintValue">
    <my-form-element v-for="item in data.Items" :key="item.Key" :owner="owner" :data="item" />
  </el-form>
</template>

<script>
import { useAppStore } from "@/stores/app";
import { useDasDevices } from "@/stores/das-devices";
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: 'MyForm',
  inject: ['viewMode'],
  props: {
    owner: {
      type: Object,
      default: undefined,
    },
    data: {
      type: Object,
      default: undefined,
    },
  },
  setup() {
    const appStore = useAppStore();
    const dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    return {
      appStore,
      dasDevices,
      dev,
      provideKeys,
    };
  },
  data() {
    const supportRefresh = this.viewMode != provideKeys.viewModePrintValue && this.data?.rOids?.length > 0;
    const toolbarVisible = this.viewMode != provideKeys.viewModePrintValue;
    const scrollable = this.viewMode != provideKeys.viewModePrintValue;

    const formModel = {};
    const rules = [];
    const submitAction = (this.data?.Actions?.submit?.accessKeys?.set) ? (this.data?.Actions?.submit || undefined) : undefined;
    const self = this;
    this.data.Items.forEach(item => {
      if (item.InputRules) {
        rules[item.Key] = item.InputRules;
      }
    })
    return {
      loading: true,
      formModel,
      rules,
      submitAction,
      toolbarVisible,
      scrollable,
      supportRefresh,
    };
  },
  mounted() {
    this.getParameterValues();
  },
  methods: {
    handleSubmit: async function () {
      const action = this.submitAction;
      if (action) {
        const self = this;
        const ref = this.$refs.formRef;
        if (ref) {
          ref.validate(async (valid, fields) => {
            if (valid) {
              self.loading = true;
              await self.doSubmitAction(action);
              self.loading = false;
            }
          })
        }
      }
    },
    doSubmitAction: async function (action) {
      const self = this;
      if (this.data?.Style?.confirmMessage) {
        this.appStore.openConfirmDialog({
          title: this.data?.Style?.confirmTitle || 'Confirm',
          content: this.data?.Style?.confirmMessage || '',
          callback: ok => {
            if (ok) {
              self.doAction(action);
            }
          }
        })
      } else {
        this.doAction(action);
      }
    },
    doAction: async function (action) {
      this.dev.doAction(action, {
        owner: this.owner,
      })
    },
    getParameterValues: async function (showMessage = false) {
      try {
        if (!this.supportRefresh) return;
        if (this.data?.oids?.length > 0) {
          this.loading = true;
          await this.dev.params.getParameterValues({
            oids: this.data.oids,
            values: this.data.defaultValues,
            showMessage,
          })
          this.loading = false;
        }
      } finally {
        this.loading = false;
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.my-form-height {
  // max-height: calc(100vh - var(--header-height) - var(--main-page-header-height) - var(--view-page-header-height) - 34px);
}

.form-view-page {
  & :deep(.el-form-item) {
    .el-form-item__label {
      pointer-events: none;
      height: auto !important;
      word-spacing: normal;
      word-wrap: break-word;
      word-break: break-word;
    }
  }
}
</style>
