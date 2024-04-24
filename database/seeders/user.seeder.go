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

func SeedUsers() {
	const name = "SeedUsers"
	var (
		seeder         models.Seeder
		newUuid        uuid.UUID
		hashedPassword []byte
	)
	result := database.DB.Where("name = ?", name).First(&seeder)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}

	newUuid, _ = uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
	hashedPassword, _ = utils.HashPassword("Qwerty123")
	var user1 models.User
	user1.Id = newUuid
	user1.Name = "irsyad"
	user1.Email = "irsyad@email.com"
	user1.Password = string(hashedPassword)

	users := []models.User{
		user1,
	}
	usersLength := len(users)
	for i := 1; i <= 30-usersLength; i++ {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := utils.HashPassword("Qwerty123")
		user := models.User{
			BaseModel: models.BaseModel{Id: randomUuid},
			Name:      "User " + fmt.Sprint(i),
			Email:     "user" + fmt.Sprint(i) + "@email.com",
			Password:  string(hashedPassword),
		}
		users = append(users, user)
	}

	result = database.DB.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	database.DB.Create(&models.Seeder{
		Name:      name,
		CreatedAt: time.Now(),
	})

	log.Println(fmt.Sprint(name) + " executed")
}
