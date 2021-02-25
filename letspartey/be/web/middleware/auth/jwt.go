package auth

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/cipher"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/web"
)

var (
	ErrTokenExpired      error = errors.New("TOKEN EXPIRED")
	ErrTokenNotValidYet  error = errors.New("TOKEN INACTIVE")
	ErrTokenMalformed    error = errors.New("MALFORMED TOKEN")
	ErrTokenInvalid      error = errors.New("INVALID TOKEN")
	ErrTokenInsufficient error = errors.New("INSUFFICIENT TOKEN")

	_signKey         []byte        = nil
	_issuer          string        = ""
	_expireInMinutes time.Duration = 0
)

type JwtAuthClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type JwtAuth struct {
}

func (*JwtAuth) Init(args map[string]interface{}) error {
	if val, ok := args["sign_key"]; ok {
		if decrypted, err := cipher.DefaultDecrypt(val.(string)); err != nil {
			logger.Error("failed to decrypt jwt signkey: %v", err)
			return errors.New("Signkey decryption fail")
		} else {
			_signKey = []byte(decrypted)
		}
	} else {
		logger.Errorf("no jwt signkey config")
		return errors.New("No jwt signkey")
	}

	if val, ok := args["issuer"]; ok {
		var err error
		if _issuer, err = cipher.DefaultDecrypt(val.(string)); err != nil {
			logger.Error("failed to decrypt jwt issuer: %v", err)
			return errors.New("Issuer decryption fail")
		}
	} else {
		logger.Errorf("no jwt issuer config")
		return errors.New("No jwt issuer")
	}

	if val, ok := args["expire"]; ok {
		_expireInMinutes = time.Duration(val.(int))
	} else {
		logger.Errorf("no jwt expire config")
		return errors.New("No jwt expire")
	}

	return nil
}

func (*JwtAuth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusForbidden, &web.BaseResponse{
				Code: web.ECTokenNoToken,
			})
			c.Abort()
			logger.Error("No token found")
			return
		}

		claims, err := VerifyToken(token, c.Request.URL)
		if err != nil {
			switch err {
			case ErrTokenExpired:
				logger.Error("token has expired")
				c.JSON(http.StatusForbidden, &web.BaseResponse{
					Code:    web.ECTokenExpire,
					Message: web.EMTokenExpire,
				})
				c.Abort()
				return

			case ErrTokenInvalid:
				logger.Error("token is invalid")
				c.JSON(http.StatusForbidden, &web.BaseResponse{
					Code:    web.ECTokenInvalid,
					Message: web.EMTokenInvalid,
				})
				c.Abort()
				return

			case ErrTokenMalformed:
				logger.Error("token is malformed")
				c.JSON(http.StatusForbidden, &web.BaseResponse{
					Code:    web.ECTokenMalform,
					Message: web.EMTokenMalform,
				})
				c.Abort()
				return

			case ErrTokenNotValidYet:
				logger.Error("token is inactive")
				c.JSON(http.StatusForbidden, &web.BaseResponse{
					Code:    web.ECTokenInactive,
					Message: web.EMTokenInactive,
				})
				c.Abort()
				return

			default:
				break
			}
		}

		c.Set("claims", claims)
		logger.Info("token verify success")
	}
}

func (*JwtAuth) GenerateToken(username string) (string, error) {
	claims := JwtAuthClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Minute * _expireInMinutes).Unix()),
			Issuer:    _issuer,
		},
		username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(_signKey)

	if err != nil {
		logger.Errorf("failed to generate token for user %s: %v", username, err)
		return "", err
	}

	logger.Infof("generated token for user %s: %s", username, token)

	return signed, nil
}

func VerifyToken(tokenString string, url *url.URL) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtAuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return _signKey, nil
		},
	)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*JwtAuthClaims); ok && token.Valid {
		logger.Infof("username: %s, expire: %d", claims.Username, claims.StandardClaims.ExpiresAt)
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
