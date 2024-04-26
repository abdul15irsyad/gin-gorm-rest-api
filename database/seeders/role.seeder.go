package seeders

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func RoleSeeder() {
	// check seeder
	const name = "RoleSeeder"
	var (
		seeder models.Seeder
	)
	result := database.DB.Where("name = ?", name).First(&seeder)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return
	}

	// seeder data
	roles := []models.Role{}
	newUuid, _ := uuid.Parse("b0e18329-a7c9-4bea-9efc-72b34818ff14")
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
	result = database.DB.Create(&roles)
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
