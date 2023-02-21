package userCenter

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/app/userCenter/controllers"
	"monaToolBox/middleware"
	"net/http"
)

func UserCenterRouter(root *gin.RouterGroup) {
	// usercenter

	root.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})

	root.POST("/register", controllers.Register)
	root.POST("/login", controllers.Login)

	userLoggedGroup := root.Group("").Use(middleware.JwtAuth)
	userLoggedGroup.POST("/logout", controllers.Logout)
	userLoggedGroup.GET("/userInfo", controllers.GetUserInfo)

	//userLoggedGroup := root.Group(":userId").Use(middleware.JwtAuth)
	//userLoggedGroup.GET("/userInfo", controllers.GetUserInfo)
	//userLoggedGroup.POST("/logout", controllers.Logout)

}
