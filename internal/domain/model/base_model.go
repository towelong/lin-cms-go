package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID         int        `json:"id" db:"primaryKey"`
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time  `json:"update_time"`
	DeleteTime gorm.DeletedAt `json:"delete_time"`
}
