package dto

type UserLoginDTO struct {
	Username string `json:"username" validate:"required" label:"username"`
	Password string `json:"password" validate:"required" label:"password"`
}

type QueryUserDTO struct {
	GroupId int `json:"group_id" form:"group_id" validate:"omitempty" label:"分组ID"`
	Page    int `json:"page" form:"page" validate:"omitempty,min=0" label:"分页数"`
	Count   int `json:"count" form:"count" validate:"required,gte=1" label:"每页数量"`
}
