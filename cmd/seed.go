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

	db, _ := databaseConfig.DB.DB()
	db.Close()
	os.Exit(0)
}
