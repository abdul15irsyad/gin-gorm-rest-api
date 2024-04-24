package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	UserRoute := route.Group("/users", middlewares.Auth)
	UserRoute.GET("/", controllers.GetAllUser)
	UserRoute.GET("/:id", controllers.GetUser)
	UserRoute.POST("/", controllers.CreateUser)
	UserRoute.PATCH("/:id", controllers.UpdateUser)
	UserRoute.DELETE("/:id", controllers.DeleteUser)
}
