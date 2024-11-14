package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	authMiddleware     *middlewares.AuthMiddleware
	authController     *controllers.AuthController
	authUserController *controllers.AuthUserController
}

func NewAuthRoute(authMiddleware *middlewares.AuthMiddleware, authController *controllers.AuthController, authUserController *controllers.AuthUserController) *AuthRoute {
	return &AuthRoute{authMiddleware: authMiddleware, authController: authController, authUserController: authUserController}
}

func (ar *AuthRoute) Init(route *gin.Engine) {
	// auth route
	authRoute := route.Group("/auth")
	authRoute.POST("/login", ar.authController.Login)
	authRoute.POST("/register", ar.authController.Register)
	authRoute.GET("/refresh-token", ar.authController.RefreshToken)
	authRoute.POST("/forgot-password", ar.authController.ForgotPassword)
	authRoute.POST("/reset-password", ar.authController.ResetPassword)
	// auth user route
	authUserRoute := authRoute.Group("/user", ar.authMiddleware.Auth)
	authUserRoute.GET("/", ar.authUserController.AuthUser)
	authUserRoute.PATCH("/", ar.authUserController.UpdateAuthUser)
	authUserRoute.PATCH("/password", ar.authUserController.UpdateAuthUserPassword)
}
