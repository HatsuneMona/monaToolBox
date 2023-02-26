package types

import (
	"monaToolBox/global/validator"
	"time"
)

type TinyUrlInfo struct {
	Id              uint      `json:"id"`
	TinyUrl         string    `json:"tiny_url"`
	OriginalUrl     string    `json:"original_url"`
	Pv              uint      `json:"pv"`
	CreateTime      time.Time `json:"create_time"`
	LimitAccessTime time.Time `json:"limit_access_time"`
}

type AddTinyUrlForm struct {
	TinyUrl         string    `json:"tiny_url" form:"tiny_url" binding:"required"`
	OriginalUrl     string    `json:"original_url" form:"original_url" binding:"required"`
	LimitAccessTime time.Time `json:"limit_access_time" form:"limit_access_time" binding:""`
}

func (a AddTinyUrlForm) GetMessages() validator.ValidatorMassages {
	return validator.ValidatorMassages{
		"tiny_url.required":     "短链不能为空",
		"original_url.required": "目标Url不能为空",
	}
}
