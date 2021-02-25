package auth

import (
	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/logger"
)

type AuthMethod int

const (
	AuthNone AuthMethod = 0
	AuthJwt  AuthMethod = 1
)

type AuthInterface interface {
	Init(map[string]interface{}) error
	AuthMiddleware() gin.HandlerFunc
	GenerateToken(string) (string, error)
}

type BaseAuth struct {
}

func (*BaseAuth) Init(map[string]interface{}) error {
	logger.Warn("base auth interface does absolutely nothing")
	return nil
}

func (*BaseAuth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (*BaseAuth) GenerateToken(string) (string, error) {
	logger.Warn("base auth interface generates nil token")
	return "", nil
}
