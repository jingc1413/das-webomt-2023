import { defineStore } from 'pinia'
import apix from "@/api"
import { ElMessage } from 'element-plus';
import { translator as t } from '@/i18n';


export const useAccounts = defineStore('accounts', {
  state: () => ({
    users: [],
    roles: [],
    passRegexp: '^(?=(?:[^a-z]*[a-z]){1,})(?=(?:[^A-Z]*[A-Z]){1,})(?=(?:[^0-9!@_]*[0-9!@_]){1,})[A-Za-z0-9!@_]{8,12}$',
    passDesc: [
      "Password must be a minimum of 8 characters in length",
      "Password must be a maximum of 12 characters in length",
      "Password must contain uppercase letters and lowercase letters",
      "Password must contain at least 1 out of 2 characters: numbers, and/or special characters (!@_)"
    ],
  }),
  getters: {
    getPassRegexp:(state)=>state.passRegexp,
    getPassDesc:(state)=>state.passDesc,
    getRoleNameList:(state)=>state.roles.map(r=>r.Name),
  },
  actions: {
    setup: async function () {
    },
    getUserList: async function () {
      await apix.getUserList().then(users => {
        this.users = users;
      }).catch(e => {
        console.log(e);
      });
    },
    createUser: async function (user) {
      return new Promise((resolve, reject) => {
        apix.createUser(user).then(resp => {
          ElMessage.success(t("tip.UserCreatedSuccessfully"));
          this.getUserList();
          return resolve(true);
        }).catch(e => {
          ElMessage.error(t("tip.RequestFailed"));
          return resolve(false);
        });
      })
    },
    editUser: async function (user) {
      return new Promise((resolve, reject) => {
        apix.editUser(user).then(result => {
          ElMessage.success(t("tip.UserModifiedSuccessfully"));
          this.getUserList();
          return resolve(true);
        }).catch(e => {
          ElMessage.error(t("tip.RequestFailed"));
          return resolve(false);
        })
      })
    },
    deleteUser: async function (user) {
      return new Promise((resolve, reject) => {
        apix.deleteUser(user).then(result => {
          ElMessage.success(t("tip.UserDeletedSuccessfully"));
          this.getUserList();
          return resolve(true);
        }).catch(e => {
          ElMessage.error(t("tip.RequestFailed"));
          return resolve(false);
        })
      })
    },
    changeCurrentUserPassword: function (password, NewPassword) {
      return new Promise((resolve, reject) => {
        apix.changeCurrentUserPassword(password, NewPassword).then(result => {
          ElMessage.success("Change password Successfully");
          return resolve(true);
        }).catch(e => {
          if (e?.message) {
            ElMessage.error(e?.message);
          } else {
            ElMessage.error("Incorrect current password");
          }
          return reject(false);
        })
      })
    },
    getRoleList: function () {
      apix.getRoleList().then(res => {
        this.roles = res;
      }).catch(e => {
        console.log(e);
      });
    },
  },
});