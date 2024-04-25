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

func UserSeeder() {
	const name = "UserSeeder"
	var (
		seeder         models.Seeder
		newUuid        uuid.UUID
		hashedPassword []byte
	)
	result := database.DB.Where("name = ?", name).First(&seeder)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}

	var files []models.File
	database.DB.Find(&files)

	users := []models.User{}
	newUuid, _ = uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
	hashedPassword, _ = utils.HashPassword("Qwerty123")
	users = append(users, models.User{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "irsyad",
		Email:     "irsyad@email.com",
		Password:  string(hashedPassword),
		ImageId:   &files[0].Id,
	})
	usersLength := len(users)
	for i := 1; i <= 30-usersLength; i++ {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := utils.HashPassword("Qwerty123")
		randomFile := *utils.RandomSlice(files)
		user := models.User{
			BaseModel: models.BaseModel{Id: randomUuid},
			Name:      "User " + fmt.Sprint(i),
			Email:     "user" + fmt.Sprint(i) + "@email.com",
			Password:  string(hashedPassword),
		}
		if *utils.RandomSlice([]bool{true, false}) {
			user.ImageId = &randomFile.Id
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
