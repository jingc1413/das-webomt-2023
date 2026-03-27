<template>
  <el-row>
    <el-col :span="12">
      <h4 v-if="data.Name">{{ data.Name }}</h4>
    </el-col>
    <el-col :span="12" v-if="viewMode !== provideKeys.viewModePrintValue"
      style="display: flex;justify-content: flex-end;">
      <div class="toolbar">
        <el-button v-if="isAdmin" v-hasPermissionAnd="['api.iam.users.create']" type="primary"
          @click="createUser()">Create</el-button>
        <el-button type="primary" v-hasPermissionAnd="['api.iam.users.list']" plain
          @click="getUserList()">Refresh</el-button>
      </div>
    </el-col>
  </el-row>
  <el-table :data="tableData" :border="false" table-layout="auto" style="width: 100%;" stripe
    :class="{ 'my-table-height': viewMode != provideKeys.viewModePrintValue }">
    <el-table-column prop="Name" label="Name" />
    <el-table-column prop="Roles" label="Role">
      <template #default="scope">
        {{ scope.row['Roles'].join(', ') }}
      </template>
    </el-table-column>
    <el-table-column prop="Action" label="">
      <template #default="scope">
        <template v-if="scope.row['Name']">
          <el-link v-if="isAdmin" type="primary" @click="editUser(scope.row)">Modify</el-link>
          <el-link v-if="isAdmin && scope.row['Name'] !== 'admin'" v-hasPermissionAnd="['api.iam.users.delete']"
            type="primary" @click="deleteUser(scope.row)" style="margin-left: 12px;">Delete</el-link>
        </template>
      </template>
    </el-table-column>
  </el-table>
  <el-dialog v-model="createDialogVisible" title="Create User" width="640px">
    <el-form v-if="createDialogVisible" ref="createFormRef" :rules="createRules" :model="userData" size="small"
      label-width="160" label-position="right" style="width: 100%; max-width: 800px;">
      <el-form-item prop="Name" label="Name">
        <el-input v-model="userData.Name" autocomplete="new-password" :minlength="5" :maxlength="12" show-word-limit
          :readonly="userData.nameReadonly" />
      </el-form-item>
      <el-form-item prop="Roles" label="Role">
        <el-select v-model="userData.Roles">
          <el-option v-for="item in roleList" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item prop="Password" label="Password">
        <el-input v-model="userData.Password" autocomplete="new-password" type="password" :minlength="8" :maxlength="12"
          show-word-limit show-password />
      </el-form-item>
      <el-form-item prop="Password2" label="Confirm Password">
        <el-input v-model="userData.Password2" autocomplete="new-password" type="password" :minlength="8" :maxlength="12"
          show-word-limit show-password />
      </el-form-item>
      <div class="tip">
        <p>Password must be a minimum of 8 characters in length</p>
        <p>Password must be a maximum of 12 characters in length</p>
        <p>Password must contain uppercase letters and lowercase letters</p>
        <p>Password must contain at least 1 out of 2 characters: numbers, and/or special characters (!@_)</p>
      </div>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="createDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleCreateUser()">Submit</el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="editDialogVisible" title="Modify User" width="640px">
    <el-form v-if="editDialogVisible" ref="editFormRef" :rules="editRules" :model="userData" size="small"
      label-width="160" label-position="right" style="width: 100%; max-width: 800px;">
      <el-form-item prop="Name" label="Name">
        <el-input v-model="userData.Name" autocomplete="new-password" :minlength="5" :maxlength="12" readonly />
      </el-form-item>
      <el-form-item prop="Roles" label="Role">
        <el-select v-model="userData.Roles">
          <el-option v-for="item in roleList" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="editDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleEditUser(ruleFormRef)">Submit</el-button>
      </span>
    </template>
  </el-dialog>
  <el-dialog v-model="deleteConfirmDialogVisible" title="Delete User" width="40%">
    <span>Confirm to delete the user</span>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="deleteConfirmDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleDeleteUser()">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
</template>
  
<script>
import { ElMessage } from "element-plus";
import { useAccounts } from "@/stores/account";
import { useAuthStore } from "@/stores/auth";
import provideKeys from '@/utils/provideKeys.js'

