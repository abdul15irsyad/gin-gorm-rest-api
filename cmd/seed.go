package main

import (
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/seeders"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	databaseConfig := configs.NewDatabaseConfig()

	seeders.FileSeeder(databaseConfig.DB)
	seeders.RoleSeeder(databaseConfig.DB)
	seeders.UserSeeder(databaseConfig.DB)
	log.Println("all seeders executed")

	sqlDB, _ := databaseConfig.DB.DB()
	sqlDB.Close()
	os.Exit(0)
}
