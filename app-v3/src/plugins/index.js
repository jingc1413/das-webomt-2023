import permissions from './permissions'
import modal from './modal'
export default function installPlugins(app){

  // Authentication
  app.config.globalProperties.$permissions = permissions;
  app.config.globalProperties.$msgModal = modal;

  // Authentication
  app.provide('$permissions', permissions)

}
