package services

import (
	"gin-gorm-rest-api/config"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type UserService struct {
	databaseConfig *config.DatabaseConfig
}

func NewUserService(databaseConfig *config.DatabaseConfig) *UserService {
	return &UserService{databaseConfig: databaseConfig}
}

func (us *UserService) GetUser(id uuid.UUID) (models.User, error) {
	var user models.User
	result := us.databaseConfig.DB.Preload(clause.Associations).First(&user, "id = ?", id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	user.AfterLoad()
	return user, nil
}

func (us *UserService) GetUserBy(options dto.GetDataByOptions) (models.User, error) {
	var user models.User
	query := us.databaseConfig.DB.Preload(clause.Associations).Where(options.Field+" = ?", options.Value)
	if options.ExcludeId != nil {
		query = query.Where("id != ?", *options.ExcludeId)
	}
	result := query.First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	user.AfterLoad()
	return user, nil
}

func (us *UserService) GetPaginatedUsers(page int, limit int, search *string) ([]models.User, int, float64, error) {
	var users []models.User
	offset := (page - 1) * limit

	query := us.databaseConfig.DB.Model(&models.User{})
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	result := query.Preload(clause.Associations).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users)
	for i := 0; i < len(users); i++ {
		users[i].AfterLoad()
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return []models.User{}, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return users, int(count), totalPages, result.Error
}
