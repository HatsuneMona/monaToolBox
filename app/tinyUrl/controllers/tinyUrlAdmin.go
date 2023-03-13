package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	. "monaToolBox/app/tinyUrl/service"
	"monaToolBox/app/tinyUrl/types"
	"monaToolBox/global"
	"monaToolBox/global/response"
	"monaToolBox/global/validator"
	"monaToolBox/models"
	"strconv"
	"time"
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
			response.FailByError(c, global.HandlerErrors.ValidateError)
			return
		}
		err, tuInfo := TinyUrlService.GetById(id)
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

		tinyUrlInfo := models.TinyUrl{
			TinyUrl:         form.TinyUrl,
			OriginalUrl:     form.OriginalUrl,
			LimitAccessTime: form.LimitAccessTime,
		}
		if tinyUrlInfo.LimitAccessTime.IsZero() {
			tinyUrlInfo.LimitAccessTime = time.Date(2099, 01, 01, 00, 00, 00, 00, time.Local)
		}

		// 检查存在性
		if err, tinyUrlInfos := TinyUrlService.GetByTinyRouteList([]string{form.TinyUrl}, true); err != nil {
			global.Log.Error("TinyUrlService.GetByTinyRouteList server error.", zap.Error(err), zap.Strings("input", []string{form.TinyUrl}))
			response.ServiceFail(c)
			return

		} else if len(tinyUrlInfos) > 0 {
			response.FailByError(c, types.HandlerErrors.Conflict, form.TinyUrl)
			return
		}

		err, tinyUrls := TinyUrlService.AddNewTinyUrlBatch([]models.TinyUrl{tinyUrlInfo})
		if err != nil {
			global.Log.Error("TinyUrlService.AddNewTinyUrlBatch server error.", zap.Error(err), zap.Any("input", []string{form.TinyUrl}))
			response.FailByError(c, global.ServiceErrors.ServiceError)
			return
		}

		response.Success(c, types.AddTinyUrlResp{Id: tinyUrls[0].Id})
		return
	}
}

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
