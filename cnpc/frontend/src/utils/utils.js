import Config from "../config/index";
import requests from "./requests";
import requestsv2 from "./requestsv2";
import { Loading, Message, MessageBox } from "element-ui";
import fileDownload from "js-file-download";
import axios from "axios";
import {
  getToken,
  getEcho,
  getUserInfo,
  getAddr,
  removeToken,
} from "../utils/dataStorage.js";
import https from "https";

const httpsAgent = new https.Agent({
  rejectUnauthorized: false,
});

export function versionCheck() {
  let usr = getUserInfo();
  if (!usr || usr.version != Config.version) {
    removeToken();
    window.location = "/#/login";
  }
}

export function currentUrlToParams(key = null) {
  let paramsUrl = window.location.href.split("?");
  if (paramsUrl.length < 2) return key ? null : {};
  let paramsArr = paramsUrl[1].split("&");
  let paramsData = {};
  paramsArr.forEach((r) => {
    let data = r.split("=");
    paramsData[data[0]] = data[1];
  });
  if (key)
    return Object.prototype.hasOwnProperty.call(paramsData, key)
      ? paramsData[key]
      : null;
  return paramsData;
}

export function clone(obj) {
  let strData = JSON.stringify(obj);
  return JSON.parse(strData);
}

export function isValidIP(ip) {
  return /^(?!0)(?!.*\.$)((1?\d?\d|25[0-5]|2[0-4]\d)(\.|$)){4}$/.test(ip);
}

export function isValidPort(port) {
  let valid = !isNaN(port);
  if (valid) {
    let p = parseInt(port);
    valid = p > 0 && p < 65536;
  }
  return valid;
}

export function allRequests(cb, gets = [], posts = [], puts = []) {
  let p = [],
    i = 0;
  for (i = 0; gets && i < gets.length; i++) {
    p.push(
      axios.get(gets[i].url, {
        params: gets[i].params,
        headers: {
          token: getToken(),
          Accept: "*/*",
          //'public-addr': getAddr()
        },
        baseURL: Config.apiUrl + "/" + Config.apiPrefix,
        httpsAgent: Config.devEnv ? httpsAgent : null,
      })
    );
  }
  for (i = 0; posts && i < posts.length; i++) {
    p.push(
      axios.post(posts[i].url, posts[i].data, {
        headers: {
          token: getToken(),
          Accept: "*/*",
          //'public-addr': getAddr()
        },
        baseURL: Config.apiUrl + "/" + Config.apiPrefix,
        httpsAgent: Config.devEnv ? httpsAgent : null,
      })
    );
  }
  for (i = 0; puts && i < puts.length; i++) {
    p.push(
      axios.put(puts[i].url, puts[i].data, {
        headers: {
          token: getToken(),
          Accept: "*/*",
          //"public-addr": getAddr(),
        },
        baseURL: Config.apiUrl + "/" + Config.apiPrefix,
        httpsAgent: Config.devEnv ? httpsAgent : null,
      })
    );
  }

  axios
    .all(p)
    .then((res) => {
      if (cb && cb.success) cb.success(res);

      if (cb && cb.finally) cb.finally(res);
    })
    .catch((errs) => {
      if (cb && cb.fail) cb.fail(errs);

      if (cb && cb.finally) cb.finally(errs);
    });
}

export function doRequest(param, callback) {
  if (param.url != "/login") versionCheck();
  if (param.loading && !window.loadingInstance) {
    window.loadingInstance = Loading.service({
      text: param.loadingText ? param.loadingText : "加载中...",
    });
  }

  requests(
    Object.assign({}, param, {
      url: Config.baseUrl + param.url,
    })
  )
    .then((res) => {
      if (callback) {
        if (
          Object.prototype.hasOwnProperty.call(callback, "obj") &&
          Object.prototype.hasOwnProperty.call(callback, "src") &&
          Object.prototype.hasOwnProperty.call(callback, "dst")
        ) {
          callback.obj[callback.dst] = res[callback.src];
        }

        if (callback.success) {
          callback.success(res);
        }

        if (callback.finally) {
          callback.finally(param.userData);
        }
      }
      if (param.loading) {
        setTimeout((_) => {
          if (window.loadingInstance) {
            window.loadingInstance.close();
            window.loadingInstance = undefined;
          }
        }, 200);
      }
    })
    .catch((err) => {
      if (callback) {
        if (callback.fail) callback.fail(err);

        if (callback.finally) {
          callback.finally(param.userData);
        }
      }

      if (param.loading) {
        setTimeout((_) => {
          if (window.loadingInstance) {
            window.loadingInstance.close();
            window.loadingInstance = undefined;
          }
        }, 200);
      }
    });
}

