package types

import (
	"monaToolBox/global/validator"
)

type Register struct {
	LoginName string `json:"loginName" form:"loginName" binding:"required"`
	Name      string `json:"name" form:"name" binding:"required"`
	Mobile    string `json:"mobile" form:"mobile" binding:"mobile"`
	Password  string `json:"password" form:"name" binding:"required"`
}

func (r Register) GetMessages() validator.ValidatorMassages {
	return validator.ValidatorMassages{
		"loginName.required": "登录用户名不能为空",
		"name.required":      "用户昵称不能为空",
		"mobile.mobile":      "手机号码不符合规范",
		"password.required":  "密码不能为空",
	}
}
