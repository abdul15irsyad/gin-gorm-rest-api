package models

import "time"

type Student struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Year      int64     `json:"year"`
	CreatedAt time.Time	`json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
