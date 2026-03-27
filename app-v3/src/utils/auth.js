import Cookies from 'js-cookie'

//登录状态
const TokenKey = 'login_state'
//登录角色
const TokenRole = 'login_role'
const webLanguage = 'webomt_language'
const superRoot = "superRoot"
const SessionKey = "session"

export function clearSession() {
  return Cookies.remove(SessionKey, "")
}

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getRole() {
  return Cookies.get(TokenRole)
}

export function setRole(token) {
  return Cookies.set(TokenRole, token)
}

export function removeRole() {
  return Cookies.remove(TokenRole)
}
export function clearToken() {
  removeToken()
  removeRole()
}
export function getSuperRoot() {
  //return sessionStorage.getItem(superRoot);
  return Cookies.get(superRoot)
}
export function setSuperRoot(token) {
  return Cookies.set(superRoot, token)
}
export function removeSuperRoot() {
  return Cookies.remove(superRoot)
}

export function getLanguage() {
  let lang = Cookies.get(webLanguage)
  if(lang == 'en' || lang == 'zh'){
    return lang
  }else{
    setLanguage('en')
    return 'en'
  }
}
export function setLanguage(token) {
  return Cookies.set(webLanguage, token)
}
export function removeLanguage() {
  return Cookies.remove(webLanguage)
}
