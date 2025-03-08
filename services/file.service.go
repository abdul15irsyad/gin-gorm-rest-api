package services

import (
	"fmt"
	"gin-gorm-rest-api/lib"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"math"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileService struct {
	db *gorm.DB
}

func NewFileService(libDB *lib.LibDatabase) *FileService {
	return &FileService{db: libDB.Database}
}

func (fs *FileService) GetFile(id uuid.UUID) (models.File, error) {
	var file models.File
	result := fs.db.Preload(clause.Associations).First(&file, "id = ?", id)
	if result.Error != nil {
		return models.File{}, result.Error
	}
	file.AfterLoad()
	return file, nil
}

func (fs *FileService) GetPaginatedFiles(page int, limit int, search *string) ([]models.File, int, float64, error) {
	var files []models.File
	offset := (page - 1) * limit

	query := fs.db.Model(&models.File{})
	if search != nil && *search != "" {
		query = query.Where("filename ILIKE ? OR original_filename ILIKE ? OR mime ILIKE ?", "%"+*search+"%", "%"+*search+"%", "%"+*search+"%")
	}
	result := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&files)
	for i := 0; i < len(files); i++ {
		files[i].AfterLoad()
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return []models.File{}, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return files, int(count), totalPages, result.Error
}

func (fs *FileService) UploadAndCreateFile(ctx *gin.Context, file *multipart.FileHeader) (models.File, bool) {
	filename := fmt.Sprint(time.Now().UnixMicro()) + "-" + utils.Slugify(strings.Split(file.Filename, ".")[0]) + "." + strings.Split(file.Filename, ".")[1]
	err := ctx.SaveUploadedFile(file, "./assets/uploads/"+filename)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return models.File{}, false
	}
	mime := file.Header.Get("Content-Type")
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return models.File{}, false
	}
	newFile := models.File{
		BaseModel:        models.BaseModel{Id: randomUuid},
		Path:             "/uploads",
		Filename:         filename,
		OriginalFilename: file.Filename,
		Mime:             mime,
	}
	result := fs.db.Save(&newFile)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return models.File{}, false
	}
	newFile.AfterLoad()

	return newFile, true
}

func (fs *FileService) DeleteFile(id uuid.UUID) {
	fs.db.Delete(&models.User{}, id)
}
