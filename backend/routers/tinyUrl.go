package routers

import "github.com/gin-gonic/gin"

func tinyUrlAdminRouter(root *gin.RouterGroup) {
	root.GET("/list")
	root.GET("/:id")

	root.POST("/add")
	root.POST("/modify/:id")

	root.DELETE("/:id")

}

func tinyUrlRouter(root *gin.RouterGroup) {
	root.GET("/:tinyUrl")
}
