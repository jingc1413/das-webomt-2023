<template>
  <el-dialog v-model="dialog.visible" :title="dialog.title" width="60%" class="file_view_dialog" :append-to-body="true" :show-close="dialog.supportCancel">
    <el-row justify="center" v-loading="loading">
      <el-col :span="24">
        <el-scrollbar style="height:60vh;background-color: var(--el-color-info-light-8);border-radius: 4px;padding: 8px;">
          <pre>{{ fileBody }}</pre>
        </el-scrollbar>
      </el-col>
    </el-row>
    <template #footer>
      <span class="dialog-footer">
        <el-button v-if="dialog.supportCancel" @click="handleCancel(false)">Cancel</el-button>
        <el-button v-if="dialog.supportSave" type="primary" @click="handleSave()">Save</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { useAppStore } from '@/stores/app';


export default {
  name: 'MyViewDialogDialog',
  props: {},
  setup() {
    const appStore = useAppStore();
    return {
      appStore,
    }
  },
  data() {
    return {
      loading: false,
      fileBody: null
    }
  },
  computed: {
    dialog() {
      return this.appStore?.viewFileDialog;
    }
  },
  watch: {
    'dialog.visible'(val) {
      console.log('dialog.visible', val);
      if (val == true) {
        this.handleLoading()
      } else {
        this.loading = false;
        this.fileBody = null;
      }
    }
  },
  methods: {
    handleLoading() {
      this.loading = true;
      console.log('loading', this.loading);
      this.dialog?.handleLoadingData().then(res => {
        if (res) {
          this.fileBody = res;
          this.loading = false;
        }
      }).catch((e) => {
        console.error(e);
      })
    },
    handleCancel(ok=false) {
      this.appStore.closeViewFileDialog(ok);
    },
    handleSave() {
      if (this.dialog?.handleSave) {
        this.dialog?.handleSave();
      }
    }
  }
}
</script>
<style lang="scss">
.file_view_dialog {
  margin-top: 15vh;
}
</style>