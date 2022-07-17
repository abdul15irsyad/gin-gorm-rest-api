package main

import (
	"belajar-gin/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// make api version
	v1 := router.Group("/v1")

	v1.GET("/", handlers.RootHandler)
	v1.GET("/mahasiswa", handlers.GetAllMahasiswa)
	v1.GET("/mahasiswa/:id", handlers.GetMahasiswa)

	router.Run(":8040")
}
