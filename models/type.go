package models

import "gorm.io/gorm"

type GetDataByOptions struct {
	DB        *gorm.DB
	Field     string
	Value     string
	ExcludeId *string
}
