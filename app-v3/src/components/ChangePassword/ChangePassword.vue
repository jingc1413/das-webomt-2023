<template>
  <el-dialog :append-to-body="true" v-model="changePasswordDialogVisible" :title="viewTitle" width="700" align-center
    @closed="cancelPassword()" :close-on-click-modal="false" :close-on-press-escape="false"  :show-close="allowClose">
    <el-row v-loading="passwordLoading">
      <el-form ref="userFormRef" :model="userForm" :rules="userFormRules" label-width="auto" label-position="right">
        <el-form-item v-if="tipText">
          <el-alert :title="tipText" type="warning" :closable="false" />
        </el-form-item>

        <template v-if="changeType == 'current'">
          <el-form-item label="Current Password" prop="oldPassword">
            <el-input class="table-s-inp" v-model.trim="userForm.oldPassword" show-password>
            </el-input>
          </el-form-item>
        </template>

        <el-form-item label="New Password" prop="password">
          <el-input class="table-s-inp" v-model.trim="userForm.password" show-password>
          </el-input>
        </el-form-item>

        <el-form-item label="Confirm Password" prop="secPassword">
          <el-input class="table-s-inp" v-model.trim="userForm.secPassword" show-password>
          </el-input>
        </el-form-item>
        <el-space fill>
          <el-alert type="info" show-icon :closable="false">
            <p v-for="(infoItem, index) in accountsStore.passDesc" :key="index">
              {{ infoItem }}
            </p>
          </el-alert>
        </el-space>
      </el-form>
    </el-row>
    <template #footer>
      <div>
        <el-button type="default" v-if="allowClose" :disabled="passwordLoading" :loading="passwordLoading"
          @click="cancelPassword()">Cancel</el-button>
        <el-button type="primary" :disabled="passwordLoading" :loading="passwordLoading"
          @click="submitPassword()">Submit</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>
import { useAccounts } from '@/stores/account';

export default {
  name: 'ChangePassword',
  props: {
    isOpen: {
      default: false,
      type: Boolean
    },
    userId: {
      default: '',
      type: String
    },
    changeType: {
      default: 'current',
      type: String
    },
    tipText: {
      default: '',
      type: String
    },
    viewTitle: {
      default: 'Change Password Dialog',
      type: String
    },
    allowClose: {
      default: true,
      type: Boolean
    }
  },
  setup(props) {
    let accountsStore = useAccounts()
    return {
      accountsStore
    }
  },
  data() {

    return {
      changePasswordDialogVisible: false,
      passwordLoading: false,
      userForm: {
        oldPassword: '',
        password: '',
        secPassword: '',
      },
      userFormRules: {
        oldPassword: [{ required: true, trigger: 'change', message: 'Current Password is required' }],
        password: [{ required: true, trigger: 'change', message: 'New Password is required' }, {pattern: new RegExp(this.accountsStore.passRegexp), trigger: 'change', message: 'Password is weak' }],
        secPassword: [{ required: true, trigger: 'change', message: 'Confirm Password is required' }, {trigger: 'change', validator: this.validateSecPassword }],
      }
    }
  },
  watch: {
    isOpen() {
      this.changePasswordDialogVisible = this.isOpen;
      if (this.isOpen) {
        this.initFormData()
      }
    }
  },
  mounted() {
    this.changePasswordDialogVisible = this.isOpen;
    setTimeout(() => {
      this.$refs['userFormRef'].clearValidate();
    }, 50);
  },
  methods: {
    validateSecPassword(rule, value, callback) {
      if ((this.userForm.password !== value)) {
        callback(new Error("Password don't match"))
      } else {
        callback()
      }
    },
    initFormData() {
      this.userForm.oldPassword = '';
      this.userForm.password = '';
      this.userForm.secPassword = '';
      setTimeout(() => {
        this.$refs['userFormRef']?.clearValidate()
      }, 50);
    },
    submitPassword() {
      this.$refs['userFormRef']?.validate((valid, fields) => {
        if (valid) {
          this.passwordLoading = true
          this.changeType == 'otherUser' ? (this.changeOtherPassword()) : (this.changeSelfPassword())
        } else {
          // modal.msgError('Password is weak')
        }
      })
    },
    changeOtherPassword() {
      //....
    },
    changeSelfPassword() {
      let info = {
        newPassword: this.userForm.password,
        oldPassword: this.userForm.oldPassword
      }
      this.accountsStore.changeCurrentUserPassword(this.userForm.oldPassword, this.userForm.password).then((res) => {
        this.passwordLoading = false
        this.cancelPassword(false, true)
      }).catch(err => {
        console.error('modifyPassword');
        console.error(err);
        this.passwordLoading = false
      })
    },
    cancelPassword(isTip = false, change = false) {
      if (this.passwordLoading) return
      this.passwordLoading = false;
      this.changePasswordDialogVisible = false;
      this.$emit('dialogHide', change);
    }
  }
}


</script>

<style lang='scss' scoped></style>