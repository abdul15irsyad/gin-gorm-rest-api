package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetFile(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getFileDto dto.GetFileDto
	ctx.ShouldBind(&getFileDto)
	getFileDto.Id = paramId
	validationErrors := utils.Validate(getFileDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	file, err := services.GetFile(database.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get file",
		"data":    file,
	})
}

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
	files, total, totalPage, err := services.GetPaginatedFiles(database.DB, page, limit, &search)
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

	newFile, ok := services.UploadAndCreateFile(ctx, file, database.DB)
	if !ok {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create file",
		"data":    newFile,
	})
}

func DeleteFile(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getFileDto dto.GetFileDto
	ctx.ShouldBind(&getFileDto)
	getFileDto.Id = paramId
	validationErrors := utils.Validate(getFileDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	file, err := services.GetFile(database.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	database.DB.Delete(&file)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete file",
	})
}
