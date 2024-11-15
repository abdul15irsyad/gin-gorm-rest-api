package routes

import (
	"gin-gorm-rest-api/handlers"

	"github.com/gin-gonic/gin"
)

type RootRoute struct {
	rootHandler *handlers.RootHandler
}

func NewRootRoute(rootHandler *handlers.RootHandler) *RootRoute {
	return &RootRoute{rootHandler: rootHandler}
}

func (rr *RootRoute) Init(route *gin.Engine) {
	route.GET("/", rr.rootHandler.Root)
}
