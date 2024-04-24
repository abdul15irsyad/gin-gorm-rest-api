package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine) {
	AuthRoute := route.Group("/auth")
	AuthRoute.POST("/login", controllers.Login)
	AuthRoute.POST("/register", controllers.Register)
	AuthRoute.GET("/user", middlewares.Auth, controllers.AuthUser)
	AuthRoute.GET("/refresh-token", controllers.RefreshToken)
}
