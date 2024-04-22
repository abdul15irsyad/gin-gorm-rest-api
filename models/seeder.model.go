package models

import "time"

type Seeder struct {
	Id        uint `gorm:"primary_key;auto_increment"`
	Name      string
	CreatedAt time.Time
}
