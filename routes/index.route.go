package routes

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {
	RootRoutes(router)
	StudentRoutes(router)
	UserRoutes(router)
}
