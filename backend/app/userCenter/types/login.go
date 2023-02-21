package types

import "monaToolBox/global/validator"

type LoginForm struct {
	LoginName string `json:"loginName" form:"loginName" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
}

func (l LoginForm) GetMessages() validator.ValidatorMassages {
	return validator.ValidatorMassages{
		"loginName.required": "登录用户名不能为空",
		"password.required":  "密码不能为空",
	}
}

type LoginResp struct {
	UserId      uint   `json:"user_id"`
	AccessToken string `json:"access_token"`
}
