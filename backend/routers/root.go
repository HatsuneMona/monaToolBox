package routers

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/middleware"
	"monaToolBox/routers/userCenter"
	"net/http"
)

func GinRootRouter() *gin.Engine {
	root := gin.Default()

	root.Use(middleware.AllowCors)

	root.GET(
		"/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		},
	)

	userCenter.UserCenterRouter(root.Group("/usercenter"))

	return root
}
