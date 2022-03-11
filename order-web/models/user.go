package models

type UserInfo struct {
	Id       int32  `json:"id"`
	Birthday uint64 `json:"birthday"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Password string `json:"password"`
	Role     int32  `json:"role"`
}

type UserInfoList struct {
	Total    int32      `json:"total"`
	UserInfo []UserInfo `json:"user_info"`
}

type Page struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
