import authPermission from '@/utils/hasPermission.js';
import {useAuthStore} from '@/stores/auth'

function authRole(role) {
  const roles = useAuthStore().loginUserRoles
  if (role && role.length > 0) {
    return roles.some(v => {
      return v == role
    })
  } else {
    return false
  }
}

export default {
  // Verify that the user has a permission
  hasPermission(permission) {
    return authPermission(permission);
  },
  hasPermissionOr(permissions) {
    if (typeof permissions == 'string') {
      return authPermission(permissions)
    } else {
      return permissions.some(item => {
        return authPermission(item)
      })
    }
  },
  hasPermissionAndAll(permissions) {
    if (typeof permissions == 'string') {
      return authPermission(permissions)
    } else {
      return permissions.every(item => {
        return authPermission(item)
      })
    }
  },
  // Verify that the user has a role
  hasRole(role) {
    return authRole(role);
  },
  hasRoleOr(roles) {
    return roles.some(item => {
      return authRole(item)
    })
  },
  hasRoleAnd(roles) {
    return roles.every(item => {
      return authRole(item)
    })
  }
}
