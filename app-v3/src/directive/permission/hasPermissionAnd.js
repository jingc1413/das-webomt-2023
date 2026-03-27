import hasPermissions from '@/utils/hasPermission.js';

export default {
  mounted(el, binding, vnode) {
    const { value } = binding
    if (value && value.length > 0) {
      const permissions = value;

      const permissionFlag = permissions.every(permission => {
        return hasPermissions(permission);
      })
      if (!permissionFlag) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      //  console.log('hasPermissionAnd mounted value', value)
      //  throw new Error(`Set the value of the operation permission tag`)
    }
  }
}
