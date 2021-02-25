import Storage from "good-storage";

export function getUserInfo() {
  var infoStr = Storage.get("user-info", "");
  if (infoStr.length > 0) {
    return JSON.parse(infoStr);
  }
  return null;
}

export function setUserInfo(userInfo) {
  Storage.set("user-info", JSON.stringify(userInfo));
  Storage.set("token", userInfo.token);
}

export function removeUserInfo() {
  Storage.remove("user-info");
  Storage.remove("token");
}

export function getToken() {
  return Storage.get("token");
}

export function setToken(token) {
  Storage.set("token", token);
}