export function doRequestv2(param, callback) {
  if (param.loading && !window.loadingInstance) {
    window.loadingInstance = Loading.service({
      text: param.loadingText ? param.loadingText : "加载中...",
    });
  }

  requestsv2(
    Object.assign({}, param, {
      url: Config.baseUrl + param.url,
    })
  )
    .then((res) => {
      if (callback) {
        if (
          Object.prototype.hasOwnProperty.call(callback, "obj") &&
          Object.prototype.hasOwnProperty.call(callback, "src") &&
          Object.prototype.hasOwnProperty.call(callback, "dst")
        ) {
          callback.obj[callback.dst] = res.data[callback.src];
        }

        if (callback.success) {
          callback.success(res);
        }

        if (callback.finally) {
          callback.finally(param.userData);
        }
      }
      if (param.loading) {
        setTimeout((_) => {
          if (window.loadingInstance) {
            window.loadingInstance.close();
            window.loadingInstance = undefined;
          }
        }, 200);
      }
    })
    .catch((err) => {
      if (callback) {
        if (callback.fail) callback.fail(err);

        if (callback.finally) {
          callback.finally(param.userData);
        }
      }
      if (param.loading) {
        setTimeout((_) => {
          if (window.loadingInstance) {
            window.loadingInstance.close();
            window.loadingInstance = undefined;
          }
        }, 200);
      }
    });
}

export function download(param, callback) {
  requests(
    Object.assign({}, param, {
      url: Config.baseUrl + param.url,
    })
  )
    .then((res) => {
      const url = window.URL.createObjectURL(new Blob([res]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", callback.saveas); //or any other extension
      document.body.appendChild(link);
      link.click();
    })
    .catch((err) => {
      if (callback && callback.fail) {
        callback.fail(err);
      }
    });
}

export function downloadFile(param, callback) {
  requestsv2(
    Object.assign({}, param, {
      url: Config.baseUrl + param.url,
      responseType: "blob",
    })
  )
    .then((res) => {
      let fn = decodeURIComponent(
        escape(
          res.headers["content-disposition"].substring(
            "attachment;filename=".length
          )
        )
      );
      if (fn && res.data) {
        fileDownload(res.data, fn);
      }
      if (callback && callback.success) {
        callback.success(res);
      }

      if (callback && callback.finally) {
        callback.finally();
      }
    })
    .catch((err) => {
      if (callback && callback.fail) {
        callback.fail(err);
      }

      if (callback && callback.finally) {
        callback.finally();
      }
    });
}

export function beautifyNumber(num) {
  if (isNaN(num)) return num + "";

  let s = num + "";
  let symbol = "";
  if (s[0] === "-") {
    s = s.substr(1);
    symbol = "-";
  }
  let dot = s.indexOf(".");
  let int = "",
    deci = "";
  if (dot >= 0) {
    int = s.substring(0, dot);
    deci = s.substring(dot, s.length);
  } else {
    int = s;
  }
  dot = int.length - 1;
  let nint = int.length,
    res = "",
    c = 0;
  for (; dot >= 0; dot--) {
    res = int[dot] + res;
    c += 1;
    if (c === 3) {
      c = 0;
      if (dot !== 0) res = "," + res;
    }
  }

  return symbol + res + deci;
}

export function beautifyNumberChn(num, ignoreIfLess = 0) {
  if (isNaN(num)) return num + "";

  const yi = 100000000;
  const wan = 10000;
  const qian = 1000;
  const bai = 100;
  const shi = 10;

  if (num >= yi) {
    return (num / yi).toFixed(2) + "亿";
  } else if (num >= wan && wan >= ignoreIfLess) {
    return (num / wan).toFixed(2) + "万";
  } else if (num >= qian && qian >= ignoreIfLess) {
    return (num / qian).toFixed(2) + "千";
  } else if (num >= bai && bai >= ignoreIfLess) {
    return (num / bai).toFixed(2) + "百";
  } else if (num >= shi && shi >= ignoreIfLess) {
    return (num / shi).toFixed(2) + "十";
  }

  return num + "";
}

export function toggleKey(obj, key) {
  let suffix = "-toggle";
  let val = obj[key];
  if (val) {
    if (val.endsWith(suffix))
      obj[key] = val.substring(0, val.length - suffix.length);
    else obj[key] = val + suffix;
  }
}

export function deepCopy(obj) {
  if (obj) return JSON.parse(JSON.stringify(obj));
  else return null;
}

export function randomNumber(min, max, pricision = 0) {
  let rand = min + Math.random() * Math.floor(max - min);
  return parseFloat(rand.toFixed(pricision));
}

export function message(type, msg, showClose = true) {
  if (typeof type === "object") {
    Message(type);
  } else {
    Message({
      type: type,
      message: msg,
      showClose: showClose,
    });
  }
}

export function confirm(type, title, msg, cbyes, cbno) {
  MessageBox.confirm(msg, title, {
    type: type,
  })
    .then(cbyes ? cbyes : () => {})
    .catch(cbno ? cbno : () => {});
}

export function amIAdmin() {
  let u = getUserInfo();
  return u.role == 1;
  // return u == "admin";
}

export function syncSleep(ms) {
  let end = Date.now() + ms;
  while (Date.now() < end) continue;
}

export function mergeNavigation(navMain, navSub) {
  let nav = deepCopy(navMain);
  Object.keys(navSub).forEach((k) => {
    nav[k].show = navSub[k].show;
    let cMain = nav[k].children;
    let cSub = navSub[k].children;
    if (cMain && cSub) {
      cSub.forEach((s) => {
        cMain.forEach((m) => {
          if (m.key == s.key) {
            m.show = s.show;
          }
        });
      });
    }
  });
  return nav;
}
