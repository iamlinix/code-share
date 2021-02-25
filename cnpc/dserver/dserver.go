package main

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/handlers"
	"cnpc.com.cn/cnpc/dserver/handlers/basic"
	"cnpc.com.cn/cnpc/dserver/handlers/eff"
	"cnpc.com.cn/cnpc/dserver/handlers/obs"
	"cnpc.com.cn/cnpc/dserver/handlers/sales"
	"cnpc.com.cn/cnpc/dserver/handlers/sys"
	"cnpc.com.cn/cnpc/dserver/handlers/vendor"
	"cnpc.com.cn/cnpc/dserver/middleware/auth"
	"cnpc.com.cn/cnpc/dserver/orm"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"github.com/gin-gonic/contrib/static"
	//"github.com/gin-contrib/cors"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var confPath = flag.String("conf-dir", "./conf", "configuration directory")

func HiddenCallerEncoder(_ zapcore.EntryCaller,
	enc zapcore.PrimitiveArrayEncoder) {

	enc.AppendString("hidden")
}

func initConfig() {

	/* init config file */
	viper.SetConfigName("dserver")
	viper.AddConfigPath(*confPath)
	switch runtime.GOOS {
	case "windows":
		viper.AddConfigPath("F:/Projects/cnpc/dserver/conf")
	case "linux":
		viper.AddConfigPath("/opt/cnpc/conf")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("failed read configurations, error: ", err.Error())
		os.Exit(-5)
	}

	/* config watchdog */
	viper.WatchConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})
}

func initLogger() {

	cfg := zap.NewProductionConfig()
	logPath := viper.GetString(common.ConfKeyLogPath)
	cfg.OutputPaths = []string{
		logPath,
	}
	cfg.ErrorOutputPaths = []string{
		logPath,
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeCaller = HiddenCallerEncoder
	level := viper.GetString(common.ConfKeyLogLevel)
	cfg.Level = zap.NewAtomicLevelAt(common.LogLevelMapping[level])

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to init logger: %s", err.Error())
		os.Exit(-2)
	}

	zap.ReplaceGlobals(logger)
}

func init() {

	flag.Parse()

	initConfig()

	initLogger()
	zaps.InitZaps()

	handlers.GinInit()
}

