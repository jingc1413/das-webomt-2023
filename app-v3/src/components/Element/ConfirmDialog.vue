<template>
  <el-dialog v-model="dialog.visible" :title="dialog.title" width="40%">
    <span v-if="dialog.content">
      {{ dialog.content }}
    </span>
    <el-form v-if="dialog.needInput" ref="formRef" :rules="formRules" :model="dialog.formModel" size="small"
      label-position="right" style="width: 100%; max-width: 800px;">
      <el-form-item prop="inputValue" :label="dialog.inputLabel">
        <el-input v-model="dialog.formModel.inputValue" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button v-if="dialog.supportCancel" @click="handleCancel()">Cancel</el-button>
        <el-button type="primary" @click="handleConfirm()">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { useAppStore } from '@/stores/app';


export default {
  name: 'MyConfirmDialog',
  props: {},
  setup() {
    const appStore = useAppStore();
    return {
      appStore,
    }
  },
  data() {
    const dialog = this.appStore?.confirmDialog;
    return {
      dialog,
    }
  },
  computed: {
    formRules() {
      return {
        inputValue: this.dialog?.inputRule,
      };
    },
  },
  mounted() {
    const ref = this.$refs.formRef;
    if (ref) {
      ref.clearValidate();
    }
  },

  methods: {
    handleCancel() {
      this.appStore.closeConfirmDialog(false);
    },
    handleConfirm() {
      const self = this;
      if (this.dialog.needInput) {
        const ref = this.$refs.formRef;
        if (ref) {
          ref.validate(async (valid, fields) => {
            if (valid) {
              self.appStore.closeConfirmDialog(true);
            }
          })
        }
      } else {
        this.appStore.closeConfirmDialog(true);
      }

    }
  }
}
</script>
<style lang="scss" scoped></style>