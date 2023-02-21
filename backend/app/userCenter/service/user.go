package service

import (
	"errors"
	"monaToolBox/app/userCenter/types"
	"monaToolBox/global"
	"monaToolBox/models"
	"monaToolBox/utils"
)

type userService struct{}

var UserService = new(userService)

// Register 注册
func (userService *userService) Register(params types.Register) (err error, user models.User) {
	mobileCheck := global.DB.Where("mobile = ?", params.Mobile).Select("id").First(&models.User{})

	if mobileCheck.RowsAffected != 0 {
		err = errors.New("手机号已存在")
		return
	}

	loginNameCheck := global.DB.Where("loginName = ?", params.LoginName).Select("id").First(&models.User{})
	if loginNameCheck.RowsAffected != 0 {
		err = errors.New("登录用户名已存在")
		return
	}

	user = models.User{
		LoginName: params.LoginName,
		Name:      params.Name,
		Mobile:    params.Mobile,
		Password:  utils.BcryptMake([]byte(params.Password)),
	}

	err = global.DB.Create(&user).Error
	return
}

// Login 登录service
func (userService *userService) Login(params types.LoginForm) (err error, user *models.User) {
	err = global.DB.Where(" login_name = ? ", params.LoginName).First(&user).Error

	if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New("用户名不存在或密码错误")
	}

	return
}

func (userService *userService) GetUserInfoById(userId int) (err error, user models.User) {
	err = global.DB.First(&user, userId).Error
	if err != nil {
		err = errors.New("未找到该用户")
	}
	return
}
