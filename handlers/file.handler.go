package handlers

import (
	"errors"
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHandler struct {
	fileService    *services.FileService
	databaseConfig *configs.DatabaseConfig
}

func NewFileHandler(fileService *services.FileService, databaseConfig *configs.DatabaseConfig) *FileHandler {
	return &FileHandler{fileService: fileService, databaseConfig: databaseConfig}
}

func (fc *FileHandler) GetFile(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getFileDto dtos.GetFileDto
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
	file, err := fc.fileService.GetFile(id)
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

func (fc *FileHandler) GetAllFiles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	files, total, totalPage, err := fc.fileService.GetPaginatedFiles(page, limit, &search)
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

func (fc *FileHandler) CreateFile(ctx *gin.Context) {
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

	newFile, ok := fc.fileService.UploadAndCreateFile(ctx, file)
	if !ok {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create file",
		"data":    newFile,
	})
}

func (fc *FileHandler) DeleteFile(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getFileDto dtos.GetFileDto
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
	file, err := fc.fileService.GetFile(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	fc.databaseConfig.DB.Delete(&file)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete file",
	})
}
