package models

import (
	"strconv"
)

type User struct {
	ID
	LoginName string `json:"loginName" gorm:"type:varchar(32);not_null;comment:登录用户名"`
	Name      string `json:"name" gorm:"type:varchar(32);not_null;comment:显示用户名"`
	Mobile    string `json:"mobile" gorm:"type:varchar(16);not_null;comment:用户手机号"`
	Password  string `json:"password" gorm:"type:varchar(128);not_null;comment:密码"`
	CommonTime
	DeleteTime
}

func (u User) GetUid() string {
	return strconv.Itoa(int(u.Id))
}
