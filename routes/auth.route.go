package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine) {
	// auth route
	authRoute := route.Group("/auth")
	authRoute.POST("/login", controllers.Login)
	authRoute.POST("/register", controllers.Register)
	authRoute.GET("/refresh-token", controllers.RefreshToken)
	authRoute.POST("/forgot-password", controllers.ForgotPassword)
	authRoute.POST("/reset-password", controllers.ResetPassword)
	// auth user route
	authUserRoute := authRoute.Group("/user", middlewares.Auth)
	authUserRoute.GET("/", controllers.AuthUser)
	authUserRoute.PATCH("/", controllers.UpdateAuthUser)
	authUserRoute.PATCH("/password", controllers.UpdateAuthUserPassword)
}
