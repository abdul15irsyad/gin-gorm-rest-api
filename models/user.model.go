package models

type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"-" gorm:"select:false;not null"`
}
