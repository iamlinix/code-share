let navigation = {
  merch: {
    key: "nk-merch-effi",
    name: "非油商品效率",
    path: "/efficiency",
    icon: "fas fa-boxes",
  },
  orders: {
    key: "nk-orders",
    name: "订单分析",
    path: "/orders",
    icon: "fas fa-file-signature",
  },
  providers: {
    key: "nk-providers",
    name: "供应商贡献度",
    path: "/providers",
    icon: "fas fa-truck-moving",
  },
  prices: {
    key: "nk-prices",
    name: "价格带分析",
    path: "/prices",
    icon: "fas fa-align-left",
  },
  specs: {
    key: "nk-specs",
    name: "规格带分析",
    path: "/specs",
    icon: "fas fa-align-right",
  },
  skybook: {
    key: "nk-bible",
    name: "天书模型",
    path: "/bible",
    icon: "fas fa-bible",
  },
};

navigation = {
  configs: {
    key: "nk-configs",
    name: "系统配置",
    icon: "fas fa-users-cog",
    show: 1,
    children: [
      {
        key: "nk-users",
        name: "用户配置",
        path: "/users",
        icon: "fas fa-users",
        show: 1,
      },
    ],
  },
  // sales: {
  //     key: 'nk-sales',
  //     name: '销售跟踪',
  //     icon: 'fab fa-amazon-pay',
  //     children: [{
  //         key: 'nk-sales-kpi',
  //         name: '销售KPI',
  //         path: '/sales-kpi',
  //         icon: 'fas fa-chart-line'
  //     }]
  // },
  // prices: {
  //     key: 'nk-prices',
  //     name: '价格带分析',
  //     path: '/prices',
  //     icon: 'fas fa-align-left'
  // },
  // vendorsMain: {
  //     key: 'nk-vendors-main',
  //     name: '供应商分析',
  //     icon: 'fas fa-truck-loading',
  //     children: [
  //         {
  //             key: 'nk-vendors',
  //             name: '供应商配置',
  //             path: '/vendors',
  //             icon: 'fas fa-cogs'
  //     },
  //     {
  //             key: 'nk-vendor-mon-rpt',
  //             name: '进销存月度报表',
  //             path: '/vendor-monthly-report',
  //             icon: 'fas fa-scroll'
  //     }
  //     ]
  // },
  vendors: {
    key: "vendor-main",
    name: "供应商",
    icon: "fas fa-briefcase",
    show: 1,
    children: [
      {
        key: "nk-vendors-main",
        name: "供应商配置",
        path: "/vendors",
        icon: "fas fa-truck-loading",
        show: 1,
      },
      {
        key: "nk-bible-view",
        name: "查看天书",
        path: "/bible-view",
        icon: "fas fa-bible",
        show: 1,
      },
      {
        key: "nk-bible",
        name: "天书配置",
        path: "/bible",
        icon: "fas fa-bible",
        show: 1,
      },
      {
        key: "nk-vendor-mon-rpt",
        // name: "进销存月度报表",
        name: "供应商结算",
        path: "/vendor-monthly-report",
        icon: "fas fa-scroll",
        show: 1,
      },
      {
        key: "nk-vendor-kpi",
        name: "供应商考核",
        path: "/vendor-kpi",
        icon: "fas fa-sort-alpha-down",
        show: 1,
      },
    ],
  },
  sales: {
    key: "sales-main",
    name: "销售数据",
    icon: "fas fa-chart-pie",
    show: 1,
    children: [
      {
        key: "nk-kpi-settings",
        name: "销售目标设置",
        path: "/kpisettings",
        icon: "fas fa-sliders-h",
        show: 1,
      },
      {
        key: "nk-kpi-analyse",
        name: "销售指标分析",
        path: "/kpianalyse",
        icon: "fas fa-chart-line",
        show: 1,
      },
      {
        key: "nk-promotions",
        name: "销售数据统计",
        path: "/promotions",
        icon: "fas fa-ad",
        show: 1,
      },
      {
        key: "nk-daily-avg-view",
        name: "单店日均评级",
        path: "/daily-average-view",
        icon: "fas fa-sort-numeric-down",
        show: 1,
      },
      {
        key: "nk-daily-avg",
        name: "单店日均配置",
        path: "/daily-average",
        icon: "fas fa-sort-numeric-down",
        show: 1,
      },
      {
        key: "nk-no-sales",
        name: "未销售商品清单",
        path: "/no-sales",
        icon: "fas fa-box-open",
        show: 1,
      },
      {
        key: "nk-sales-focus",
        name: "关注商品",
        path: "/sales-focus",
        icon: "fas fa-crosshairs",
        show: 1,
      },
      {
        key: "nk-single-store",
        name: "单站数据",
        path: "/single-store",
        icon: "fas fa-gas-pump",
        show: 1,
      },
    ],
  },
  effis: {
    key: "effis-main",
    name: "商品效率",
    icon: "fas fa-file-invoice-dollar",
    show: 1,
    children: [
      //     {
      //     key: 'nk-effi-config',
      //     name: '商品效率参数',
      //     path: '/effi-config',
      //     icon: 'fas fa-cog'
      // },
      {
        key: "nk-effi-stale-goods",
        name: "滞销商品清单",
        path: "/effi-stale-goods",
        icon: "fas fa-heart-broken",
        show: 1,
      },
      {
        key: "nk-effi-stock-alarm",
        name: "库存效率预警",
        path: "/effi-stock-alarm",
        icon: "fas fa-bell",
        show: 1,
      },
      {
        key: "nk-effi-knockout",
        name: "淘汰商品清单",
        path: "/effi-knockouts",
        icon: "fas fa-cut",
        show: 1,
      },
    ],
  },
  // promitions: {
  //     key: 'nk-promotions',
  //     name: '销售数据统计',
  //     path: '/promotions',
  //     icon: 'fas fa-ad'
  // },
  // dailyAvg: {
  //     key: 'nk-daily-avg',
  //     name: '单店日均评级',
  //     path: '/daily-average',
  //     icon: 'fas fa-sort-numeric-down'
  // },
  // noSales: {
  //     key: 'nk-no-sales',
  //     name: '未销售商品清单',
  //     path: '/no-sales',
  //     icon: 'fas fa-box-open'
  // }
};

export default navigation;
