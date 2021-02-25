import Vue from "vue";
import App from "./App.vue";
import ElementUI from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import locale from "element-ui/lib/locale/lang/en";
import "./assets/fonts.css";
//import axios from "axios";

//axios.defaults.baseURL = "http://192.168.31.183:10086";

Vue.use(ElementUI, { locale });

Vue.config.productionTip = false;

new Vue({
  render: (h) => h(App),
}).$mount("#app");
