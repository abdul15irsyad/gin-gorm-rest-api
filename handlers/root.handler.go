package handlers

import (
	"gin-gorm-rest-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootHandler struct {
	logService *services.LogService
}

func NewRootHandler(logService *services.LogService) *RootHandler {
	return &RootHandler{logService}
}

func (rh *RootHandler) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "gin rest api with gorm",
	})
}
