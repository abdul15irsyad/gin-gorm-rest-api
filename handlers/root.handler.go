package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootHandler struct{}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (rc *RootHandler) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "gin rest api with gorm",
	})
}
