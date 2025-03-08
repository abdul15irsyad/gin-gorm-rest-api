package services

import (
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/lib"
	"gin-gorm-rest-api/models"
	"math"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(libDB *lib.LibDatabase) *RoleService {
	return &RoleService{db: libDB.Database}
}

func (rs *RoleService) GetRole(id uuid.UUID) (models.Role, error) {
	var role models.Role
	result := rs.db.Preload(clause.Associations).First(&role, "id = ?", id)
	if result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (rs *RoleService) GetRoleBy(options dtos.GetDataByOptions) (models.Role, error) {
	var role models.Role
	query := rs.db.Preload(clause.Associations).Where(options.Field+" = ?", options.Value)
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

	query := rs.db.Model(&models.Role{})
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

func (rs *RoleService) CreateRole(createRoleDto dtos.CreateRoleDto) (models.Role, error) {
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return models.Role{}, err
	}
	role := models.Role{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      createRoleDto.Name,
		Slug:      slug.Make(createRoleDto.Name),
		Desc:      createRoleDto.Desc,
	}
	rs.db.Save(&role)
	role, _ = rs.GetRole(role.Id)

	return role, nil
}

func (rs *RoleService) UpdateRole(id uuid.UUID, updateRoleDto dtos.UpdateRoleDto) (models.Role, error) {
	var role models.Role
	role.Name = updateRoleDto.Name
	role.Slug = slug.Make(updateRoleDto.Name)
	role.Desc = updateRoleDto.Desc
	rs.db.Save(&role)
	role, err := rs.GetRole(role.Id)
	if err != nil {
		return models.Role{}, nil
	}

	return role, nil
}

func (rs *RoleService) DeleteRole(id uuid.UUID) {
	rs.db.Delete(&models.Role{}, id)
}
