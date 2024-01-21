package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func RootRoutes(route *gin.Engine) {
	route.GET("/", controllers.Root)
}
