package models

type Role struct {
	BaseModel
	Name string  `json:"name" gorm:"not null"`
	Slug string  `json:"slug" gorm:"not null;uniqueIndex:idx_roles_slug,where:deleted_at IS NULL"`
	Desc *string `json:"desc" gorm:"comment:description"`
}
