package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func FileRoutes(route *gin.Engine) {
	fileRoute := route.Group("/files", middlewares.Auth)
	fileRoute.GET("/", controllers.GetAllFiles)
	fileRoute.GET("/:id", controllers.GetFile)
	fileRoute.POST("/", controllers.CreateFile)
	fileRoute.DELETE("/:id", controllers.DeleteFile)
}
