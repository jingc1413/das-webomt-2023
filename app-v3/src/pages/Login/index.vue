<template>
  <div class="login-container" ref="loginContainerRef">
    <el-form ref="loginFormRef" id="loginForm" :model="loginForm" :rules="loginRules" class="login-form"
      auto-complete="on" label-position="left">
      <div class="title-container">
        <h3 class="title">{{ loginTitle }}</h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <el-input v-model="loginForm.username" placeholder="Username" name="username" type="text" tabindex="1"
          auto-complete="on" class="login_input" />
      </el-form-item>

      <el-form-item prop="password">
        <span class="svg-container">
          <svg-icon icon-class="password" />
        </span>
        <el-input ref="passwordRef" v-model="loginForm.password" :type="passwordType" placeholder="Password"
          name="password" tabindex="2" auto-complete="on" @keyup.enter="handleLogin" />
        <span class="show-pwd" @click="showPwd">
          <svg-icon :iconClass="passwordType === 'password' ? 'eye' : 'eye-open'" />
        </span>
      </el-form-item>

      <el-button :loading="loading" type="primary" style="width: 100%; margin-bottom: 30px;height: 40px;"
        @click.prevent="handleLogin">Login</el-button>

      <div class="tips">
        <span class="err_info">{{ messageBody }}</span>
      </div>
    </el-form>

    <my-change-password v-if="passwordDialogVisible" :isOpen="passwordDialogVisible" :allowClose="false"
      @dialogHide="closeChangePassword" />
  </div>
</template>

<script setup>
import { ref, reactive, watch, nextTick, toRaw } from "vue";
import { useAuthStore } from "@/stores/auth.js";
import { onMounted } from "vue";
import { useBase64 } from '@vueuse/core'
import { useAppStore } from "@/stores/app";
import { computed } from "vue";

const loginContainerRef = ref(null);

const loginFormRef = ref(null);
const passwordRef = ref(null);

const authStore = useAuthStore();

const loginForm = reactive({
  username: "",
  password: "",
});
const messageBody = ref(null);

const appStore = useAppStore();
let loginTitle = computed(() => appStore.loginTitle);


let validateUsername = (rule, value, callback) => {
  callback();
};
let validatePassword = (rule, value, callback) => {
  callback();
};

const loginRules = {
  username: [
    {
      required: true,
      message: "Please Enter Username!",
      trigger: "blur",
      validator: validateUsername,
    },
  ],
  password: [
    {
      required: true,
      message: "Please Enter Password!",
      trigger: "blur",
      validator: validatePassword,
    },
  ],
};

const loading = ref(false);
const passwordType = ref("password");

const showPwd = () => {
  if (passwordType.value === "password") {
    passwordType.value = "";
  } else {
    passwordType.value = "password";
  }
  nextTick(() => {
    passwordRef.value && passwordRef.value.focus();
  });
};

const handleLogin = () => {
  loginFormRef.value.validate(async (isValid, fields) => {
    if (isValid) {
      loading.value = true;
      const response = await authStore.login(loginForm.username, loginForm.password);
      if (response == undefined) {
        messageBody.value = response;
        loading.value = false;
        return;
      }

      let userInfo = await authStore.getCurrentUser();
      if (!userInfo) {
        loading.value = false;
        return
      }
      // if (userInfo.FirstTimeLogin == true) {
      //   openPasswordDialog();
      //   return;
      // }
      appStore.gotoDefaultRoute();
    } else {
      return false;
    }
  });
};

onMounted(() => {
  appStore.getLoginBackground(appStore.appInfo?.Schema).then(res => {
    if (res) {
      const reader = new FileReader();
      reader.onloadend = () => {
        loginContainerRef.value.style.backgroundImage = `url(${reader.result})`;
        document.getElementById('loginForm').style.backgroundColor = '#06060680';
      };
      reader.readAsDataURL(res);
    }
  })
})

//------First  Login---------

let passwordDialogVisible = ref(false);

function openPasswordDialog() {
  appStore.openConfirmDialog({
    title: 'Tip',
    content: "You will be required to change your password on first login",
    supportCancel: false,
    callback: ok => {
      if (ok) {
        passwordDialogVisible.value = true;
      }
    }
  })
}

function closeChangePassword(change) {
  if (change) {
    appStore.openConfirmDialog({
      title: 'Tip',
      content: "Your password has been change,please login again",
      supportCancel: false,
      callback: ok => {
        if (ok) {
          authStore.logout().then(()=>{
            window.location.reload()
          })
        }
      }
    })
  }
}

</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #283443;
$light_gray: #fff;
$cursor: #fff;
$dark_gray: #889aa4;
$primary_text: #303133;

@supports (-webkit-mask: none) and (not (cater-color: $primary_text)) {
  .login-container .el-input input {
    color: $primary_text;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    // display: inline-block;
    height: 52px;
    width: 85%;

    .el-input__wrapper {
      background: transparent;
      border: 0px;
      appearance: none;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $primary_text;
      height: 52px;
      caret-color: $primary_text;
      box-shadow: 0 0 0 0 transparent inset !important;
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    // background: rgba(0, 0, 0, 0.1);
    background: #fafafa;
    border-radius: 4px;
    color: #F5F7FA;
  }
}
</style>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;
$primary_text: #303133;

.login-container {
  min-height: 100%;
  width: 100%;
  background: $bg;
  background-size: cover;
  background-repeat: no-repeat;
  background-attachment: fixed;
  overflow: hidden;

  .login-form {
    position: relative;
    margin: 0 8% auto auto;
    overflow: hidden;
    // width: 25%;
    max-width: 450px;
    padding: calc(10% + 40px) 40px 30px;
    height: 100vh;
    background: rgba(255, 255, 255, 0.2);
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .err_info {
    display: block;
    color: red;
    font-size: 14px;
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
      padding: 0;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 14px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
