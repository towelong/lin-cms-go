package model

// User 用户表
type User struct {
	BaseModel
	Username string `json:"username" validate:"required" label:"用户名"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}
