package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

const (
	ECOK int    = 0
	EMOK string = "OK"

	ECTokenNoToken  int    = 10001
	EMTokenNoToken  string = "token not found"
	ECTokenExpire   int    = 10002
	EMTokenExpire   string = "token expired"
	ECTokenInactive int    = 10003
	EMTokenInactive string = "token is inactive"
	ECTokenInvalid  int    = 10004
	EMTokenInvalid  string = "invalid token"
	ECTokenMalform  int    = 10005
	EMTokenMalform  string = "malformed token"

	ECGenUnknown       int    = 20001
	EMGenUnknown       string = "unknkown error"
	ECGenCorruptBody   int    = 20002
	EMGenCorruptBody   string = "corrupted request data"
	ECGenIncorrectBody int    = 20003
	EMGenIncorrectBody string = "incorrect request data"

	ECUserWrongPassword int    = 30001
	EMUserWrongPassword string = "wrong password"
	ECUserNotExist      int    = 30002
	EMUserNotExist      string = "user not exist"

	ECResourceNotFound int    = 40001
	EMResourceNotFound string = "resource not found"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers",
			"Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, Content-Disposition")
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
