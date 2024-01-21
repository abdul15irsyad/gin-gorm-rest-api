package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(route *gin.Engine) {
	routerGroup := route.Group("/students")
	routerGroup.GET("/", controllers.GetAllStudent)
	routerGroup.GET("/:id", controllers.GetStudent)
}
