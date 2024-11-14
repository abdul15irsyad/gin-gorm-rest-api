package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type RoleRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	roleController *controllers.RoleController
}

func NewRoleRoute(authMiddleware *middlewares.AuthMiddleware, roleController *controllers.RoleController) *RoleRoute {
	return &RoleRoute{authMiddleware: authMiddleware, roleController: roleController}
}

func (rr *RoleRoute) Init(route *gin.Engine) {
	roleRoute := route.Group("/roles", rr.authMiddleware.Auth)
	roleRoute.GET("/", rr.roleController.GetAllRoles)
	roleRoute.GET("/:id", rr.roleController.GetRole)
	roleRoute.POST("/", rr.roleController.CreateRole)
	roleRoute.PATCH("/:id", rr.roleController.UpdateRole)
	roleRoute.DELETE("/:id", rr.roleController.DeleteRole)
}
