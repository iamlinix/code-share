package models

type Post struct {
	BaseModel
	Name    string `gorm:"column:name;type:varchar(32)" json:"name"`
	Creator string `gorm:"column:creator;type:varchar(64)" json:"creator"`
	Desc    string `gorm:"column:desc;type:text" json:"desc"`
}

type PostImage struct {
	PostID int    `gorm:"column:post_id;type:int" json:"postID"`
	URL    string `gorm:"column:url;type:varchar(255)" json:"url"`
}
