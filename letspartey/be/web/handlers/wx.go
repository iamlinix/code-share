package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
	"iamlinix.com/partay/web"
	"iamlinix.com/partay/web/middleware/auth"
	"iamlinix.com/partay/web/wechat"
)

type WxCode2SessionRequest struct {
	Code string `json:"code"`
}

type WxUserInfoRequest struct {
	OpenID string `json:"openId"`
}

type WxUserInfoResponse struct {
	web.BaseResponse
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

func HdlrWxCode2Session(c *gin.Context) {
	var session WxCode2SessionRequest
	if err := UnmarshalRequestData(c, &session); err != nil {
		logger.Errorf("wx openid body read error: %#v", err)
		return
	}

	if wxResp, err := wechat.WxApply4OpenId(session.Code); err != nil {
		c.JSON(http.StatusNotFound, &wechat.WxLoginResp{
			BaseResponse: web.BaseResponse{
				Code:    web.ECResourceNotFound,
				Message: web.EMResourceNotFound,
			},
		})
		return
	} else {
		token, err := auth.Get().GenerateToken(wxResp.OpenId)
		if err != nil {
			logger.Errorf("error generating token for openid: %v", err)
			c.JSON(http.StatusInternalServerError, &web.BaseResponse{
				Code:    web.ECGenUnknown,
				Message: web.EMGenUnknown,
			})
			return
		}

		wxResp.Token = token
		c.JSON(http.StatusOK, wxResp)
	}
}

func HdlrWxUserInfo(c *gin.Context) {
	var user models.WxUser
	if err := UnmarshalRequestData(c, &user); err != nil {
		logger.Errorf("wx login body read error: %#v", err)
		return
	}

	user.IP = c.ClientIP()
	if err := models.WxUserReport(&user); err != nil {
		logger.Errorf("error running model wx user login: %v", err)
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

func HdlrWxSysInfo(c *gin.Context) {
	var device models.WxDevice
	if err := UnmarshalRequestData(c, &device); err != nil {
		logger.Errorf("wx login body read error: %#v", err)
		return
	}

	go models.WxDeviceReport(&device)
	c.JSON(http.StatusOK, &web.BaseResponse{
		Code:    web.ECOK,
		Message: web.EMOK,
	})
}

func HdlrWxGetUserInfo(c *gin.Context) {
	var request WxUserInfoRequest
	if err := UnmarshalRequestData(c, &request); err != nil {
		logger.Errorf("wx login body read error: %#v", err)
		return
	}

	nickname, avatar, err := models.WxGetUserBasic(request.OpenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	if len(nickname) == 0 || len(avatar) == 0 {
		c.JSON(http.StatusNotFound, &web.BaseResponse{
			Code:    web.ECResourceNotFound,
			Message: web.EMResourceNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, &WxUserInfoResponse{
		BaseResponse: web.BaseResponse{
			Code:    web.ECOK,
			Message: web.EMOK,
		},
		NickName: nickname,
		Avatar:   avatar,
	})
}
