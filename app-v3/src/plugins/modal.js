import { ElMessage, ElMessageBox, ElNotification, ElLoading } from 'element-plus'
import { translator } from '../i18n';
let loadingInstance;

export default {
  /**
   * ElMessage.info
   */
  msg(content) {
    ElMessage.info(content)
  },
  /**
   * ElMessage.info
   */
  msgError(content) {
    ElMessage.error(content)
  },
  /**
   * ElMessage.info
   */
  msgSuccess(content) {
    ElMessage.success(content)
  },
  /**
   * ElMessage.info
   */
  msgWarning(content) {
    ElMessage.warning(content)
  },
  // ElMessageBox.alert
  alert(content) {
    ElMessageBox.alert(content, translator("tip.tip"))
  },
  // ElMessageBox.alert error
  alertError(content) {
    ElMessageBox.alert(content, translator("tip.tip"), { type: 'error' })
  },
  // ElMessageBox.alert success
  alertSuccess({content, callback}) {
    ElMessageBox.alert(content, translator("tip.tip"), { type: 'success', callback })
  },
  // ElMessageBox.alert warning
  alertWarning(content) {
    ElMessageBox.alert(content, translator("tip.tip"), { type: 'warning' })
  },
  // ElNotification.info
  notify(content) {
    ElNotification.info(content)
  },
  // ElNotification.error
  notifyError(content) {
    ElNotification.error(content);
  },
  // ElNotification.success
  notifySuccess(content) {
    ElNotification.success(content)
  },
  // ElNotification.warning
  notifyWarning(content) {
    ElNotification.warning(content)
  },
  // ElMessageBox.confirm
  confirm(content, title='Tip', options) {
    let confirmOptions = {
      confirmButtonText: translator("buttons.confirm"),
      cancelButtonText: translator("buttons.cancel"),
      type: "warning",
      ...options
    }
    return ElMessageBox.confirm(content, title, confirmOptions)
  },
  /**
   * ElMessageBox.prompt
   */
  prompt(content) {
    return ElMessageBox.prompt(content, translator("tip.tip"), {
      confirmButtonText: translator("buttons.confirm"),
      cancelButtonText: translator("buttons.cancel"),
      type: "warning",
    })
  },
  /**
   * open loading
   */
  loading(content) {
    loadingInstance = ElLoading.service({
      lock: true,
      text: content,
      background: "rgba(0, 0, 0, 0.7)",
    })
  },
  /**
   * closeLoading
   */
  closeLoading() {
    loadingInstance.close();
  }
}
