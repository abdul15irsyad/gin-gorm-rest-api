package routes

import (
	"gin-gorm-rest-api/handlers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type RoleRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	roleHandler    *handlers.RoleHandler
}

func NewRoleRoute(authMiddleware *middlewares.AuthMiddleware, roleHandler *handlers.RoleHandler) *RoleRoute {
	return &RoleRoute{authMiddleware: authMiddleware, roleHandler: roleHandler}
}

func (rr *RoleRoute) Init(route *gin.Engine) {
	roleRoute := route.Group("/roles", rr.authMiddleware.Auth)
	roleRoute.GET("/", rr.roleHandler.GetAllRoles)
	roleRoute.GET("/:id", rr.roleHandler.GetRole)
	roleRoute.POST("/", rr.roleHandler.CreateRole)
	roleRoute.PATCH("/:id", rr.roleHandler.UpdateRole)
	roleRoute.DELETE("/:id", rr.roleHandler.DeleteRole)
}