func StartTimer(f func()) {

	zaps.Info(">>> dserver start new timer")
	go func() {
		for {
			f()

			now := time.Now()
			next := now.Add(time.Second * 30)
			t := time.NewTimer(next.Sub(now))

			<-t.C
		}
	}()
	zaps.Info("<<< dserver new timer done")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers",
			"Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,Content-Disposition,Local-Addr,Public-Addr")
		c.Header("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, Content-Type, Content-Disposition")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func main() {

	var err error

	/* read from config file first */
	port := viper.GetInt(common.ConfKeyServerPort)
	port2 := viper.GetInt(common.ConfKeyServerPort2)
	bind := viper.GetString(common.ConfKeyServerBind)
	ssl := viper.GetBool(common.ConfKeyServerSSL)
	cert := viper.GetString(common.ConfKeyServerCert)
	key := viper.GetString(common.ConfKeyServerKey)

	webToken := viper.GetString(common.ConfKeyServerAuthWeb)
	zaps.Info("dserver need token: ", webToken)

	if port < 1024 || port > 65535 {
		zaps.Fatal("port range should be within [1024, 65535]")
		return
	}

	address := fmt.Sprintf("%s:%d", bind, port)
	address2 := fmt.Sprintf("%s:%d", bind, port2) //non-ssl

	zaps.Info("dserver listening on；", address)
	zaps.Info("dserver listening on；", address2)

	/* connect database */
	if viper.GetString(common.ConfKeyDBDriver) == "mysql" {
		err = orm.DBConnect(viper.GetString(common.ConfKeyDBDriver),
			viper.GetString(common.ConfKeyDBAddr),
			viper.GetString(common.ConfKeyDBPort),
			viper.GetString(common.ConfKeyDBUser),
			viper.GetString(common.ConfKeyDBPass),
			viper.GetString(common.ConfKeyDBName))
	} else {
		zaps.Fatal("umsupport database drive")
	}

	if err != nil {
		zaps.Fatal("failed to initialize database, error: ", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	//r.Use(auth.IPCheck())

	//r.Use(cors.Default())
	r.Use(Cors())

	path := viper.GetString(common.ConfKeyDistPath)
	r.Use(static.Serve("/", static.LocalFile(path, true)))

	r.POST("/login", sys.LoginHandler)

	/* API for WEB */
	web := r.Group("/v1/web")
	if webToken == "true" {
		web.Use(auth.JWTAuth())
	}
	{
		////////////////////////////////////////////////////////////////
		// BASIC API
		/* basic API */
		web.GET("/basic/material/:code", basic.BasicMaterialHandler)
		web.GET("/basic/vendor/:code", basic.BasicVendorHandler)
		web.GET("/basic/plant/:code", basic.BasicPlantHandler)

		/* class API */
		web.GET("/basic/class/init", basic.BasicClassInitHandler)
		web.GET("/basic/class-dict/init", basic.BasicClassDictInitHandler)

		web.GET("/class/main", basic.MainClassHandler)
		web.GET("/class/mid", basic.MidClassHandler)
		web.GET("/class/mid/:mainclass", basic.MidClassByMainHandler)
		web.GET("/class/sub", basic.SubClassHandler)
		web.GET("/class/sub/:midclass", basic.SubClassByMidHandler)

		/* org API */
		web.GET("/org/branch", basic.BranchOrgHandler)
		web.GET("/org/plant", basic.PlantOrgHandler)

		////////////////////////////////////////////////////////////////
		// OBSOLETE API
		/* scenes API */
		web.POST("/scenes/add", obs.ScenesAddHandler)
		web.POST("/scenes/update", obs.ScenesUpdateHandler)
		web.GET("/scenes/del/:id", obs.ScenesDelHandler)
		web.GET("/scenes/list", obs.ScenesListHandler)
		web.GET("/scenes/detail/:name", obs.ScenesDetailHandler)

		/* price zone cfg API */
		web.POST("/pricezone/config/add", obs.PriceZoneCfgAddHandler)
		web.POST("/pricezone/config/update", obs.PriceZoneCfgUpdateHandler)
		web.GET("/pricezone/config/del/:classcode", obs.PriceZoneCfgDelHandler)
		web.GET("/pricezone/config/list", obs.PriceZoneCfgListHandler)
		web.GET("/pricezone/config/detail/:classcode", obs.PriceZoneCfgDetailHandler)

		/* price zone API */
		web.POST("/pricezone", obs.PriceZoneHandler)
		web.POST("/pricezone/single", obs.PriceZoneSingleHandler)
		web.POST("/pricezone/rank", obs.PriceZoneRankHandler)

		/* specs zone API */
		web.POST("/specszone", obs.SpecsZoneHandler)
		web.POST("/specszone/rank/:zone", obs.SpecsZoneRankHandler)

		////////////////////////////////////////////////////////////////
		// EFFICIENCY API
		/* Metric Threshold Config API */
		web.POST("/eff/metric/config/update", eff.EffMetricCfgUpdateHandler)
		web.GET("/eff/metric/config/list/:id", eff.EffMetricCfgListHandler)

		/* Sales Activity Rate Material API */
		web.POST("/eff/unsalable/material", eff.EffUnsalableMatlHandler)
		web.POST("/eff/unsalable/material/export", eff.EffUnsalableMatlExportHandler)

		/* Inventory Efficiency Material API */
		web.POST("/eff/inventory/material", eff.EffInventoryMatlHandler)
		web.POST("/eff/inventory/material/export", eff.EffInventoryMatlExportHandler)

		/* Test Market Material Config API */
		web.POST("/eff/test-market/material/config/add", eff.EffTestMarketMatlCfgAddHandler)
		web.POST("/eff/test-market/material/config/update", eff.EffTestMarketMatlCfgUpdateHandler)
		web.GET("/eff/test-market/material/config/del/:code", eff.EffTestMarketMatlCfgDelHandler)
		web.GET("/eff/test-market/material/config/clear", eff.EffTestMarketMatlCfgClearHandler)
		web.GET("/eff/test-market/material/config/list", eff.EffTestMarketMatlCfgListHandler)
		web.POST("/eff/test-market/matrial/config/import", eff.EffTestMarketMatlCfgImportHandler)
		web.GET("/eff/test-market/matrial/config/export", eff.EffTestMarketMatlCfgExportHandler)

		/* Test Market Material API */
		web.POST("/eff/test-market/material/fodate", eff.EffTestMarketMatlFODateHandler)
		web.POST("/eff/test-market/material", eff.EffTestMarketMatlHandler)
		web.POST("/eff/test-market/material/single", eff.EffTestMarketMatlSingleHandler)
		web.POST("/eff/test-market/material/export", eff.EffTestMarketMatlExportHandler)

		////////////////////////////////////////////////////////////////
		// VENDOR API
		/* bible API */
		web.POST("/erp/bible/add", vendor.VendorBibleCfgAddHandler)
		web.POST("/erp/bible/update", vendor.VendorBibleCfgUpdateHandler)
		web.GET("/erp/bible/del/:code", vendor.VendorBibleCfgDelHandler)
		web.GET("/erp/bible/clear", vendor.VendorBibleCfgClearHandler)
		web.GET("/erp/bible/list", vendor.VendorBibleCfgListHandler)
		web.GET("/erp/bible/detail/:code", vendor.VendorBibleCfgDetailHandler)
		web.POST("/erp/bible/import", vendor.VendorBibleCfgImportHandler)
		web.GET("/erp/bible/export", vendor.VendorBibleCfgExportHandler)

		/* vendor config */
		web.POST("/erp/vendor/config/add", vendor.ERPVendorCfgAddHandler)
		web.POST("/erp/vendor/config/update", vendor.ERPVendorCfgUpdateHandler)
		web.GET("/erp/vendor/config/del/:code", vendor.ERPVendorCfgDelHandler)
		web.GET("/erp/vendor/config/clear", vendor.ERPVendorCfgClearHandler)
		web.GET("/erp/vendor/config/list", vendor.ERPVendorCfgListHandler)
		web.GET("/erp/vendor/config/detail/:code", vendor.ERPVendorCfgDetailHandler)
		web.GET("/erp/vendor/config/export", vendor.ERPVendorCfgExportHandler)
		web.POST("/erp/vendor/config/import", vendor.ERPVendorCfgImportHandler)

		/* vendor exclude material config */
		web.POST("/erp/vendor/ex-matl/config/add", vendor.ERPExMatlCfgAddHandler)
		web.POST("/erp/vendor/ex-matl/config/update", vendor.ERPExMatlCfgUpdateHandler)
		web.GET("/erp/vendor/ex-matl/config/del/:code", vendor.ERPExMatlCfgDelHandler)
		web.GET("/erp/vendor/ex-matl/config/clear", vendor.ERPExMatlCfgClearHandler)
		web.GET("/erp/vendor/ex-matl/config/list", vendor.ERPExMatlCfgListHandler)
		web.GET("/erp/vendor/ex-matl/config/export", vendor.ERPExMatlCfgExportHandler)
		web.POST("/erp/vendor/ex-matl/config/import", vendor.ERPExMatlCfgImportHandler)

		/* vendor: purchase-sales-inv */
		web.POST("/erp/vendors/report/add", vendor.ERPVendorReportAddHandler)
		web.POST("/erp/vendors/report/update", vendor.ERPVendorReportUpdateHandler)
		web.GET("/erp/vendors/report/history/:type/:month/:vendor", vendor.ERPVendorReportHistory)
		web.GET("/erp/vendors/report/del/:month", vendor.ERPVendorReportDelHandler)
		web.GET("/erp/vendors/report/list", vendor.ERPVendorReportListHandler)
		web.GET("/erp/vendors/report/detail/:month/all", vendor.ERPVendorReportDetailAllHandler)
		web.GET("/erp/vendors/report/detail/:month/active", vendor.ERPVendorReportDetailActiveHandler)
		web.POST("/erp/vendors/report/export", vendor.ERPVendorReportExportHandler)

		web.GET("/erp/materials/report/del/:month", vendor.ERPMaterialReportDelHandler)
		web.GET("/erp/materials/report/list/:month/all", vendor.ERPMaterialReportListByMonthAllHandler)
		web.GET("/erp/materials/report/list/:month/active", vendor.ERPMaterialReportListByMonthActiveHandler)
		web.GET("/erp/materials/report/detail/:month/:vendor", vendor.ERPMaterialReportListByMonthVendorHandler)
		web.GET("/erp/materials/report/export/:month/all", vendor.ERPMaterialReportExportAllHandler)
		web.GET("/erp/materials/report/export/:month/active", vendor.ERPMaterialReportExportActiveHandler)

		web.GET("/erp/vendors/report/account-statement/:begin-date/:end-date", vendor.ERPVendorAccountStatementExportHandler)
		web.GET("/erp/vendors/report/payment-details/:month", vendor.ERPVendorPaymentDetailsExportHandler)
		/*签认单*/
		web.GET("/erp/vendors/report/endorsement-form/:begin-date/:end-date", vendor.ERPVendorEndorsementFormExportHandler)
		//台账
		web.GET("/erp/vendors/report/standing-book/:month", vendor.ERPVendorStandingBookExportHandler)

		web.POST("/erp/materials/purchase", vendor.ERPMatlPurchaseHandler)
		web.POST("/erp/materials/zifpurd", vendor.ERPMatlZifpurdHandler)
		web.POST("/erp/materials/sales", vendor.ERPMatlSalesHandler)
		web.POST("/erp/materials/openstock", vendor.ERPMatlOpenInvHandler)
		web.POST("/erp/materials/closestock", vendor.ERPMatlCloseInvHandler)
		web.POST("/erp/materials", vendor.ERPMaterialHandler)
		web.POST("/erp/materials/export", vendor.ERPMatlExportHandler)

		web.POST("/erp/vendors", vendor.ERPVendorHandler)
		web.POST("/erp/vendors/export", vendor.ERPVendorExportHandler)

		/* vendor rating config */
		web.POST("/vendor/rating/config/update", vendor.VendorRatingCfgUpdateHandler)
		web.GET("/vendor/rating/config/list", vendor.VendorRatingCfgListHandler)

		/* vendor rating */
		web.POST("/vendor/rating/report", vendor.VendorRatingReportHandler)
		web.POST("/vendor/rating/report/add", vendor.VendorRatingReportAddHandler)
		//web.POST("/vendor/rating/report/update", vendor.VendorRatingReportUpdateHandler)
		web.GET("/vendor/rating/report/del/:id", vendor.VendorRatingReportDelHandler)
		web.GET("/vendor/rating/report/list", vendor.VendorRatingReportListHandler)
		web.GET("/vendor/rating/report/detail/:id", vendor.VendorRatingReportDetailHandler)
		web.GET("/vendor/rating/report/export/:id", vendor.VendorRatingReportExportHandler)

		////////////////////////////////////////////////////////////////
		// SALES API
		/* sales kpi */
		web.POST("/sales/kpi", sales.SalesKPIHandler)
		web.POST("/sales/class", sales.SalesClassHandler)
		web.POST("/sales/plant/rank", sales.SalesPlantRankHandler)
		web.POST("/sales/material/rank", sales.SalesMaterialRankHandler)
		web.POST("/sales/kpi/analyse", sales.SalesKPIAnalyseHandler)

		/* sales plant level cfg API */
		//XXX TODO
		web.POST("/plant/sales-zone/config/add", sales.SalesPlantLevelCfgAddHandler)
		//web.POST("/sales/plant/level/config/add", sales.SalesPlantLevelCfgAddHandler)
		web.POST("/plant/sales-zone/config/update", sales.SalesPlantLevelCfgUpdateHandler)
		//web.POST("/sales/plant/level/config/update", sales.SalesPlantLevelCfgUpdateHandler)
		web.GET("/plant/sales-zone/config/del/:id", sales.SalesPlantLevelCfgDelHandler)
		//web.GET("/sales/plant/level/config/del/:id", sales.SalesPlantLevelCfgDelHandler)
		web.GET("/plant/sales-zone/config/list", sales.SalesPlantLevelCfgListHandler)
		//web.GET("/sales/plant/level/config/list", sales.SalesPlantLevelCfgListHandler)
		web.GET("/plant/sales-zone/config/detail/:id", sales.SalesPlantLevelCfgDetailHandler)
		//web.GET("/sales/plant/level/config/detail/:id", sales.SalesPlantLevelCfgDetailHandler)

		/* sales plant level API */
		web.POST("/plant/sales-zone", sales.SalesPlantLevelHandler)
		web.POST("/plant/single-sales-zone", sales.SalesSinglePlantLevelHandler)
		//web.POST("/sales/plant/level", sales.SalesPlantLevelHandler)

		/* sales plant car-service API */
		web.POST("/plant/car-service", sales.SalesPlantCarServiceHandler)
		//web.POST("/sales/plant/car-service", sales.SalesPlantCarServiceHandler)

		/* sales plant config api */
		web.POST("/sales/plant/config/add", sales.SalesPlantCfgAddHandler)
		web.POST("/sales/plant/config/update", sales.SalesPlantCfgUpdateHandler)
		web.GET("/sales/plant/config/del/:code", sales.SalesPlantCfgDelHandler)
		web.GET("/sales/plant/config/clear", sales.SalesPlantCfgClearHandler)
		web.GET("/sales/plant/config/list", sales.SalesPlantCfgListHandler)

		/* sales plant api */
		web.POST("/sales/plant/kpi", sales.SalesPlantKPIHandler)
		web.POST("/sales/plant/material/rank", sales.SalesPlantMaterialRankHandler)
		web.POST("/sales/plant/date", sales.SalesPlantDateHandler)
		web.POST("/sales/plant/class", sales.SalesPlantClassHandler)

		//XXX TODO
		/* sales material plant no-sales cfg API */
		web.POST("/plant/material/monitor/config/add", sales.SalesMatlNoSalesCfgAddHandler)
		//web.POST("/sales/material/no-sales/config/add", sales.SalesMatlNoSalesCfgAddHandler)
		web.POST("/plant/material/monitor/config/update", sales.SalesMatlNoSalesCfgUpdateHandler)
		//web.POST("/sales/material/no-sales/config/update", sales.SalesMatlNoSalesCfgUpdateHandler)
		web.POST("/plant/material/monitor/config/del", sales.SalesMatlNoSalesCfgDelHandler)
		//web.POST("/sales/material/no-sales/config/del", sales.SalesMatlNoSalesCfgDelHandler)
		web.GET("/plant/material/monitor/config/list", sales.SalesMatlNoSalesCfgListHandler)
		//web.GET("/sales/material/no-sales/config/list", sales.SalesMatlNoSalesCfgListHandler)

		/* sales material plant no-sales API */
		web.POST("/plant/material/no-sales", sales.SalesMatlNoSalesHandler)
		//web.POST("/sales/material/no-sales", sales.SalesMatlNoSalesHandler)

		/* sales material config api */
		web.POST("/sales/material/config/add", sales.SalesMaterialCfgAddHandler)
		web.POST("/sales/material/config/update", sales.SalesMaterialCfgUpdateHandler)
		web.GET("/sales/material/config/del/:code", sales.SalesMaterialCfgDelHandler)
		web.GET("/sales/material/config/clear", sales.SalesMaterialCfgClearHandler)
		web.GET("/sales/material/config/list", sales.SalesMaterialCfgListHandler)

		/* sales material api */
		web.POST("/sales/material/kpi", sales.SalesMaterialKPIHandler)
		web.POST("/sales/material/plant/rank", sales.SalesMaterialPlantRankHandler)
		web.POST("/sales/material/date", sales.SalesMaterialDateHandler)

		/* kpi */
		web.GET("/kpi/list", sales.KpiListHandler)
		web.GET("/kpi/bymonth", sales.KpiByMonth)
		web.POST("/kpi/bymonth", sales.KpiByMonthV2)
		web.POST("/kpi/update", sales.KpiUpdate)
		web.POST("/kpi/delete", sales.KpiDelete)
		web.GET("/kpi/fystart", sales.GetFinanceYearStart)
		web.POST("/kpi/fystart", sales.SetFinanceYearStart)
		web.GET("/kpi/fmonths", sales.GetFinanceMonthList)
		web.POST("kpi/fmonth", sales.SetFinanceMonth)
		web.POST("/kpi/payments", sales.GetPaymentTypes)

		///////////////////////////////////////////////////////////////
		// SYSTEM API
		/* user API */
		web.POST("/user/add", sys.UserAddHandler)
		web.POST("/user/update", sys.UserUpdateHandler)
		web.POST("/user/role-update", sys.UserRoleUpdateHandler)
		web.GET("/user/del/:id", sys.UserDelHandler)
		web.GET("/user/list", sys.UserListHandler)
		web.GET("/user/detail", sys.UserDetailHandler)
		web.GET("/user/view-perm/:username", sys.UserViewPermHandler)
		web.GET("/user/org-perm/:username", sys.UserOrgPermHandler)
		web.POST("/user/view-perm", sys.UpdateUserViewPermHandler)
		web.POST("/user/org-perm", sys.UpdateUserOrgPermHandler)
		web.POST("/user/reset", sys.UserResetPasswordHandler)
	}

	wto := viper.GetInt(common.ConfKeyServerWTO)
	rto := viper.GetInt(common.ConfKeyServerRTO)
	server := &http.Server{
		Addr:         address,
		Handler:      r,
		WriteTimeout: time.Duration(wto) * time.Second,
		ReadTimeout:  time.Duration(rto) * time.Second,
	}

	//non-ssl
	server2 := &http.Server{
		Addr:           address2,
		Handler:        r,
		WriteTimeout:   time.Duration(wto) * time.Second,
		ReadTimeout:    time.Duration(rto) * time.Second,
		MaxHeaderBytes: 4 << 20,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-sigChan
		zaps.Warn("system signal received, ready to quit: ", s)
		handlers.GinShutdown()
		ctx, _ := context.WithTimeout(context.Background(),
			10*time.Second)
		server.Shutdown(ctx)
		server2.Shutdown(ctx)
	}()

	// write pid file
	err = ioutil.WriteFile(viper.GetString(common.ConfKeyPidFile),
		[]byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		zaps.Fatal("failed to write pid file, error: ", err)
		os.Exit(-3)
	}

	if ssl {
		go func() {
			if len(cert) <= 0 || len(key) <= 0 {
				zaps.Fatal("invalid certificate/private key file")
				return
			}

			zaps.Infof("dserver ssl starts running @ %s", address)
			zaps.Fatal(server.ListenAndServeTLS(cert, key))
		}()
		go func() {
			zaps.Infof("dserver starts running @ %s", address2)
			zaps.Fatal(server2.ListenAndServe())
		}()

		select {}
	} else {

		zaps.Infof("dserver starts running @ %s", address)
		zaps.Fatal(server.ListenAndServe())
	}
}
