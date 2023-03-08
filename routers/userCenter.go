package routers

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/app/userCenter/controllers"
)

func userCenterRouter(root *gin.RouterGroup) {
	// userCenter

	root.POST("/register", controllers.Register)
	root.POST("/login", controllers.Login)

}

func userCenterAdminRouter(root *gin.RouterGroup) {
	root.POST("/logout", controllers.Logout)
	root.GET("/userInfo", controllers.GetUserInfo)
}
