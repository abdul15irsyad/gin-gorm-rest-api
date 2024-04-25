package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine) {
	authRoute := route.Group("/auth")
	authRoute.POST("/login", controllers.Login)
	authRoute.POST("/register", controllers.Register)
	authRoute.GET("/user", middlewares.Auth, controllers.AuthUser)
	authRoute.GET("/refresh-token", controllers.RefreshToken)
}
