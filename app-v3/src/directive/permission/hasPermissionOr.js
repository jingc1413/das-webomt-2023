 

import hasPermissions from '@/utils/hasPermission.js';

export default {
  mounted(el, binding, vnode) {
    const { value } = binding

    if (value && value.length > 0) {
      const permissions = value;

      let permissionFlag = false;
      if (typeof permissions == 'string') {
        permissionFlag = hasPermissions(permissions)
      } else {
        permissionFlag = permissions.some(item => {
          return hasPermissions(item)
        })
      }
      
      if (!permissionFlag) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    } else {
      // console.log('hasPermi mounted value', value)
      // throw new Error(`Set the value of the operation permission tag`)
    }
  }
}
