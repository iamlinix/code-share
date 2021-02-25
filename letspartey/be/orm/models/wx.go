package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"iamlinix.com/partay/db"
	"iamlinix.com/partay/fs"
	"iamlinix.com/partay/logger"
)

type WxOpen struct {
	BaseModel
	OpenID      string     `gorm:"column:open_id;unique;unique_index" json:"openId"`
	UnionID     string     `gorm:"column:union_id;type:varchar(64)" json:"unionId"`
	SessionKey  string     `gorm:"column:session_key;type:varchar(64)" json:"sessionKey"`
	Password    string     `gorm:"column:password;not null" json:"password"`
	Avatar      string     `gorm:"column:avatar;type:varchar(255)" json:"avatarUrl"`
	LocalAvatar string     `gorm:"column:local_avatar;type:varchar(255)" json:"localAvatar"`
	NickName    string     `gorm:"column:nickname;type:varchar(64)" json:"nickName"`
	AuthorizeAt *time.Time `gorm:"column:authorize_at" json:"authorize_at"`
}

type WxUser struct {
	BaseModel
	OpenID      string `gorm:"column:open_id;index" json:"openId"`
	NickName    string `gorm:"column:nickname;type:varchar(64)" json:"nickName"`
	Avatar      string `gorm:"column:avatar;type:varchar(255)" json:"avatarUrl"`
	LocalAvatar string `gorm:"column:local_avatar;type:varchar(255)" json:"localAvatar"`
	Country     string `gorm:"column:country;type:varchar(64)" json:"contry"`
	Province    string `gorm:"column:province;type:varchar(64)" json:"province"`
	City        string `gorm:"column:city;type:varchar(64)" json:"city"`
	Gender      int    `gorm:"column:gender" json:"gender"`
	PhoneNumber string `gorm:"column:phone_number" json:"phone_number"`
	EMail       string `gorm:"column:email" json:"email"`
	IP          string `gorm:"column:ip;type:varchar(15)"`
}

type SafeArea struct {
	Left   int `gorm:"column:sa_left" json:"left"`
	Right  int `gorm:"column:sa_right" json:"right"`
	Top    int `gorm:"column:sa_top" json:"top"`
	Bottom int `gorm:"column:sa_bottom" json:"bottom"`
	Width  int `gorm:"column:sa_width" json:"width"`
	Height int `gorm:"column:sa_height" json:"height"`
}

type WxDevice struct {
	BaseModel
	OpenID                      string   `gorm:"column:open_id;index" json:"openId"`
	SDKVer                      string   `gorm:"column:sdk_version;type:varchar(16)" json:"SDKVersion"`
	BatteryLevel                int      `gorm:"column:battery_level" json:"batteryLevel"`
	BenchMark                   int      `gorm:"column:benchmark_level" json:"benchmarkLevel"`
	Brand                       string   `gorm:"column:brand;type:varchar(16)" json:"brand"`
	Orientation                 string   `gorm:"column:device_orientation;type:varchar(16)" json:"deviceOrientation"`
	ErrMsg                      string   `gorm:"column:err_msg:type:text" json:"errMsg"`
	FontSize                    int      `gorm:"column:font_size" json:"fontSizeSetting"`
	Language                    string   `gorm:"column:language;type:varchar(8)" json:"language"`
	Model                       string   `gorm:"column:model" json:"model"`
	DevicePixelRatio            float32  `gorm:"column:device_ratio" json:"devicePixelRatio"`
	PixelRatio                  float32  `gorm:"column:pixel_ratio" json:"pixelRatio"`
	Platform                    string   `gorm:"column:platform;type:varchar(32)" json:"platform"`
	ScreenHeight                int      `gorm:"column:screen_height" json:"screenHeight"`
	ScreenWidth                 int      `gorm:"column:screen_width" json:"screenWidth"`
	StatusBarHeight             int      `gorm:"column:status_bar_height" json:"statusBarHeight"`
	System                      string   `gorm:"column:system" json:"system"`
	Version                     string   `gorm:"column:version" json:"version"`
	WindowHeight                int      `gorm:"column:window_height" json:"windowHeight"`
	WindowWidth                 int      `gorm:"column:window_width" json:"windowWidth"`
	AlbumAuthorized             bool     `gorm:"column:album_authed;type:tinyint" json:"albumAuthorized"`
	CameraAuthorized            bool     `gorm:"column:camera_authed;type:tinyint" json:"cameraAuthorized"`
	LocationAuthorized          bool     `gorm:"column:location_authed;type:tinyint" json:"locationAuthorized"`
	MicrophoneAuthorized        bool     `gorm:"column:microphone_authed;type:tinyint" json:"microphoneAuthorized"`
	NotificationAuthorized      bool     `gorm:"column:notify_authed;type:tinyint" json:"notificationAuthorized"`
	NotificationAlertAuthorized bool     `gorm:"column:notify_alert_authed;type:tinyint" json:"notificationAlertAuthorized"`
	NotificationBadgeAuthorized bool     `gorm:"column:notify_badge_authed;type:tinyint" json:"notificationBadgeAuthorized"`
	NotificationSoundAuthorized bool     `gorm:"column:notify_sound_authed;type:tinyint" json:"notificationSoundAuthorized"`
	BluetoothEnabled            bool     `gorm:"column:bluetooth_enable;type:tinyint" json:"bluetoothEnabled"`
	LocationEnabled             bool     `gorm:"column:location_enable;type:tinyint" json:"locationEnabled"`
	WifiEnabled                 bool     `gorm:"column:wifi_enable;type:tinyint" json:"wifiEnabled"`
	Theme                       string   `gorm:"column:theme;type:varchar(16)"`
	SafeArea                    SafeArea `gorm:"embedded" json:"safeArea"`
}

