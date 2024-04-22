package routes

import (
	"gin-gorm-rest-api/controllers"

	"github.com/gin-gonic/gin"
)

func RootRoutes(route *gin.Engine) {
	route.GET("/", controllers.Root)
}
