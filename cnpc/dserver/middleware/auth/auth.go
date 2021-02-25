package auth

import (
	"cnpc.com.cn/cnpc/dserver/common"
	"cnpc.com.cn/cnpc/dserver/zaps"

	"errors"
	"net/http"
	"net/url"
	"time"

	"cnpc.com.cn/cnpc/dserver/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenExpired      error = errors.New("Token is expired")
	TokenNotValidYet  error = errors.New("Token not active yet")
	TokenMalformed    error = errors.New("That's not even a token")
	TokenInvalid      error = errors.New("Couldn't handle this token:")
	TokenInsufficient error = errors.New("Token insufficient :")
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int    `json:"role"`
}

func VerifyToken(c *gin.Context, tokenString string,
	url *url.URL) (jwt.Claims, error) {

	laddr := c.Request.Header.Get("Local-Addr")
	paddr := c.Request.Header.Get("Public-Addr")

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("maygodblessyou"), nil
		})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				zaps.Error("<<< token validation error")
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		zaps.Warnf("lip: %s, pip: %s, username: %s, access url: %s",
			//c.ClientIP(), claims.Username, url.Path)
			laddr, paddr, claims.Username, url.Path)

		if claims.Role != 1 {
			if url.Path == "/v1/web/user/reset" {
				zaps.Warn("<<< API ONLY FOR Admin")
				return nil, TokenInsufficient
			}
		}

		return claims, nil
	}

	zaps.Error("<<< verify token failed")

	return nil, TokenInvalid
}

//XXX TODO: regresh token
/*
func RefreshToken(tokenString string) (string, error) {

	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{},
			func(token *jwt.Token) (interface{}, error) {
		return []byte("maygodblessyou"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenInvalid
}
*/

func JWTAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		laddr := c.Request.Header.Get("Local-Addr")
		paddr := c.Request.Header.Get("Public-Addr")
		url := c.Request.URL
		db := orm.GetDB()
		user := "none"
		ip := c.ClientIP()

		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": -1,
				"msg":    "no auth token",
			})
			c.Abort()
			zaps.Warnf("<<< no token found: (lip: %s, "+
				"pip: %s, url: %s)",
				laddr, paddr, url.Path)
			if rows, err := db.Query("INSERT INTO access_log (url, username, ip, ts) VALUES (?, ?, ?, CURRENT_TIMESTAMP)", url.String(), user, ip); err != nil {
				zaps.Error("failed to log access log: ", err)
			} else {
				rows.Close()
			}
			return
		}

		//zaps.Info("got token: ", token)
		user = "expire"

		claims, err := VerifyToken(c, token, c.Request.URL)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusForbidden, gin.H{
					"status": -1,
					"msg":    "token expired",
				})
				c.Abort()
				zaps.Warnf("<<< token expired: (lip: %s, "+
					"pip: %s, url: %s)",
					laddr, paddr, url.Path)
				if rows, err := db.Query("INSERT INTO access_log (url, username, ip, ts) VALUES (?, ?, ?, CURRENT_TIMESTAMP)", url.String(), user, ip); err != nil {
					zaps.Error("failed to log access log: ", err)
				} else {
					rows.Close()
				}
				return
			}

			c.JSON(http.StatusForbidden, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})

			c.Abort()
			zaps.Warnf("<<< verify token failed: (lip: %s, "+
				"pip: %s, url: %s)",
				laddr, paddr, url.Path)
			return
		}

		user = claims.(*AuthClaims).Username
		if rows, err := db.Query("INSERT INTO access_log (url, username, ip, ts) VALUES (?, ?, ?, CURRENT_TIMESTAMP)", url.String(), user, ip); err != nil {
			zaps.Error("failed to log access log: ", err)
		} else {
			rows.Close()
		}
		c.Set("claims", claims)

		zaps.Debug("<<< verify token success")
	}
}

///////////////////////////////////////////////////////////////////////////////
//

func GenerateToken(user common.User) (string, error) {

	zaps.Info(">>> generate token start")

	authSigningKey := []byte("maygodblessyou")

	claims := AuthClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 24 * 365).Unix()),
			Issuer:    "cnpc.com.cn",
		},
		user.Username,
		user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(authSigningKey)

	if err != nil {
		zaps.Error("generate token failed: ", err)
		return "", err
	}

	zaps.Info("<<< generate token: ", token)

	return ss, nil
}
