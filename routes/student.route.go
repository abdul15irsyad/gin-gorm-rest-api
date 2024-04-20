package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(route *gin.Engine) {
	studentRoute := route.Group("/students")
	studentRoute.GET("/", controllers.GetAllStudent)
	studentRoute.GET("/:id", controllers.GetStudent)
	studentRoute.POST("/", controllers.CreateStudent)
	studentRoute.PATCH("/:id", controllers.UpdateStudent)
	studentRoute.DELETE("/:id", controllers.DeleteStudent)
}
