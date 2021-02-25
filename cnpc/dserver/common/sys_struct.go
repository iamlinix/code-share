package common

///////////////////////////////////////////////////////////////////////////////
//

/* user info struct */
type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	Role        int    `json:"role"`
	ViewPerm    string `json:"viewPerm"`
	OrgPerm     string `json:"orgPerm"`
	CreateTime  string `json:"createTime"`
}

/* user add response data */
type UserAddResp struct {
	Code   int   `json:"code"`
	UserID int64 `json:"userID"`
}

/* user list response data */
type UserListResp struct {
	Code  int    `json:"code"`
	Total int    `json:"total"`
	Users []User `json:"users"`
}

/* user one response data */
type UserDetailResp struct {
	Code int  `json:"code"`
	User User `json:"user"`
}

/* user login response data */
type AuthLoginResp struct {
	Code     int    `json:"code"`
	UserID   int64  `json:"userID"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Role     int    `json:"role"`
	ViewPerm string `json:"viewPerm"`
	OrgPerm  string `json:"orgPerm"`
	Mesg     string `json:"message"`
}
