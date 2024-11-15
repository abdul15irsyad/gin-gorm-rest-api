package routes

import (
	"gin-gorm-rest-api/handlers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type FileRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	fileHandler    *handlers.FileHandler
}

func NewFileRoute(authMiddleware *middlewares.AuthMiddleware, fileHandler *handlers.FileHandler) *FileRoute {
	return &FileRoute{authMiddleware: authMiddleware, fileHandler: fileHandler}
}

func (fr *FileRoute) Init(route *gin.Engine) {
	fileRoute := route.Group("/files", fr.authMiddleware.Auth)
	fileRoute.GET("/", fr.fileHandler.GetAllFiles)
	fileRoute.GET("/:id", fr.fileHandler.GetFile)
	fileRoute.POST("/", fr.fileHandler.CreateFile)
	fileRoute.DELETE("/:id", fr.fileHandler.DeleteFile)
}
