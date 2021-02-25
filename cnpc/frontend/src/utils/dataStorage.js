import Storage from "good-storage";
import Config from "../config/index.js";
import CryptoJS from "crypto-js";

export function getUserInfo() {
  let obj = null;
  try {
    obj = JSON.parse(
      CryptoJS.AES.decrypt(
        Storage.get(Config.userInfoKey),
        Config.cipherKey
      ).toString(CryptoJS.enc.Utf8)
    );
  } catch (error) {
    console.log(error);
  }
  return obj;
}

export function setUserInfo(user) {
  Storage.set(
    Config.userInfoKey,
    CryptoJS.AES.encrypt(user, Config.cipherKey).toString()
  );
  return user;
}

export function setUserID(id) {
  Storage.set(Config.userIDKey, id);
}

export function getUserID() {
  return Storage.get(Config.userIDKey) || -1;
}

export function getToken() {
  let token = Storage.get(Config.tokenKey);
  if (token)
    return CryptoJS.AES.decrypt(token, Config.cipherKey).toString(
      CryptoJS.enc.Utf8
    );
  else return "";
}

export function getAddr() {
  let addr = Storage.get(Config.addrKey);
  if (addr)
    return CryptoJS.AES.decrypt(addr, Config.cipherKey).toString(
      CryptoJS.enc.Utf8
    );
  else return "0.0.0.0";
}

export function setAddr(addr) {
  Storage.set(
    Config.addrKey,
    CryptoJS.AES.encrypt(addr, Config.cipherKey).toString()
  );
}

export function setToken(token) {
  return Storage.set(
    Config.tokenKey,
    CryptoJS.AES.encrypt(token, Config.cipherKey).toString()
  );
  //return Cookies.set(Config.tokenKey, token ,{ expires: Config.cookiesExpires })
}

export function removeToken() {
  return Storage.remove(Config.tokenKey);
}

export function isLogin() {
  return (getToken() || "").length > 5;
}

export function isEchoed() {
  return (Storage.get(Config.echoKey) || "").length > 8;
}

export function setEcho(echo) {
  Storage.set(Config.echoKey, echo);
}

export function getEcho() {
  return Storage.get(Config.echoKey);
}
