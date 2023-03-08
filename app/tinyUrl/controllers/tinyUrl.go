package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	. "monaToolBox/app/tinyUrl/service"
	"monaToolBox/app/tinyUrl/types"
	"monaToolBox/global"
	"monaToolBox/global/response"
	"net/http"
)

// TinyUrlRedirect 访问短链，进行重定向
func TinyUrlRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		tinyStr := c.Param("tinyUrl")

		if tinyStr == "" {
			response.FailByError(c, types.HandlerErrors.NotFound)
			return
		}

		err, tinyList := TinyUrlService.GetByTinyRouteList([]string{tinyStr})
		if err != nil {
			global.Log.Error("TinyUrlService.GetByTinyRouteList server error.", zap.Error(err), zap.Strings("input", []string{tinyStr}))
			response.ServiceFail(c)
			return
		} else if len(tinyList) < 1 {
			response.FailByError(c, types.HandlerErrors.NotFound)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, tinyList[0].OriginalUrl)
	}
}
