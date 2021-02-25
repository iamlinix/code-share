package rdt

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"errors"
	"time"
)


func RedirectToTLS() gin.HandlerFunc {

	return func(c *gin.Context) {

		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:"10.191.7.105:8081",
		})

		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}

		c.Next()
	}
}


