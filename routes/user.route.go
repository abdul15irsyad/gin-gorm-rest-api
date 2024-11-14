package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	userController *controllers.UserController
}

func NewUserRoute(authMiddleware *middlewares.AuthMiddleware, userController *controllers.UserController) *UserRoute {
	return &UserRoute{userController: userController, authMiddleware: authMiddleware}
}

func (ur *UserRoute) Init(route *gin.Engine) {
	userRoute := route.Group("/users", ur.authMiddleware.Auth)
	userRoute.GET("/", ur.userController.GetAllUsers)
	userRoute.GET("/:id", ur.userController.GetUser)
	userRoute.POST("/", ur.userController.CreateUser)
	userRoute.PATCH("/:id", ur.userController.UpdateUser)
	userRoute.DELETE("/:id", ur.userController.DeleteUser)
}
