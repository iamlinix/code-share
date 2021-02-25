import Vue from "vue";
import "xe-utils";
import App from "./App.vue";
import ElementUI from "element-ui";
import Config from "./config/index";
import "./assets/css/style.scss";
import router from "./router/";
import { isLogin } from "./utils/dataStorage";
import store from "./store/";

import "@fortawesome/fontawesome-free/css/all.css";
import "@fortawesome/fontawesome-free/js/all.js";
import Donut from "vue-css-donut-chart";
import "vue-css-donut-chart/dist/vcdonut.css";

import { library } from "@fortawesome/fontawesome-svg-core";
import {
  faLongArrowAltDown,
  faLongArrowAltUp,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faLongArrowAltDown);
library.add(faLongArrowAltUp);

Vue.config.productionTip = false;
Vue.prototype.$Config = Config;

import VueTippy, { TippyComponent } from "vue-tippy";
import "@/assets/fonts/roboto.all.css";
import VXETable from "vxe-table";
import "vxe-table/lib/style.css";
import VXETablePluginExportXLSX from "vxe-table-plugin-export-xlsx";

Vue.use(VXETable);
VXETable.use(VXETablePluginExportXLSX);
Vue.use(VueTippy);
Vue.component("tippy", TippyComponent);
Vue.use(ElementUI);
Vue.use(Donut);
Vue.component("font-awesome-icon", FontAwesomeIcon);

router.beforeEach((to, from, next) => {
  window.document.title = to.meta.title
    ? to.meta.title + "-" + Config.siteName
    : Config.siteName;

  if (!isLogin() && to.path !== "/login") {
    next({ path: "/login" });
  } else {
    next();
  }
});

router.afterEach((transition) => {});

Config.router = router;

new Vue({
  el: "#app",
  router,
  store,
  render: (h) => h(App),
});
