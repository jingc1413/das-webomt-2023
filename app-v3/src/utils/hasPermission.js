 
import {useAuthStore} from '@/stores/auth'

export default function (permission) {
  if (typeof permission == 'string') {
    return hasPermissionsString(permission);
  } else {
    return hasPermissionsArray(permission);
  }
}

function hasPermissionsArray(permissionFlag) {
  const hasPermissions = permissionFlag.every(permission => {
    return hasPermissionsString(permission);
  })
  return hasPermissions
}

function hasPermissionsString(permissionFlag) {
  const permissions = useAuthStore().permissions;
  const extensionsPermissions = useAuthStore().extensionsPermissions;
  let flag = false;
  flag = extensionsPermissions.some(item=>{
    return new RegExp(item).test(permissionFlag)
  })
  if (flag) {
    return true;
  }
  return permissions.includes(permissionFlag);
}

