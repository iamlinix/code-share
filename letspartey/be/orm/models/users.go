package models

import (
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/logger"
)

type User struct {
	BaseModel
	Name     string `gorm:"column:username;type:varchar(64);unique;not null"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	Avartar  string `gorm:"column:avatar;type:varchar(255);"`
}

func (*User) TableName() string {
	return "users"
}

func CheckUserPassword(username, password string) (*User, error) {
	var user User
	_, err := db.Get().ExecuteStruct(func() interface{} { return &user },
		"SELECT username, password, avatar FROM users WHERE username = ? AND password = ?",
		username, password)
	if err != nil {
		logger.Errorf("error checking user passord: %v", err)
		return nil, err
	}

	if user.Name == username && user.Password == password {
		return &user, nil
	}

	logger.Errorf("user login mismatch: %s,%s, %s,%s", username, user.Name, password, user.Password)
	return nil, nil
}
