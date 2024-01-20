package configs

import (
	"belajar-gin/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database() {
	dsn := "host=localhost user=postgres password=Starmoon15 dbname=belajar_gin port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}
	
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		panic("Failed migrate!")
	}

	DB = db
}