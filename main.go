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
		log.Fatal("error loading `.env` file: " + err.Error())
	}
	database.InitDatabase()

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.SetTrustedProxies([]string{})
	router.Static("/assets", "./assets")

	// init routes
	routes.AuthRoutes(router)
	routes.RootRoutes(router)
	routes.UserRoutes(router)
	routes.FileRoutes(router)

	// listen on port
	Port := os.Getenv("PORT")
	Env := os.Getenv("ENV")
	log.Println("environment=" + Env)
	router.Run("localhost:" + fmt.Sprint(Port))
}
