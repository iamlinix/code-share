package wechat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"iamlinix.com/partay/json"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/orm/models"
	"iamlinix.com/partay/web"
)

type WxLoginResp struct {
	web.BaseResponse
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
	SessionKey string `json:"session_key"`
	Token      string `json:"token"`
}

const (
	WxCode2SessionUrl string = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	WxLoginErrorCodeSuccess     int = 0
	WxLoginErrorCodeBusy        int = -1
	WxLoginErrorCodeInvalidCode int = 40029
	WxLoginErrorCodeFrequent    int = 45011
)

var _appId string = ""
var _secret string = ""

func WxInit(appId, secret string) {
	_appId = appId
	_secret = secret
}

func WxApply4OpenId(code string) (*WxLoginResp, error) {
	var loginResp WxLoginResp
	resp, err := http.Get(fmt.Sprintf(WxCode2SessionUrl, _appId, _secret, code))
	if err != nil {
		logger.Errorf("wx auth api request error: %v", err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("wx auth api body read error: %v", err)
		return nil, err
	}

	if err = json.JsonUnmarshal(data, &loginResp); err != nil {
		logger.Errorf("error unmarshaling wx login resp data: %v", err)
		return nil, err
	}

	logger.Infof("wx login response: %v", loginResp)
	if loginResp.Code != WxLoginErrorCodeSuccess {
		logger.Errorf("wx login error: %d:%s", loginResp.Code, loginResp.Message)
		return nil, errors.New(loginResp.Message)
	}

	go models.WxOpened(loginResp.OpenId, loginResp.UnionId, loginResp.SessionKey)

	return &loginResp, nil
}
