package main

import (
	"belajar-gin/database"
	"belajar-gin/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading `.env` file: %s", err)
	}
	router := gin.Default()
	database.Database()

	// init routes
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	routes.InitRoutes(router)

	// listen on port
	PORT := os.Getenv("PORT")
	router.Run("localhost:" + PORT)
}
