package routes

import (
	"gin-gorm-rest-api/handlers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	authMiddleware  *middlewares.AuthMiddleware
	authHandler     *handlers.AuthHandler
	authUserHandler *handlers.AuthUserHandler
}

func NewAuthRoute(authMiddleware *middlewares.AuthMiddleware, authHandler *handlers.AuthHandler, authUserHandler *handlers.AuthUserHandler) *AuthRoute {
	return &AuthRoute{authMiddleware, authHandler, authUserHandler}
}

func (ar *AuthRoute) Init(route *gin.Engine) {
	// auth route
	authRoute := route.Group("/auth")
	authRoute.POST("/login", ar.authHandler.Login)
	authRoute.POST("/register", ar.authHandler.Register)
	authRoute.GET("/refresh-token", ar.authHandler.RefreshToken)
	authRoute.POST("/forgot-password", ar.authHandler.ForgotPassword)
	authRoute.POST("/reset-password", ar.authHandler.ResetPassword)
	// auth user route
	authUserRoute := authRoute.Group("/user", ar.authMiddleware.Auth)
	authUserRoute.GET("/", ar.authUserHandler.AuthUser)
	authUserRoute.PATCH("/", ar.authUserHandler.UpdateAuthUser)
	authUserRoute.PATCH("/password", ar.authUserHandler.UpdateAuthUserPassword)
}
