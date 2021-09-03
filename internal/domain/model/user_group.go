package model

type UserGroup struct {
	ID      int `json:"id" db:"primaryKey"`
	UserID  int `json:"user_id"`
	GroupID int `json:"group_id"`
}
