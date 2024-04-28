package models

import (
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	BaseModel
	Name     string     `json:"name" gorm:"not null"`
	Email    string     `json:"email" gorm:"not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL"`
	Password string     `json:"-" gorm:"select:false;not null"`
	RoleId   uuid.UUID  `json:"-" gorm:"not null"`
	Role     Role       `json:"role" gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ImageId  *uuid.UUID `json:"-"`
	Image    *File      `json:"image" gorm:"foreignKey:ImageId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (user *User) AfterLoad() {
	if user != nil {
		if user.ImageId != nil {
			user.Image.AfterLoad()
		}
	}
}

func GetUser(db *gorm.DB, id uuid.UUID) (User, error) {
	var user User
	result := db.Preload(clause.Associations).First(&user, "id = ?", id)
	if result.Error != nil {
		return User{}, result.Error
	}
	user.AfterLoad()
	return user, nil
}

func GetUserByField(db *gorm.DB, field string, value string, excludeId *string) (User, error) {
	var user User
	query := db.Preload(clause.Associations).Where(field+" = ?", value)
	if excludeId != nil {
		query = query.Where("id != ?", *excludeId)
	}
	result := query.First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	user.AfterLoad()
	return user, nil
}

func GetPaginatedUsers(db *gorm.DB, page int, limit int, search *string) ([]User, int, float64, error) {
	var users []User
	offset := (page - 1) * limit

	query := db.Model(&User{})
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	result := query.Preload(clause.Associations).Limit(limit).Offset(offset).Order("created_at DESC").Find(&users)
	for i := 0; i < len(users); i++ {
		users[i].AfterLoad()
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return []User{}, 0, 0, err
	}
	totalPages := math.Ceil(float64(count) / float64(limit))

	return users, int(count), totalPages, result.Error
}
