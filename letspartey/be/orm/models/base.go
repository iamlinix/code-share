package models

import "time"

type BaseModel struct {
	ID    int        `gorm:"column:id;auto_increment;primary_key" json:"id"`
	CTime *time.Time `gorm:"column:ctime;type:datetime;default:CURRENT_TIMESTAMP" json:"ctime"`
	MTime *time.Time `gorm:"column:mtime;type:datetime;default:CURRENT_TIMESTAMP" json:"mtime"`
	DTime *time.Time `gorm:"column:dtime;type:datetime" json:"dtime"`
}
