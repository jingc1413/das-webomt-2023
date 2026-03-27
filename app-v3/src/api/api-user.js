import * as base from "./api-base";

export async function getRoleList() {
  return base.httpGet(base.ApiBase + "/iam/roles");
}

export async function getUserList() {
  return base.httpGet(base.ApiBase + "/iam/users");
}

export async function createUser(user) {
  return base.httpPost(base.ApiBase + "/iam/users", user);
}

export async function deleteUser(user) {
  return base.httpDelete(base.ApiBase + "/iam/users/" + user.Name);
}

export async function editUser(user) {
  return base.httpPost(base.ApiBase + "/iam/users/" + user.Name, user);
}
