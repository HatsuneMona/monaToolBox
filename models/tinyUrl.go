package models

import "time"

type TinyUrl struct {
	ID
	TinyUrl         string    `json:"tiny_url" gorm:"type:varchar(32);uniqueIndex:idx_tinyUrl;not null;comment:短链路由"`
	OriginalUrl     string    `json:"original_url" gorm:"type:varchar(32);not null;comment:原始网址"`
	Pv              uint      `json:"pv" gorm:"not null;default:0;comment:点击量"`
	LimitAccessTime time.Time `json:"limit_access_time" gorm:"not null;comment:过了这个时间点后不允许访问"`
	CommonTime
	DeleteTime
}
