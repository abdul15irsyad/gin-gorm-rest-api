package routes

import (
	"belajar-gin/controllers"

	"github.com/gin-gonic/gin"
)

func Root( route *gin.Engine) {
	route.GET("/", controllers.Root)
}
