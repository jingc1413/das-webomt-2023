import settings from "@/settings";
import * as base from "./api-base";

export async function login(username, password) {
  const data = { Username: username, Password: password };
  return base.httpPost(base.ApiBase + "/auth/login", data);
}

export async function logout() {
  return base.httpGet(base.ApiBase + "/auth/logout");
}

export async function getCurrentUser() {
  if (settings.nodeTest) {
    return base.httpGet("/mock/auth/rules.json");
  }
  return base.httpGet(base.ApiBase + "/current");
}

export async function changeCurrentUserPassword(password, password2) {
  const data = { Password: password, NewPassword: password2 };
  return base.httpPost(base.ApiBase + "/current/change-password", data);
}