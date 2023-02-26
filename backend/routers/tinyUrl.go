package routers

import (
	"github.com/gin-gonic/gin"
	"monaToolBox/app/tinyUrl/controllers"
)

func tinyUrlAdminRouter(root *gin.RouterGroup) {
	// /tinyUrl
	root.GET("/list")
	root.GET("/:id", controllers.GetInfo())

	root.POST("/add")
	root.POST("/modify/:id")

	root.DELETE("/:id")

}

func tinyUrlRouter(root *gin.RouterGroup) {
	root.GET("/:tinyUrl")
}
