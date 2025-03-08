package main

import (
	"gin-gorm-rest-api/lib"
	"gin-gorm-rest-api/seeders"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	libDB := lib.NewDatabase()

	seeders.FileSeeder(libDB.Database)
	seeders.RoleSeeder(libDB.Database)
	seeders.UserSeeder(libDB.Database)
	log.Println("all seeders executed")

	db, _ := libDB.Database.DB()
	db.Close()
	os.Exit(0)
}
