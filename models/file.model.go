package models

import (
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	BaseModel
	Path             string `json:"path" gorm:"not null"`
	Filename         string `json:"filename" gorm:"not null"`
	OriginalFilename string `json:"originalFilename" gorm:"not null"`
	Mime             string `json:"mime" gorm:"not null"`
	Url              string `json:"url" gorm:"-"`
}

func (file *File) AfterLoad() {
	file.Url = "/assets" + file.Path + "/" + file.Filename
}

func GetPaginatedFiles(db *gorm.DB, page int, limit int, search *string) ([]File, int, float64, error) {
	var files []File
	offset := (page - 1) * limit

	query := db.Model(&File{})
	if search != nil && *search != "" {
		query = query.Where("name ilike ? or email ilike ?", "%"+*search+"%", "%"+*search+"%")
	}
	result := query.Limit(limit).Offset(offset).Find(&files)
	for i := 0; i < len(files); i++ {
		files[i].AfterLoad()
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return files, int(count), totalPages, result.Error
}

func UploadAndCreateFile(ctx *gin.Context, file *multipart.FileHeader, db *gorm.DB) (*File, bool) {
	filename := fmt.Sprint(time.Now().UnixMicro()) + "-" + file.Filename
	err := ctx.SaveUploadedFile(file, "./assets/uploads/"+filename)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return nil, false
	}
	mime := file.Header.Get("Content-Type")
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return nil, false
	}
	newFile := File{
		BaseModel:        BaseModel{Id: randomUuid},
		Path:             "/uploads",
		Filename:         filename,
		OriginalFilename: file.Filename,
		Mime:             mime,
	}
	result := db.Save(&newFile)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return nil, false
	}
	newFile.AfterLoad()

	return &newFile, true
}
