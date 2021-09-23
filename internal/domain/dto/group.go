package dto

type NewGroupDTO struct {
	Name          string `json:"name" validate:"required,gte=1,lte=60" label:"分组名"`
	Info          string `json:"info" validate:"omitempty,lte=255" label:"描述"`
	PermissionIds []int  `json:"permission_ids"`
}

type UpdateGroupDTO struct {
	Name string `json:"name" validate:"required,gte=1,lte=60" label:"分组名"`
	Info string `json:"info" validate:"omitempty,lte=255" label:"描述"`
}
