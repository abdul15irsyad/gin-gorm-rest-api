package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (rh *RootHandler) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "gin rest api with gorm",
	})
}

func (rh *RootHandler) Error(ctx *gin.Context) {
	panic(errors.New("error example"))
}
