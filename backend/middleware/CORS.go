package middleware

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/global"
	"net/http"
)

func AllowCors(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	if c.Request.Method == http.MethodOptions {
		global.Log.Info("Method Options, allow all CORS.")
		c.AbortWithStatus(http.StatusNoContent)
	}
}
