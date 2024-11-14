package routes

import (
	"gin-gorm-rest-api/controllers"

	"github.com/gin-gonic/gin"
)

type RootRoute struct {
	rootController *controllers.RootController
}

func NewRootRoute(rootController *controllers.RootController) *RootRoute {
	return &RootRoute{rootController: rootController}
}

func (rr *RootRoute) Init(route *gin.Engine) {
	route.GET("/", rr.rootController.Root)
}
