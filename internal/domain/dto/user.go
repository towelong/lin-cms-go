package dto

type BasePage struct {
	Page  int `json:"page" form:"page" validate:"omitempty,min=0" label:"分页数"`
	Count int `json:"count" form:"count" validate:"required,gte=1" label:"每页数量"`
}

type UserLoginDTO struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

type QueryUserDTO struct {
	GroupId int `json:"group_id" form:"group_id" validate:"omitempty" label:"分组ID"`
	BasePage
}

type UsrID struct {
	ID int `uri:"id" json:"id"  validate:"required,gt=0" label:"用户编号"`
}

type ResetPasswordDTO struct {
	NewPassword     string `json:"new_password" validate:"required" label:"新密码"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" label:"确认密码"`
}

type RegisterDTO struct {
	Username        string `json:"username" validate:"required,gte=2,lte=10" label:"用户名"`
	Email           string `json:"email" validate:"omitempty,email" label:"邮箱"`
	GroupIds        []int  `json:"group_ids"`
	Password        string `json:"password" validate:"required" label:"密码"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" label:"确认密码"`
}

type ChangePasswordDTO struct {
	NewPassword     string `json:"new_password" validate:"required" label:"新密码"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" label:"确认密码"`
	OldPassword     string `json:"old_password" validate:"required"`
}

type UpdateInfoDTO struct {
	Email    string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Nickname string `json:"nickname" validate:"omitempty,gte=2,lte=10" label:"昵称"`
	Username string `json:"username" validate:"omitempty,gte=2,lte=10" label:"用户名"`
	Avatar   string `json:"avatar" validate:"omitempty,lte=500" label:"头像"`
}

type UpdateGroupsDTO struct {
	GroupIds []int `json:"group_ids" validate:"required,dive,min=1" label:"GroupIds"`
}
