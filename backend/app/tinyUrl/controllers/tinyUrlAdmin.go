package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"monaToolBox/app/tinyUrl/service"
	"monaToolBox/app/tinyUrl/types"
	"monaToolBox/global"
	"monaToolBox/global/response"
	"monaToolBox/global/validator"
	"strconv"
)

func List() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			global.Log.Warn("tinyUrl.GetInfo input id error.", zap.String("id", c.Param("id")))
			response.FailByError(c, global.Errors.ValidateError)
			return
		}
		err, tuInfo := service.TinyUrlService.GetById(id)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				global.Log.Error("service.TinyUrlService error.", zap.Error(err))
			}
			response.ServiceFail(c, gorm.ErrRecordNotFound.Error())
			return
		}

		output := types.TinyUrlInfo{
			Id:              tuInfo.Id,
			TinyUrl:         tuInfo.TinyUrl,
			OriginalUrl:     tuInfo.OriginalUrl,
			Pv:              tuInfo.Pv,
			CreateTime:      tuInfo.CreateTime,
			LimitAccessTime: tuInfo.LimitAccessTime,
		}
		response.Success(c, output)
	}
}

func Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.AddTinyUrlForm

		if err := c.ShouldBindJSON(&form); err != nil {
			response.ValidateFail(c, validator.GetErrorMsg(form, err))
			return
		}

	}
}

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
