package handlers

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/json"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
	"iamlinix.com/partay/web"
	"iamlinix.com/partay/web/middleware/auth"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type LoginResponse struct {
	web.BaseResponse
	User
	Token string `json:"token"`
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func HdlrLogin(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("failed to read login data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenCorruptBody,
			Message: web.EMGenCorruptBody,
		})
		return
	}

	var user User
	if err = json.JsonUnmarshal(data, &user); err != nil {
		logger.Errorf("error parsing login data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return
	}

	usr, err := models.CheckUserPassword(user.Username, user.Password)
	if err != nil {
		logger.Errorf("error checking user password: %v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	if usr == nil {
		logger.Warnf("user login wrong password: %s", user.Username)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECUserWrongPassword,
			Message: web.EMUserWrongPassword,
		})
		return
	} else {
		var token string
		if auth.WebAuthEnabled {
			token, err = auth.Get().GenerateToken(user.Username)
		}

		if err != nil {
			logger.Errorf("error generating token for user %s: %v", user.Username, err)
			c.JSON(http.StatusInternalServerError, &web.BaseResponse{
				Code:    web.ECGenUnknown,
				Message: web.EMGenUnknown,
			})
			return
		}

		user.Avatar = usr.Avartar
		c.JSON(http.StatusOK, &LoginResponse{
			web.BaseResponse{
				Code:    web.ECOK,
				Message: web.EMOK,
			},
			user,
			token,
		})
	}
}

func HdlrSignUp(c *gin.Context) {
	var user User
	if err := UnmarshalRequestData(c, &user); err != nil {
		return
	}

	if _, err := db.Get().Execute("INSERT INTO users (username, password) VALUES "+
		"(?, ?)", user.Username, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	if res, err := db.Get().Execute("SELECT avatar FROM users WHERE username = ?", user.Username); err != nil {
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	} else {
		user.Avatar = string(*res[0]["avatar"].(*sql.RawBytes))
		c.JSON(http.StatusOK, &LoginResponse{
			web.BaseResponse{
				Code:    web.ECOK,
				Message: web.EMOK,
			},
			user,
			"",
		})
	}
}

func HdlrGetUser(c *gin.Context) {
	username := c.Query("username")
	var user User
	if res, err := db.Get().Execute("SELECT avatar FROM users WHERE username = ?", username); err != nil {
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	} else {
		user.Avatar = string(*res[0]["avatar"].(*sql.RawBytes))
		c.JSON(http.StatusOK, &LoginResponse{
			web.BaseResponse{
				Code:    web.ECOK,
				Message: web.EMOK,
			},
			user,
			"",
		})
	}
}

func HdlrUpdatePassword(c *gin.Context) {
	var user User
	if err := UnmarshalRequestData(c, &user); err != nil {
		return
	}

	if _, err := db.Get().Execute("UPDATE users SET password = ? WHERE username = ?", user.Password, user.Username); err != nil {
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	c.JSON(http.StatusOK, &web.BaseResponse{
		Code:    web.ECOK,
		Message: web.EMOK,
	})
}
