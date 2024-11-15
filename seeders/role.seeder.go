package seeders

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/models"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func RoleSeeder(DB *gorm.DB) {
	// check seeder
	const seederName = "RoleSeeder"
	var (
		seeder models.Seeder
	)
	result := DB.Where("name = ?", seederName).First(&seeder)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(fmt.Sprint(seederName), "already executed")
		return
	}

	// seeder data
	roles := []models.Role{}
	newUuid, _ := uuid.Parse("92833737-af8f-4af2-993d-74ec5b235109")
	role := models.Role{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "Administrator",
		Slug:      slug.Make("Administrator"),
	}
	desc := "a role that can do everything"
	role.Desc = &desc
	roles = append(roles, role)
	newUuid, _ = uuid.Parse("3ed4e622-4642-499a-b711-fb86a458f098")
	role = models.Role{
		BaseModel: models.BaseModel{Id: newUuid},
		Name:      "User",
		Slug:      slug.Make("User"),
		Desc:      nil,
	}
	roles = append(roles, role)
	result = DB.Create(&roles)
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
