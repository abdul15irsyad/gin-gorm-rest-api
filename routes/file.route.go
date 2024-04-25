package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func FileRoutes(route *gin.Engine) {
	fileRoute := route.Group("/file", middlewares.Auth)
	fileRoute.GET("/", controllers.GetAllFiles)
	fileRoute.POST("/", controllers.CreateFile)
}
