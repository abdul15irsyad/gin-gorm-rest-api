package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine) {
	AuthRoute := route.Group("/auth")
	AuthRoute.POST("/login", controllers.Login)
	AuthRoute.POST("/register", controllers.Register)
}
