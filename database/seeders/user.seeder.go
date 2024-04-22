package seeders

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedUsers() error {
	var seeder models.Seeder
	result := database.DB.Where("name = ?", "SeedUsers").First(&seeder)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	users := []models.User{
		{
			Name:     "irsyad",
			Email:    "irsyad@email.com",
			Password: "Qwerty123",
		},
	}
	for i := 1; i <= 30; i++ {
		user := models.User{
			Name:     "User " + fmt.Sprint(i),
			Email:    "user" + fmt.Sprint(i) + "@email.com",
			Password: "Qwerty123",
		}
		users = append(users, user)
	}

	for i := 0; i < len(users); i++ {
		hashedPassword, _ := utils.HashPassword(users[i].Password)
		randomUuid, _ := uuid.NewRandom()
		users[i].Id = randomUuid
		users[i].Password = string(hashedPassword)
		users[i].CreatedAt = time.Now()
		users[i].UpdatedAt = time.Now()
	}

	result = database.DB.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	database.DB.Create(&models.Seeder{
		Name:      "SeedUsers",
		CreatedAt: time.Now(),
	})

	log.Println("SeedUsers executed")
	return nil
}
