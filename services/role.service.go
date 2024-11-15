package services

import (
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/models"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type RoleService struct {
	databaseConfig *configs.DatabaseConfig
}

func NewRoleService(databaseConfig *configs.DatabaseConfig) *RoleService {
	return &RoleService{databaseConfig}
}

func (rs *RoleService) GetRole(id uuid.UUID) (models.Role, error) {
	var role models.Role
	result := rs.databaseConfig.DB.Preload(clause.Associations).First(&role, "id = ?", id)
	if result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (rs *RoleService) GetRoleBy(options dtos.GetDataByOptions) (models.Role, error) {
	var role models.Role
	query := rs.databaseConfig.DB.Preload(clause.Associations).Where(options.Field+" = ?", options.Value)
	if options.ExcludeId != nil {
		query = query.Where("id != ?", *options.ExcludeId)
	}
	result := query.First(&role)
	if result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (rs *RoleService) GetPaginatedRoles(page int, limit int, search *string) ([]models.Role, int, float64, error) {
	var roles []models.Role
	offset := (page - 1) * limit

	query := rs.databaseConfig.DB.Model(&models.Role{})
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
	}
	result := query.Preload(clause.Associations).Limit(limit).Offset(offset).Order("created_at DESC").Find(&roles)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return []models.Role{}, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return roles, int(count), totalPages, result.Error
}