export default {
  name: "MyUsersTable",
  inject: ['viewMode'],
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const accounts = useAccounts();
    const auth = useAuthStore();
    return {
      auth,
      accounts,
      provideKeys,
    };
  },
  data() {
    const self = this;
    const userData = {
      Name: "",
      Password: "",
      Password2: "",
      Roles: []
    };
    const nameRules = [{
      required: true,
      validator: (rule, value, callback) => {
        if (!value) {
          return callback('Please input name')
        }
        setTimeout(() => {
          if (value.length < 5) {
            callback('Minimum of 5 characters in length')
          } else {
            const exist = self.tableData.find(row => row.Name === value);
            if (exist) {
              callback(`The user name already exists`)
            } else {
              callback();
            }
          }
        }, 500)
      },
      trigger: 'blur',
    }]
    const passwordRules = [{
      required: true,
      validator: (rule, value, callback) => {
        if (!value) {
          return callback('Please input password')
        }
        if (value.length < 8) {
          callback('Minimum of 8 characters in length')
        } else if (value.length > 12) {
          callback('Maximum of 12 characters in length')
        } else if (!value.match(/^[A-Za-z0-9!@_]{8,12}$/)) {
          callback('Password is weak')
        } else if (!value.match(/^.*[A-Z]+.*$/)) {
          callback('Password is weak')
        } else if (!value.match(/^.*[a-z]+.*$/)) {
          callback('Password is weak')
        } else {
          let count = 0;
          if (value.match(/^.*[A-Z]+.*$/)) {
            count += 1;
          }
          if (value.match(/^.*[a-z]+.*$/)) {
            count += 1;
          }
          if (value.match(/^.*[0-9]+.*$/)) {
            count += 1;
          }
          if (value.match(/^.*[!@_]+.*$/)) {
            count += 1;
          }
          if (count >= 3) {
            callback()
          } else {
            callback('Password is weak')
          }
        }

      },
      trigger: 'blur',
    }]
    const password2Rules = [{
      required: true,
      validator: (rule, value, callback) => {
        if (!value) {
          return callback(new Error('Please input password again'))
        }
        if (value !== this.userData.Password) {
          callback(new Error("Password don't match!"))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    }]
    const createRules = {
      Name: nameRules,
      Roles: [{ required: true, message: 'Please select role', trigger: 'blur' }],
      Password: passwordRules,
      Password2: password2Rules,
    };
    const editRules = {
      Roles: [{ required: true, message: 'Please select role', trigger: 'blur' }],
    };
    const isAdmin = this.auth.loginUserName === "admin";
    return {
      createRules,
      editRules,
      loading: false,
      isAdmin,
      userData,
      createDialogVisible: false,
      editDialogVisible: false,
      deleteConfirmDialogVisible: false,
    }
  },
  computed: {
    tableData() {
      return this.accounts.users;
    },
    hasAdmin() {
      return this.tableData?.find(v => v.Name === 'admin') !== undefined;
    },
    roleList() {
      return this.accounts.getRoleNameList;
    }
  },
  mounted() {
    this.accounts.getRoleList();
    this.getUserList();
  },
  methods: {
    resetUserData() {
      this.userData = {
        Name: this.hasAdmin ? "" : "admin",
        Password: "",
        Password2: "",
      };
    },
    getUserList: async function () {
      this.loading = true;
      await this.accounts.getUserList();
      this.loading = false;
    },
    createUser: async function () {
      if (this.accounts?.users?.length >= 5) {
        ElMessage.warning("You have exceeded the maximum number of users")
        return;
      }
      this.resetUserData();
      this.$nextTick(() => {
        this.createDialogVisible = true;
      })
      const ref = this.$refs.createFormRef;
      if (ref) {
        ref.clearValidate();
      }
    },
    handleCreateUser: function () {
      const self = this;
      const ref = this.$refs.createFormRef;
      if (ref) {
        ref.validate(async (valid, fields) => {
          if (valid) {
            self.loading = true;
            let info = {
              Name: self.userData.Name,
              Password: self.userData.Password,
              Roles: self.userData.Roles,
            }
            if (typeof info.Roles === 'string') {
              info.Roles = [info.Roles];
            }
            self.accounts.createUser(info).then(res => {
              if (res) {
                self.createDialogVisible = false;
              }
            }).finally(() => {
              self.loading = false;
            })
          }
        })
      }
    },
    editUser: async function (user) {
      this.resetUserData();
      this.$nextTick(() => {
        this.userData.Roles = user.Roles.join('');
        this.userData.Name = user.Name;
        this.editDialogVisible = true;
      })
      const ref = this.$refs.editFormRef;
      if (ref) {
        ref.clearValidate();
      }
    },
    handleEditUser: async function () {
      const self = this;
      const ref = this.$refs.editFormRef;
      if (ref) {
        ref.validate(async (valid, fields) => {
          if (valid) {
            const user = self.userData;
            self.loading = true;
            let info = {
              Roles: self.userData.Roles
            }
            if (typeof info.Roles === 'string') {
              info.Roles = [info.Roles];
            }
            self.accounts.editUser(user).then(res => {
              if (res) {
                self.editDialogVisible = false;
              }
            }).finally(() => {
              self.loading = false;
            })
          }
        })
      }
    },
    deleteUser: async function (user) {
      this.deleteConfirmDialogVisible = true;
      this.deleteUserData = user;
    },
    handleDeleteUser: async function () {
      const user = this.deleteUserData;
      this.loading = true;
      let res = await this.accounts.deleteUser(user.Name);
      this.loading = false;
      if (res == true) {
        this.deleteConfirmDialogVisible = false;
      }
    },
  },
};
</script>
<style lang="scss" scoped>
.tip {
  padding: 8px 16px;
  background-color: rgba(var(--el-color-info-rgb), .1);
  border-radius: 4px;
  border-left: 5px solid var(--el-color-info);
  margin: 20px 0;
}

.tip p:not(.tip-title) {
  font-size: .9rem;
}
</style>