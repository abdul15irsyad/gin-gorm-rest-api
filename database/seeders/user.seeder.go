package seeders

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"log"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func UserSeeder(DB *gorm.DB) {
	// check seeder
	const seederName = "UserSeeder"
	var (
		seeder models.Seeder
	)
	result := DB.Where("name = ?", seederName).First(&seeder)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(fmt.Sprint(seederName), "already executed")
		return
	}

	// seeder data
	var images []models.File
	DB.Where("filename ILIKE ?", "%"+"dummy-profile"+"%").Find(&images)
	users := []models.User{}
	administratorRoleId, _ := uuid.Parse("92833737-af8f-4af2-993d-74ec5b235109")
	userRoleId, _ := uuid.Parse("3ed4e622-4642-499a-b711-fb86a458f098")
	newUuid, _ := uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
	hashedPassword, _ := utils.HashPassword("Qwerty123")
	user := models.User{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "irsyad",
		Email:     "irsyad@email.com",
		Password:  string(hashedPassword),
		ImageId:   &images[0].Id,
		RoleId:    administratorRoleId,
	}
	users = append(users, user)
	newUuid, _ = uuid.Parse("0e3cec4a-206c-4dfb-96f6-3f6b85db9543")
	hashedPassword, _ = utils.HashPassword("Qwerty123")
	user = models.User{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "abdul",
		Email:     "abdul@email.com",
		Password:  string(hashedPassword),
		ImageId:   &images[2].Id,
		RoleId:    userRoleId,
	}
	users = append(users, user)
	usersLength := len(users)
	for i := 0; i < 30-usersLength; i++ {
		randomUuid, _ := uuid.NewRandom()
		hashedPassword, _ := utils.HashPassword("Qwerty123")
		randomFile := utils.RandomSlice(images)
		name := faker.Name()
		user := models.User{
			BaseModel: models.BaseModel{Id: randomUuid},
			Name:      name,
			Email:     strings.ReplaceAll(slug.Make(name), "-", "") + "@email.com",
			Password:  string(hashedPassword),
			RoleId:    userRoleId,
		}
		if utils.RandomSlice([]bool{true, false}) {
			user.ImageId = &randomFile.Id
		}
		users = append(users, user)
	}
	result = DB.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	// add to seeder table
	DB.Create(&models.Seeder{
		Name:      seederName,
		CreatedAt: time.Now(),
	})
	log.Println(fmt.Sprint(seederName) + " executed")
}
