package models

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Name     string     `json:"name" gorm:"not null"`
	Email    string     `json:"email" gorm:"not null;uniqueIndex:idx_users_email,where:deleted_at IS NULL"`
	Password string     `json:"-" gorm:"select:false;not null"`
	RoleId   uuid.UUID  `json:"-" gorm:"not null"`
	Role     *Role      `json:"role" gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
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
