package routers

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/middleware"
	"net/http"
)

func GinRootRouter() *gin.Engine {
	root := gin.Default()

	root.Use(middleware.GinCors())
	// root.Use(middleware.AllowCors)

	root.GET(
		"/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		},
	)

	adminRouter(root.Group("/admin"))

	tinyUrlRouter(root.Group("/tu"))
	userCenterRouter(root.Group("/userCenter"))

	return root
}

func adminRouter(root *gin.RouterGroup) {
	root.Use(middleware.AdminOperateLog())
	root.Use(middleware.JwtAuth)

	userCenterAdminRouter(root.Group("/userCenter"))
	tinyUrlAdminRouter(root.Group("/tinyUrl"))

}
