package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"monaToolBox/app/userCenter/service"
	"monaToolBox/app/userCenter/types"
	"monaToolBox/global"
	"monaToolBox/global/response"
	"monaToolBox/global/validator"
	"strconv"
)

// Register 注册接口
func Register(c *gin.Context) {
	var form types.Register

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, validator.GetErrorMsg(form, err))
		return
	}

	if err, user := service.UserService.Register(form); err != nil {
		global.Log.Error(err.Error())
		response.ServiceFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// Login 登录接口
func Login(c *gin.Context) {
	var form types.LoginForm

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, validator.GetErrorMsg(form, err))
		return
	}

	if err, user := service.UserService.Login(form); err != nil {
		response.ServiceFail(c, err.Error())
	} else {
		tokenData, err := service.JwtService.CreateToken(user)
		if err != nil {
			response.ServiceFail(c, err.Error())
			return
		}
		loginResp := types.LoginResp{
			UserId:      user.Id,
			AccessToken: tokenData.AccessToken,
		}
		response.Success(c, loginResp)
	}
}

// Logout 用户登出接口
func Logout(c *gin.Context) {
	err := service.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.ServiceFail(c, "登出失败")
		return
	}
	response.Success(c, nil)
}

// GetUserInfo 获取用户信息接口
func GetUserInfo(c *gin.Context) {
	userId, err := strconv.Atoi(c.Keys["userId"].(string))
	if err != nil {
		response.ClaimsTokenFail(c)
		return
	}

	err, user := service.UserService.GetUserInfoById(userId)
	if err != nil {
		response.ServiceFail(c, err.Error())
		return
	}
	response.Success(c, user)
}
