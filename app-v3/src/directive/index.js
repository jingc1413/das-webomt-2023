import hasRole from './permission/hasRole'
import hasPermissionOr from './permission/hasPermissionOr'
import hasPermissionAnd from './permission/hasPermissionAnd'
import hasAdmin from './permission/hasAdmin'


export default function directive(app){
  app.directive('hasRole', hasRole)
  app.directive('hasPermissionOr', hasPermissionOr)
  app.directive('hasPermissionAnd', hasPermissionAnd)
  app.directive('checkAdmin', hasAdmin)

}