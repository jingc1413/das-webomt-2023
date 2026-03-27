<template>
  <el-dialog :title="data.Name" v-model="visible" width="40%">
    <my-layout v-for="item in data.Items" :key="item.Key" :owner="data" :data="item" />
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">Cancel</el-button>
        <my-button v-if="data.Submit" :owner="owner" :data="data.Submit" @click="visible = false" />
      </span>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'MyDialog',
  props: {
    owner: Object,
    data: Object,
    modelValue: Boolean,
  },
  emits: ['update:modelValue'],
  computed: {
    visible: {
      get() {
        return this.modelValue;
      },
      set(v) {
        this.$emit('update:modelValue', v);
      },
    }
  },
}
</script>
<style lang="scss" scoped></style>
