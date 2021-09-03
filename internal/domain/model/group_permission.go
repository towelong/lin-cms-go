package model

type GroupPermission struct {
	ID           int `json:"id" db:"primaryKey"`
	GroupID      int `json:"group_id"`
	PermissionID int `json:"permission_id"`
}
