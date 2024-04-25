package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllFiles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	files, total, totalPage, err := models.GetPaginatedFiles(database.DB, page, limit, &search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get all file",
		"data":    files,
		"meta": gin.H{
			"currentPage": page,
			"pageSize":    limit,
			"totalPages":  totalPage,
			"totalItems":  total,
		},
	})
}

func CreateFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	newFile, ok := models.UploadAndCreateFile(ctx, file, database.DB)
	if !ok {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create file",
		"data":    *newFile,
	})
}
