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

	seeders.FileSeeder()
	seeders.RoleSeeder()
	seeders.UserSeeder()
	log.Println("all seeders executed")

	os.Exit(0)
}
