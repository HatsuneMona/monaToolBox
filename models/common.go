package models

import (
	"gorm.io/gorm"
	"time"
)

type ID struct {
	Id uint `json:"id" gorm:"primaryKey"`
}

type CommonTime struct {
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

type DeleteTime struct {
	DeleteTime gorm.DeletedAt `json:"delete_time" gorm:"index"`
}
