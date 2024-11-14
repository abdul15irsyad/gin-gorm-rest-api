package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

func (rc *RootController) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "gin rest api with gorm",
	})
}
