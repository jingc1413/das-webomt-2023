<template>
  <el-input :class="editable === false ? 'input-readonly' : ''" v-model="inputModel" :maxlength="param?.ByteSize"
    show-word-limit :readonly="editable === false">
    <template #append>
      <el-button-group v-if="editable">
        <el-button @click="handleSubmit()">
          <el-icon class="el-icon--right">
            <Check />
          </el-icon>
        </el-button>
        <el-button @click="handleCancel()">
          <el-icon class="el-icon--right">
            <Close />
          </el-icon>
        </el-button>
      </el-button-group>
      <span v-else>
        <el-button @click="editable = true">
          <el-icon class="el-icon--right">
            <EditPen />
          </el-icon>
        </el-button>
      </span>
    </template>
  </el-input>
</template>
<script>
import { useDasDevices } from '@/stores/das-devices'

export default {
  name: 'editableInput',
  props: {
    subID: String,
    param: Object,
    value: String,
  },
  setup() {
    const dasDevices = useDasDevices()
    return {
      dasDevices,
    }
  },
  data() {
    return {
      editable: false,
      inputModel: this.value,
    }
  },
  watch: {
    value() {
      this.inputModel = this.value;
    }
  },
  methods: {
    async handleSubmit() {
      try {
        var values = [];
        values.push({ oid: this.param.PrivOid, value: this.inputModel });
        this.dasDevices.setDeviceParameterValues({
          sub: this.subID,
          values, 
          showMessage: true,
        });
      } catch (e) {
        console.error(e)
      } finally {
        this.editable = false;
      }
    },
    handleCancel() {
      this.inputModel = this.value;
      this.editable = false;
    }
  }
}


</script>
