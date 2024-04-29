package database

import (
	"gin-gorm-rest-api/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	dsn := "host=" + DB_HOST + " user=" + DB_USER + " password=" + DB_PASS + " dbname=" + DB_NAME + " port=" + DB_PORT + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	} else {
		log.Println("database connected")
	}

	// auto migrate table
	allModels := []interface{}{
		models.Seeder{},
		models.File{},
		models.Role{},
		models.User{},
		models.Token{},
	}
	for _, model := range allModels {
		err = db.AutoMigrate(&model)
		if err != nil {
			panic(err.Error())
		}
	}

	DB = db
}
