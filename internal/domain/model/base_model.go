package model

import "time"

type BaseModel struct {
	ID         int        `json:"id" db:"primaryKey"`
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time  `json:"update_time"`
	DeleteTime *time.Time `json:"delete_time"`
}
