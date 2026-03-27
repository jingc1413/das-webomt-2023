<template>
  <template v-if="visible">
    <my-form-layout v-if="data.Type.startsWith('Layout')" :owner="owner" :data="data" />
    <el-form-item v-if="data.Type === 'Label'" :prop="data.Key">
      <h4 style="padding-block: 0;">{{ data.Value }}</h4>
    </el-form-item>
    <el-form-item
      v-else-if="data.Type === 'Param' && (!data?.Style?.hidden) &&
        (viewMode == provideKeys.viewModeDefaultValue || (viewMode != provideKeys.viewModeDefaultValue && data.Access != 'wo'))"
      :prop="data.Key">
      <template #label>
        <label style="text-align: right;">{{ paramLabel }}</label>
      </template>
      <my-param-input :owner="owner" :data="data"/>
      <!-- <my-form-editable-input :owner="owner" :data="data"/> -->
    </el-form-item>
    <el-form-item v-else-if="data.Type?.startsWith('Component:ParamGroup')" :prop="data.Key">
      <template #label>
        <label style="text-align: right;">{{ paramLabel }}</label>
      </template>
      <my-param-group-input :owner="owner" :data='data' />
    </el-form-item>
    <el-form-item v-else-if="data.Type === 'Button' && viewMode == provideKeys.viewModeDefaultValue" :prop="data.Key">
      <my-button :owner="owner" :data="data" />
    </el-form-item>
    <span v-else-if="data.Type === 'Text'">{{ data.Value }}</span>
  </template>
</template>

<script>

import { useDasDevices } from '@/stores/das-devices';
import provideKeys from '@/utils/provideKeys.js'
export default {
  name: 'MyFormElement',
  props: {
    owner: Object,
    data: Object,
  },
  inject: ['viewMode'],
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
      provideKeys
    };
  },
  computed: {
    paramLabel() {
      let label = this.data.Name;
      return label;
    },
    visible() {
      if (this.data?.Style?.visibleValue && this.data?.Style?.visibleParam) {
        const param = this.dev.params.getParam(this.data?.Style?.visibleParam);
        if (param) {
          console.log(this.data?.Style?.visibleParam, this.data?.Style?.visibleValue, param.Value)
          return param?.InputValue === this.data?.Style?.visibleValue;
        }
        return false;
      }
      return true;
    },
  },
}
</script>
<style lang="scss" scoped></style>
