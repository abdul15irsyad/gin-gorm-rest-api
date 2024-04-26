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
	// check seeder
	const name = "UserSeeder"
	var (
		seeder models.Seeder
	)
	result := database.DB.Where("name = ?", name).First(&seeder)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}

	// seeder data
	var images []models.File
	database.DB.Where("filename ILIKE ?", "%"+"dummy-profile"+"%").Find(&images)
	users := []models.User{}
	newUuid, _ := uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
	hashedPassword, _ := utils.HashPassword("Qwerty123")
	administratorRoleId, _ := uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
	user := models.User{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "irsyad",
		Email:     "irsyad@email.com",
		Password:  string(hashedPassword),
		ImageId:   &images[0].Id,
		RoleId:    administratorRoleId,
	}
	users = append(users, user)
	usersLength := len(users)
	userRoleId, _ := uuid.Parse("3ed4e622-4642-499a-b711-fb86a458f098")
	for i := 1; i <= 30-usersLength; i++ {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := utils.HashPassword("Qwerty123")
		randomFile := *utils.RandomArray(images)
		user := models.User{
			BaseModel: models.BaseModel{Id: randomUuid},
			Name:      "User " + fmt.Sprint(i),
			Email:     "user" + fmt.Sprint(i) + "@email.com",
			Password:  string(hashedPassword),
			RoleId:    userRoleId,
		}
		if *utils.RandomArray([]bool{true, false}) {
			user.ImageId = &randomFile.Id
		}
		users = append(users, user)
	}
	result = database.DB.Create(&users)
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
