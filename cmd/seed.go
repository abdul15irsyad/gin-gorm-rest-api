package main

import (
	"gin-gorm-rest-api/config"
	"gin-gorm-rest-api/database/seeders"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	databaseConfig := config.NewDatabaseConfig()

	seeders.FileSeeder(databaseConfig.DB)
	seeders.RoleSeeder(databaseConfig.DB)
	seeders.UserSeeder(databaseConfig.DB)
	log.Println("all seeders executed")

	sqlDB, _ := databaseConfig.DB.DB()
	sqlDB.Close()
	os.Exit(0)
}
