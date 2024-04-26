package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(route *gin.Engine) {
	roleRoute := route.Group("/roles", middlewares.Auth)
	roleRoute.GET("/", controllers.GetAllRoles)
	roleRoute.GET("/:id", controllers.GetRole)
	// roleRoute.POST("/", controllers.CreateRole)
	// roleRoute.PATCH("/:id", controllers.UpdateRole)
	// roleRoute.DELETE("/:id", controllers.DeleteRole)
}
