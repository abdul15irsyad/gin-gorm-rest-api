package v1

import (
	"belajar-gin/handlers"

	"github.com/gin-gonic/gin"
)

func Root(route *gin.RouterGroup) {
	route.GET("/", handlers.RootHandler)
}
