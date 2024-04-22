package main

import (
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading `.env` file: %s", err)
	}
	router := gin.Default()
	database.InitDatabase()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	// init routes
	routes.RootRoutes(router)
	routes.UserRoutes(router)
	routes.AuthRoutes(router)

	// listen on port
	Port := os.Getenv("PORT")
	Env := os.Getenv("ENV")
	log.Println("environment=" + Env)
	router.Run("localhost:" + fmt.Sprint(Port))
}