func (*WxOpen) TableName() string {
	return "wx_open"
}

func (*WxDevice) TableName() string {
	return "wx_devices"
}

func (*WxUser) TableName() string {
	return "wx_users"
}

func WxOpened(openId, unionId, sessionKey string) {
	var open WxOpen
	if _, err := db.Get().ExecuteStruct(func() interface{} { return &open },
		"SELECT * FROM wx_open WHERE open_id = ?", openId); err != nil {
		logger.Errorf("error searching openid @ opened: %#v", err)
	}

	now := time.Now()
	if open.OpenID != openId {
		logger.Warnf("new open id logged in @ open: %s", openId)
		open.OpenID = openId
		open.UnionID = unionId
		open.SessionKey = sessionKey
		open.CTime = &now
		open.MTime = &now
		open.Password = fmt.Sprintf("%X", md5.Sum([]byte(openId)))
		db.Get().Create(&open)
	} else {
		logger.Warnf("authorize comes earlier than open: %s", openId)
		open.OpenID = openId
		open.UnionID = unionId
		open.SessionKey = sessionKey
		open.MTime = &now
		db.Get().Update(&open)
	}
}

func WxAuthorize(openId string) error {
	var open WxOpen
	if _, err := db.Get().ExecuteStruct(func() interface{} { return &open },
		"SELECT * FROM wx_open WHERE open_id = ?", openId); err != nil {
		logger.Errorf("error searching openid @ authorize: %#v", err)
		return err
	}

	now := time.Now()
	if len(open.OpenID) == 0 {
		logger.Warnf("new open id logged in @ authorize: %s", openId)
		open.OpenID = openId
		open.CTime = &now
		open.MTime = &now
		open.Password = fmt.Sprintf("%X", md5.Sum([]byte(openId)))
		db.Get().Create(&open)
	} else {
		logger.Warnf("authorize comes earlier than open: %s", openId)
		open.AuthorizeAt = &now
		open.MTime = &now
		db.Get().Update(&open)
	}

	return nil
}

func WxUserReport(user *WxUser) error {
	var open WxOpen
	if _, err := db.Get().ExecuteStruct(func() interface{} { return &open },
		"SELECT * FROM wx_open WHERE open_id = ?", user.OpenID); err != nil {
		logger.Errorf("error searching openid: %#v", err)
		return err
	}

	if len(open.OpenID) == 0 {
		logger.Errorf("no open id found for user report: %#v", user)
		return errors.New("No openid found for user")
	}

	if len(open.LocalAvatar) == 0 || open.Avatar != user.Avatar {
		go WxGetUserAvatar(user.OpenID, user.Avatar)
	}

	if open.NickName != user.NickName {
		go WxUpdateUserNickName(user.OpenID, user.NickName)
	}

	db.Get().Create(user)

	return nil
}

func WxDeviceReport(device *WxDevice) {
	if err := db.Get().Create(device); err != nil {
		logger.Errorf("error creating user device: %v", err)
	}
}

func WxUpdateUserNickName(openId, nickname string) {
	if _, err := db.Get().Execute("UPDATE wx_open SET nickname= ? WHERE open_id = ?", nickname, openId); err != nil {
		logger.Errorf("error updating user nickname: %s,%s", openId, nickname)
	}
}

func WxGetUserAvatar(openId, url string) {
	if resp, err := http.Get(url); err != nil {
		logger.Errorf("error getting user avatar %s: %v", openId, err)
	} else {
		if data, err := ioutil.ReadAll(resp.Body); err != nil {
			logger.Errorf("error reading user avatar %s: %v", openId, err)
		} else {
			localAvartar := fmt.Sprintf("files/images/avartars/%X", md5.Sum([]byte(openId)))
			file, err := fs.Get().Create(localAvartar)
			if err != nil {
				logger.Errorf("error creating avartar file: %v", err)
				return
			}

			if err = fs.Get().Write(file, data); err != nil {
				logger.Errorf("error writing avartar file: %v", err)
				return
			}

			if _, err := db.Get().Execute("UPDATE wx_open SET avatar = ?, local_avatar = ?, mtime = current_timestamp()", url,
				localAvartar); err != nil {
				logger.Errorf("error updating local avartar: %v", err)
			}
		}
	}
}

func WxGetUserBasic(openId string) (string, string, error) {
	var open WxOpen
	if _, err := db.Get().ExecuteStruct(func() interface{} { return &open },
		"SELECT * FROM wx_open WHERE open_id = ?", openId); err != nil {
		logger.Errorf("error searching openid @ authorize: %#v", err)
		return "", "", err
	}

	return open.NickName, open.LocalAvatar, nil
}
