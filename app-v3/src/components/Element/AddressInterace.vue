<!-- eslint-disable vue/no-v-model-argument -->
<template>
  <el-form v-loading="loading" size="small" label-width="40%" label-position="right"
    style="width: 80%; max-width: 800px;">
    <my-form-element :owner="owner" :data="data.Items[0]" />

    <el-form-item label-width="40%">
      <template #label>
        <label style="text-align: right;">Buffer</label>
      </template>
      <BuffTableView v-model:bufferData="bufferData" v-model:bufferOffset="bufferOffset"
        v-model:bufferLength="bufferLength">
      </BuffTableView>
    </el-form-item>

    <el-form-item label-width="40%">
      <el-button type="primary" plain @click="handleReadRegister()">
        Read
      </el-button>
      <el-button type="primary" plain @click="handleWriteRegister()">
        Write
      </el-button>
    </el-form-item>
    <!-- <span>{{bufferData}}, {{ bufferOffset }}, {{ bufferLength }}</span> -->
  </el-form>
</template>

<script>

import { useDasDevices } from "@/stores/das-devices";
import BuffTableView from '@/components/BuffTableView/BuffTableView.vue';

export default {
  name: "MyAddressInterfaceForm",
  components: { BuffTableView },
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
    const dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    return {
      dasDevices,
      dev,
    };
  },
  data() {
    const param = this.dev.params.getParam("TB2.P0CCC");
    return {
      loading: false,
      param,
      selectModule: "00",
      bufferLength: 4,
      bufferOffset: 0x0000,
      bufferData: [],
    };
  },
  methods: {
    handleReadRegister: async function () {
      this.loading = true;
      const selectModule = parseInt(this.selectModule, 16);
      const data = await this.dev.funcs.readRegister(selectModule, this.bufferOffset, this.bufferLength)
      this.bufferData = data;
      this.loading = false;
    },
    handleWriteRegister: async function () {
      this.loading = true;
      const selectModule = parseInt(this.selectModule, 16);
      const data = await this.dev.funcs.writeRegister(selectModule, this.bufferOffset, this.bufferLength, this.bufferData);
      this.bufferData = data;
      this.loading = false;
    },
  },
};
</script>
<style lang="scss" scoped></style>