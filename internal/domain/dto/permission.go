package dto

type DispatchPermissionDTO struct {
	GroupId      int `json:"group_id" validate:"required,min=1" label:"group_id"`
	PermissionId int `json:"permission_id" validate:"required,min=1" label:"permission_id"`
}

type DispatchPermissionsDTO struct {
	GroupId       int   `json:"group_id" validate:"required,min=1" label:"group_id"`
	PermissionIds []int `json:"permission_ids" validate:"required" label:"permission_ids"`
}
