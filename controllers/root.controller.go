package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Gin REST API with GORM",
	})
}
