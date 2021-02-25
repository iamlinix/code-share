import Vue from "vue";
import Router from "vue-router";
import { getToken, isLogin, isEchoed, setEcho } from "@/utils/dataStorage";
import Config from "../config/index";
import { doRequest } from "@/utils/utils";

Vue.use(Router);

let RouteList = [
  {
    path: "/",
    component: (resolve) => require(["@/views/layout/App.vue"], resolve),
    meta: {
      title: "",
      keepAlive: false,
      requiresAuth: true,
    },
    children: [
      {
        path: "/efficiency",
        name: "Efficiency",
        meta: {
          title: "非油商品效率",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/effi/Index.vue"], resolve),
      },
      {
        path: "/orders",
        name: "Orders",
        meta: {
          title: "订单分析",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/orders/Index.vue"], resolve),
      },
      {
        path: "/providers",
        name: "Providers",
        meta: {
          title: "供应商贡献度",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/providers/Index.vue"], resolve),
      },
      {
        path: "/prices",
        name: "Prices",
        meta: {
          title: "价格带分析",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/prices/Index.vue"], resolve),
      },
      {
        path: "/specs",
        name: "Specs",
        meta: {
          title: "规格带分析",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/specs/Index.vue"], resolve),
      },
      {
        path: "/bible-view",
        name: "BibleView",
        meta: {
          title: "查看天书",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/bible/ViewOnly.vue"], resolve),
      },
      {
        path: "/bible",
        name: "Bible",
        meta: {
          title: "天书配置",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/bible/Index.vue"], resolve),
      },
      {
        path: "/users",
        name: "Users",
        meta: {
          title: "用户配置",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/users/Index.vue"], resolve),
      },
      {
        path: "/vendors",
        name: "Vendor",
        meta: {
          title: "供应商统计",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/vendors/Index.vue"], resolve),
      },
      {
        path: "/kpisettings",
        name: "KpiSettings",
        meta: {
          title: "销售目标设置",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/kpi/Settings.vue"], resolve),
      },
      {
        path: "/kpianalyse",
        name: "KpiAnalyse",
        meta: {
          title: "销售指标分析",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/kpi/Analyse.vue"], resolve),
      },
      {
        path: "/promotions",
        name: "Promotions",
        meta: {
          title: "营销活动分析",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/promotions/Index3.vue"], resolve),
      },
      {
        path: "/vendor-monthly-report",
        name: "MonthlyReport",
        meta: {
          // title: "进销存月度报表",
          title: "供应商结算",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/vendors/MonthlyReport.vue"], resolve),
        // require(["@/views/vendors/Test.vue"], resolve),
      },
      {
        path: "/daily-average-view",
        name: "DailyAverageView",
        meta: {
          title: "单店日均评级",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/dailyAvg/ViewOnly.vue"], resolve),
        // require(["@/views/dailyAvg/Index1.vue"], resolve),
      },
      {
        path: "/daily-average",
        name: "DailyAverage",
        meta: {
          title: "单店日均配置",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/dailyAvg/Index.vue"], resolve),
        // require(["@/views/dailyAvg/Index1.vue"], resolve),
      },
      {
        path: "/no-sales",
        name: "NoSales",
        meta: {
          title: "未销售商品清单",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/nosales/Index.vue"], resolve),
      },
      {
        path: "/sales-kpi",
        name: "KPI",
        meta: {
          title: "销售KPI",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/sales/Kpi.vue"], resolve),
      },
      {
        path: "/sales-focus",
        name: "关注商品",
        meta: {
          title: "关注商品",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/focus/Index.vue"], resolve),
      },
      {
        path: "/vendor-kpi",
        name: "供应商考核",
        meta: {
          title: "供应商考核",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/vendors/Kpi.vue"], resolve),
      },
      {
        path: "/effi-config",
        name: "商品效率参数",
        meta: {
          title: "商品效率参数",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/effi/Index.vue"], resolve),
      },
      {
        path: "/effi-stale-goods",
        name: "滞销商品清单",
        meta: {
          title: "滞销商品清单",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/effi/Unsalable.vue"], resolve),
      },
      {
        path: "/effi-stock-alarm",
        name: "库存效率预警",
        meta: {
          title: "库存效率预警",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/effi/StockAlarm.vue"], resolve),
      },
      {
        path: "/effi-knockouts",
        name: "淘汰商品清单",
        meta: {
          title: "淘汰商品清单",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) => require(["@/views/effi/Knockout.vue"], resolve),
      },
      {
        path: "/single-store",
        name: "单站数据",
        meta: {
          title: "单站数据",
          keepAlive: false,
          requiresAuth: true,
        },
        component: (resolve) =>
          require(["@/views/therest/Singlestore.vue"], resolve),
      },
    ],
  },
  {
    path: "/login",
    name: "Login",
    meta: {
      title: "用户登陆",
      keepAlive: false,
    },
    components: {
      blank: (resolve) => require(["@/views/login/Index.vue"], resolve),
    },
  },
];

const router = new Router({ routes: RouteList });

router.beforeEach((to, from, next) => {
  if (to.matched.some((r) => r.meta.requiresAuth) && !isLogin()) {
    next({
      path: "/login",
      query: {
        redirect: to.fullPath,
      },
    });
  } else {
    if (to.path === "/") {
      next({ path: Config.firstPage });
    } else {
      next();
    }
  }
});

export default router;
