package routes

import (
	"gin-gorm-rest-api/handlers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	userHandler    *handlers.UserHandler
}

func NewUserRoute(authMiddleware *middlewares.AuthMiddleware, userHandler *handlers.UserHandler) *UserRoute {
	return &UserRoute{userHandler: userHandler, authMiddleware: authMiddleware}
}

func (ur *UserRoute) Init(route *gin.Engine) {
	userRoute := route.Group("/users", ur.authMiddleware.Auth)
	userRoute.GET("/", ur.userHandler.GetAllUsers)
	userRoute.GET("/:id", ur.userHandler.GetUser)
	userRoute.POST("/", ur.userHandler.CreateUser)
	userRoute.PATCH("/:id", ur.userHandler.UpdateUser)
	userRoute.DELETE("/:id", ur.userHandler.DeleteUser)
}
