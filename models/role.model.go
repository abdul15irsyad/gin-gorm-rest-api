package models

import (
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Role struct {
	BaseModel
	Name string  `json:"name" gorm:"not null"`
	Slug string  `json:"slug" gorm:"index;not null"`
	Desc *string `json:"desc" gorm:"comment:description"`
}

func GetRoleBySlug(roles *[]Role, slug string) (Role, bool) {
	for _, role := range *roles {
		if role.Slug == slug {
			return role, true
		}
	}
	return Role{}, false
}

func GetRole(db *gorm.DB, id uuid.UUID) (Role, error) {
	var role Role
	result := db.Preload(clause.Associations).First(&role, "id = ?", id)
	if result.Error != nil {
		return Role{}, result.Error
	}
	return role, nil
}

func GetPaginatedRoles(db *gorm.DB, page int, limit int, search *string) ([]Role, int, float64, error) {
	var roles []Role
	offset := (page - 1) * limit

	query := db.Model(&Role{})
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}
	result := query.Preload(clause.Associations).Limit(limit).Offset(offset).Order("created_at DESC").Find(&roles)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return []Role{}, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return roles, int(count), totalPages, result.Error
}
