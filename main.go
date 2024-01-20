package main

import (
	"belajar-gin/configs"
	"belajar-gin/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	configs.Database()

	// init routes
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	routes.Root(router);
	routes.Student(router);

	// listen on port 8040
	router.Run(":8040")
}
