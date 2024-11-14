package main

import (
	"fmt"
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

	Env := os.Getenv("ENV")
	if Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
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

	InitDependencies(router)

	// listen on port
	Port := os.Getenv("PORT")
	log.Println("environment=" + Env)
	router.Run("localhost:" + fmt.Sprint(Port))
}
