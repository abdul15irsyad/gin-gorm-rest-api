package seeders

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FileSeeder() {
	// check seeder
	const name = "FileSeeder"
	var (
		seeder  models.Seeder
		newUuid uuid.UUID
	)
	result := database.DB.Where("name = ?", name).First(&seeder)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}

	// seeder data
	uuids := []string{
		"4655d7f8-b5af-4b03-88d7-907093a358ed", "33d7a6cd-6d03-4eb6-b05c-293633c2a01d", "abf3104c-26e3-4853-aa3a-1f7b9a8cd71e", "563c38fa-17ba-4512-9f31-59eeab9e136f", "82ec1902-a999-476a-99de-83c0978e6b76",
	}
	var files []models.File
	for i := 1; i <= 5; i++ {
		newUuid, _ = uuid.Parse(uuids[i-1])
		file := models.File{
			BaseModel:        models.BaseModel{Id: newUuid},
			Path:             "/dummy",
			Filename:         "dummy-profile-" + fmt.Sprint(i) + ".jpg",
			OriginalFilename: "dummy-profile-" + fmt.Sprint(i) + ".jpg",
			Mime:             "image/jpeg",
		}
		files = append(files, file)
	}
	result = database.DB.Create(&files)
	if result.Error != nil {
		panic(result.Error)
	}

	// add to seeder table
	database.DB.Create(&models.Seeder{
		Name:      name,
		CreatedAt: time.Now(),
	})
	log.Println(fmt.Sprint(name) + " executed")
}
