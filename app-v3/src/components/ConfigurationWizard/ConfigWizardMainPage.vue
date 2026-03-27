<template>
  <el-dialog v-model="dialogVisible" :title="wizardTitle" :fullscreen="true" :show-close="true"
    :close-on-click-modal="false" @closed="cancelDialogVisible()" :close-on-press-escape="false">
    <el-row>

      <el-col class="wizard-main-page-title">
        <div>
          <span>{{ data?.Items[activeSteps]['Name'] }}</span>
        </div>

        <div>
          <el-button :disabled="activeSteps == 0" @click="handlePrevious()">Previous step</el-button>
          <el-button link>
            {{ activeSteps + 1 }} / {{ data?.Items.length ?? 1 }}
          </el-button>
          <el-button :disabled="activeSteps == (data?.Items.length-1)" @click="handleNext()">Next
            step</el-button>
        </div>

      </el-col>

      <el-col v-if="data && data?.Items">
        <el-scrollbar style="height: calc(100vh - var(--header-height) - var(--view-page-header-height));">
          <my-view-page :owner="owner" :page="data.Items[activeSteps]"
            style="height: calc(100vh - var(--header-height) - var(--view-page-header-height));" />
        </el-scrollbar>
      </el-col>
    </el-row>
  </el-dialog>
</template>

<script>

export default {
  name: 'ConfigWizardMainPage',
  props: {
    owner: Object,
    data: Object,
    wizardTitle: {
      default: 'Wizard',
      type: String
    },
    isOpen: {
      default: false,
      type: Boolean
    },
  },
  data() {
    return {
      activeSteps: 0,
      dialogVisible: false,
    }
  },
  watch: {
    isOpen() {
      this.dialogVisible = this.isOpen;
    }
  },
  mounted() {
    this.dialogVisible = this.isOpen;
  },
  methods: {
    handlePrevious() {
      this.activeSteps--;
    },
    handleNext() {
      this.activeSteps++;
    },
    cancelDialogVisible(isTip = false, change = false) {
      if (this.passwordLoading) return
      this.passwordLoading = false;
      this.changePasswordDialogVisible = false;
      this.$emit('dialogHide', change);
    }
  },
}

</script>

<style scoped lang="scss">
.wizard-main-page-title {
  display: flex;
  justify-content: space-between;
}
</style>