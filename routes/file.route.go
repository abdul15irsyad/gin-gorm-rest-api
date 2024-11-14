package routes

import (
	"gin-gorm-rest-api/controllers"
	"gin-gorm-rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

type FileRoute struct {
	authMiddleware *middlewares.AuthMiddleware
	fileController *controllers.FileController
}

func NewFileRoute(authMiddleware *middlewares.AuthMiddleware, fileController *controllers.FileController) *FileRoute {
	return &FileRoute{authMiddleware: authMiddleware, fileController: fileController}
}

func (fr *FileRoute) Init(route *gin.Engine) {
	fileRoute := route.Group("/files", fr.authMiddleware.Auth)
	fileRoute.GET("/", fr.fileController.GetAllFiles)
	fileRoute.GET("/:id", fr.fileController.GetFile)
	fileRoute.POST("/", fr.fileController.CreateFile)
	fileRoute.DELETE("/:id", fr.fileController.DeleteFile)
}
