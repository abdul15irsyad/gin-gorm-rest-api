package services

import (
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"math"
)

func GetPaginatedUsers(page int, limit int, search *string) ([]models.User, int, float64, error) {
	var users []models.User
	offset := (page - 1) * limit

	query := database.DB.Model(&models.User{})
	if search != nil && *search != "" {
		query = query.Where("name ilike ? or email ilike ?", "%"+*search+"%", "%"+*search+"%")
	}
	result := query.Limit(limit).Offset(offset).Find(&users)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return users, int(count), totalPages, result.Error
}
