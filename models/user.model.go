package models

import (
	"math"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"-" gorm:"select:false;not null"`
}

func GetPaginatedUsers(db *gorm.DB, page int, limit int, search *string) ([]User, int, float64, error) {
	var users []User
	offset := (page - 1) * limit

	query := db.Model(&User{})
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
