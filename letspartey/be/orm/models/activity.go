package models

import "time"

type Activity struct {
	BaseModel
	Creator   string     `gorm:"column:creator;type:varchar(64)" json:"creator"`
	Name      string     `gorm:"column:name;type:varchar(32)" json:"name`
	StartTime *time.Time `gorm:"column:start_time;type:datetime;" json:"startTime"`
	Desc      string     `gorm:"column:desc;type:text" json:"desc"`
}

type ActivityImage struct {
	ActivityID int    `gorm:"column:activity_id;type:int" json:"activityID"`
	URL        string `gorm:"column:url;type:varchar(255)" json:"url"`
}
