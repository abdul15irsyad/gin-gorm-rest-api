package services

import (
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/lib"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(libDB *lib.LibDatabase) *UserService {
	return &UserService{db: libDB.Database}
}

func (us *UserService) GetUser(id uuid.UUID) (models.User, error) {
	var user models.User
	result := us.db.Preload(clause.Associations).First(&user, "id = ?", id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	user.AfterLoad()
	return user, nil
}

func (us *UserService) GetUserBy(options dtos.GetDataByOptions) (models.User, error) {
	var user models.User
	query := us.db.Preload(clause.Associations).Where(options.Field+" = ?", options.Value)
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

func (us *UserService) GetUserCredential(email string) (models.User, error) {
	var authUser models.User
	result := us.db.Select([]string{"id", "password"}).First(&authUser, "email = ?", email)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return authUser, nil
}

func (us *UserService) GetPaginatedUsers(options dtos.GetUsersDto) ([]models.User, int, float64, error) {
	page := *options.Page
	limit := *options.Limit
	search := options.Search
	var users []models.User
	offset := (page - 1) * limit

	query := us.db.Model(&models.User{})
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

func (us *UserService) CreateUser(createUserDto dtos.CreateUserDto) (models.User, error) {
	hashedPassword, err := utils.HashPassword(createUserDto.Password)
	if err != nil {
		return models.User{}, err
	}
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return models.User{}, err
	}

	roleId, _ := uuid.Parse(createUserDto.RoleId)

	user := models.User{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  hashedPassword,
		RoleId:    roleId,
	}
	result := us.db.Save(&user)
	if result.Error != nil {
		return models.User{}, err
	}

	user, _ = us.GetUser(user.Id)
	return user, nil
}

func (us *UserService) UpdateUser(id uuid.UUID, updateUserDto dtos.UpdateUserDto, newImageFileId *uuid.UUID) (models.User, error) {
	var user models.User
	user.Name = updateUserDto.Name
	user.Email = updateUserDto.Email
	roleId, _ := uuid.Parse(updateUserDto.RoleId)
	user.RoleId = roleId
	if newImageFileId != nil {
		user.ImageId = newImageFileId
	}
	if updateUserDto.Password != nil {
		hashedPassword, err := utils.HashPassword(*updateUserDto.Password)
		if err != nil {
			return models.User{}, err
		}
		user.Password = hashedPassword
	}
	us.db.Save(&user)
	user, err := us.GetUser(user.Id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (us *UserService) UpdateUserPassword(user models.User, newPassword string) error {
	hashedPassword, _ := utils.HashPassword(newPassword)
	user.Password = hashedPassword
	result := us.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (us *UserService) DeleteUser(id uuid.UUID) {
	us.db.Delete(&models.User{}, id)
}
