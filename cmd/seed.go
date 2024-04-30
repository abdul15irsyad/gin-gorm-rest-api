package main

import (
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/database/seeders"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database.InitDatabase()

	seeders.FileSeeder(database.DB)
	seeders.RoleSeeder(database.DB)
	seeders.UserSeeder(database.DB)
	log.Println("all seeders executed")

	sqlDB, _ := database.DB.DB()
	sqlDB.Close()
	os.Exit(0)
}
