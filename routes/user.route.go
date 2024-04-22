package routes

import (
	"gin-gorm-rest-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	UserRoute := route.Group("/users")
	UserRoute.GET("/", controllers.GetAllUser)
	UserRoute.GET("/:id", controllers.GetUser)
	UserRoute.POST("/", controllers.CreateUser)
	UserRoute.PATCH("/:id", controllers.UpdateUser)
	UserRoute.DELETE("/:id", controllers.DeleteUser)
}
