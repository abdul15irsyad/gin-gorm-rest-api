package main

import (
	"belajar-gin/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// init routes
	routers.V1(router)

	// listen on port 8040
	router.Run(":8040")
}
