package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func Student(route *gin.Engine) {
	routerGroup := route.Group("/students")
	routerGroup.GET("/", controllers.GetAllStudent)
	routerGroup.GET("/:id", controllers.GetStudent)
}
