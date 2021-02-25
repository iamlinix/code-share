package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"iamlinix.com/partay/cipher"
	"iamlinix.com/partay/conf"
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/fs"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm"
	"iamlinix.com/partay/web"
	"iamlinix.com/partay/web/handlers"
	"iamlinix.com/partay/web/middleware/auth"
	"iamlinix.com/partay/web/middleware/throttle"
	"iamlinix.com/partay/web/wechat"
)

func main() {
	// arguments parsing
	enc := flag.String("enc", "", "string to encrypt")
	dec := flag.String("dec", "", "string to decrypt")
	key := flag.String("key", "", "encryption key")
	confFile := flag.String("config", "conf.yaml", "config file path")
	flag.Parse()

	if len(*enc) > 0 && len(*key) > 0 {
		encrypted, err := cipher.Encrypt(*enc, *key)
		if err != nil {
			logger.Infof("error encrypting: %v\n", err)
		} else {
			logger.Infof("encrypted string: [%s]\n", encrypted)
		}
		return
	}

	if len(*dec) > 0 && len(*key) > 0 {
		decrypted, err := cipher.Decrypt(*dec, *key)
		if err != nil {
			logger.Infof("error decrypting: %v\n", err)
		} else {
			logger.Infof("encrypted string: [%s]\n", decrypted)
		}
		return
	}

	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		logger.Infof("error reading config file: %v ", err)
		return
	}

	var c conf.Conf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		logger.Infof("error parsing config file: %v", err)
		return
	}

	logger.InitLogger(c.Log.App.Path, c.Log.App.Level)
	logger.Info("logger initated")

	cipher.SetGlobalKey(c.Encrypt.Code)

	logger.Infof("setting fs backend: %s", fs.BackendMap[fs.BackendType(c.FileSystem.Backend)])
	if err = fs.SetFsBackend(fs.BackendType(c.FileSystem.Backend), c.FileSystem.BaseDir); err != nil {
		logger.Errorf("error setting fs backend: %v", err)
		return
	}
	if err = fs.Get().Init(c.FileSystem.Extras); err != nil {
		logger.Errorf("error initializing fs backend: %v", err)
		return
	}

	logger.Infof("setting db backend: %s", db.BackendMap[db.BackendType(c.Database.Backend)])
	if err = db.SetDbBackend(db.BackendType(c.Database.Backend)); err != nil {
		logger.Errorf("error setting db backend: %v", err)
		return
	}

	dbUsername, err := cipher.Decrypt(c.Database.Username, c.Encrypt.Code)
	if err != nil {
		logger.Errorf("error decrypting db username: %v", err)
		return
	}

	dbPassword, err := cipher.Decrypt(c.Database.Password, c.Encrypt.Code)
	if err != nil {
		logger.Errorf("error decrypting db password: %v", err)
		return
	}

	dbDatabase, err := cipher.Decrypt(c.Database.Database, c.Encrypt.Code)
	if err != nil {
		logger.Errorf("error decrypting db database: %v", err)
		return
	}

	logger.Info("connecting database")
	if err = db.Get().Connect(c.Database.Driver, dbUsername, dbPassword, dbDatabase, c.Database.Extras, c.Database.Pool); err != nil {
		logger.Errorf("error while connecting database: %v", err)
		return
	}

	logger.Info("initating orm")
	if err = orm.InitOrmWithDbInst(db.Get()); err != nil {
		logger.Errorf("error initating orm with db instance: %v", err)
		return
	}

	if err = orm.InitAllTables(); err != nil {
		logger.Errorf("failed to init tables: %v", err)
		return
	}

	appId, err := cipher.Decrypt(c.Wx.AppId, c.Encrypt.Code)
	if err != nil {
		logger.Errorf("error decrypting appid: %v", err)
		return
	}

	appSecret, err := cipher.Decrypt(c.Wx.Secret, c.Encrypt.Code)
	if err != nil {
		logger.Errorf("error decrypting secret: %v", err)
		return
	}
	wechat.WxInit(appId, appSecret)

	if c.Web.Auth.Enable {
		auth.SetAuthMethod(auth.AuthMethod(c.Web.Auth.Method))
		auth.Get().Init(c.Web.Auth.Extras)
	}

	gin.DisableConsoleColor()
	accessFd, _ := os.Create(c.Log.Access.Path)
	gin.DefaultWriter = io.MultiWriter(accessFd)

	gin.SetMode(c.Mode)
	r := gin.Default()
	r.Use(gin.Recovery())
	if c.Web.Cors {
		r.Use(web.Cors())
	}

	if c.Web.Throttle.Enable {
		if err = throttle.InitThrottle(c.Web.Throttle.Urls, c.Web.Throttle.MaxPerSec, c.Web.Throttle.MaxBurst); err != nil {
			logger.Errorf("throttle initialization error: %v", err)
			return
		}
		r.Use(throttle.Throttle())
	}

	for _, v := range c.Web.Serves {
		logger.Infof("serving directories: %s > %s", v["k"], v["v"])
		//r.Static(v["k"], v["v"])
		r.Use(static.Serve(v["k"], static.LocalFile(v["v"], true)))
	}

	for _, v := range c.Web.Statics {
		logger.Infof("serving static urls: %s > %s", v["k"], v["v"])
		r.StaticFile(v["k"], v["v"])
	}

	api := r.Group("/api")
	{
		wx := api.Group("/wx")
		{
			wx.POST("/open", handlers.HdlrWxCode2Session)
			wx.POST("/user", handlers.HdlrWxUserInfo)
			wx.POST("/device", handlers.HdlrWxSysInfo)
			wx.POST("/basic", handlers.HdlrWxGetUserInfo)
		}

		api.GET("/ping", handlers.Ping)
		api.POST("/login", handlers.HdlrLogin)
		api.POST("/signup", handlers.HdlrSignUp)
		api.POST("/download", handlers.HdlrPostReadFile)
		api.GET("/files/*any", handlers.HdlrReadFile)
		api.GET("/whoami", handlers.HdlrWhoAmI)

		v1 := api.Group("/v1")
		if c.Web.Auth.Enable {
			v1.Use(auth.Get().AuthMiddleware())
		}
		{
			v1.GET("/ping", handlers.Ping)
			v1.GET("/activities", handlers.HdlrListActivity)
			v1.POST("/activity", handlers.HdlrCreateActivity)
			v1.GET("/posts", handlers.HdlrListPost)
			v1.POST("/post", handlers.HdlrCreatePost)
			v1.POST("/upload", handlers.HdlrUploadImage)
			v1.GET("/user", handlers.HdlrGetUser)
			v1.POST("/password", handlers.HdlrUpdatePassword)

		}
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Web.Bind, c.Web.Port),
		Handler:      r,
		WriteTimeout: time.Duration(c.Web.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.Web.ReadTimeout) * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-sigChan
		logger.Warn("system signal received, ready to quit: ", s)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	if c.Web.SSL.Enable {
		logger.Infof("running ssl on port %d", c.Web.Port)
		logger.Warn(server.ListenAndServeTLS(c.Web.SSL.CertFile, c.Web.SSL.KeyFile))
	} else {
		logger.Infof("running on port %d", c.Web.Port)
		logger.Warn(server.ListenAndServe())
	}
}
