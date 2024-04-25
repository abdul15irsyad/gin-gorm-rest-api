package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	userRoute := route.Group("/users", middlewares.Auth)
	userRoute.GET("/", controllers.GetAllUsers)
	userRoute.GET("/:id", controllers.GetUser)
	userRoute.POST("/", controllers.CreateUser)
	userRoute.PATCH("/:id", controllers.UpdateUser)
	userRoute.DELETE("/:id", controllers.DeleteUser)
}
