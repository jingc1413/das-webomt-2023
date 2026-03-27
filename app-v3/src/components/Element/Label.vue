<template>
  <div>
    <el-tooltip :content="data.Value" v-if="viewMode != provideKeys.viewModePrintValue" :trigger-keys="[]" :show-after="500">
      <span style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">{{data.Value}}</span>
    </el-tooltip>
    <span v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">{{data.Value}}</span>
  </div>
</template>

<script>
import provideKeys from '@/utils/provideKeys.js'
export default {
  name:'MyLabel',
  props: {
    owner: Object,
    data: Object,
  },
  inject: ['viewMode'],
  data() {
    return {
      provideKeys
    }
  },
}
</script>
<style lang="scss" scoped>

</style>
